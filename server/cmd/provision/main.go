package main

import (
	"context"
	"crypto/rand"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"math"
	"math/big"
	mathrand "math/rand/v2"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/adrg/xdg"
	"github.com/forgetaboutitapp/forget-about-it/server"
	dbUtils "github.com/forgetaboutitapp/forget-about-it/server/pkg/db_utils"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	uuidUtils "github.com/forgetaboutitapp/forget-about-it/server/pkg/uuid_utils"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

var ErrUserGeneration = errors.New("user generation error")

var ErrMnemonicGeneration = errors.New("mnemonic generation error")
var ErrHostRequired = errors.New("host is required: pass -host or configure it in db")

func main() {
	ctx := context.Background()
	dbLocation := flag.String("location", filepath.Join(xdg.StateHome, "forget-about-it.sqlite3"), "sqlite3 file location")
	host := flag.String("host", "", "server host")
	flag.Parse()
	server.DBFilename = *dbLocation
	fmt.Println(server.DBFilename)
	db, err := dbUtils.OpenDatabase(ctx)
	if err != nil {
		panic(err)
	}
	q := sql_queries.New(db)
	address, err := SetOrGetHost(ctx, q, *host)
	if err != nil {
		slog.Error("cannot save host", "err", err)
		os.Exit(1)
	}
	address = strings.TrimPrefix(address, "http://")

	adminUser, err := q.GetUser(ctx, 0)
	if err != nil {
		panic(err)
	}
	if len(adminUser) > 1 {
		panic(fmt.Sprintf("There should not be more than 1 user but there were %d", len(adminUser)))
	} else if len(adminUser) == 0 {
		mnemonic, id, err := AddUser(q)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Your 12 word mnemonic is: %s\n", mnemonic)
		fmt.Printf("Your copyable login is: %s\n", id)
		if err := PrintQRCode(fmt.Sprintf("%s;%s", address, id)); err != nil {
			panic(err)
		}
	} else if len(adminUser) == 1 {
		mnemonic, id, err := AddLogin(context.TODO(), q, adminUser[0])
		if err != nil {
			panic(err)
		}
		fmt.Printf("Your 12 word mnemonic is: %s\n", mnemonic)
		fmt.Printf("Your copyable login is: %s\n", id)
		if err := PrintQRCode(fmt.Sprintf("%s;%s", address, id)); err != nil {
			panic(err)
		}
	}
}

func SetOrGetHost(ctx context.Context, q *sql_queries.Queries, hostFlagValue string) (string, error) {
	hostFlagValue = strings.TrimSpace(hostFlagValue)
	if hostFlagValue != "" {
		err := q.SetConfigValue(ctx, sql_queries.SetConfigValueParams{
			Key:   "host",
			Value: hostFlagValue,
		})
		if err != nil {
			return "", err
		}
		return hostFlagValue, nil
	}

	host, err := q.GetConfigValue(ctx, "host")
	if errors.Is(err, sql.ErrNoRows) {
		hostname, _ := os.Hostname()
		return fmt.Sprintf("%s.local:80", hostname), nil
	}
	if err != nil {
		return "", err
	}
	host = strings.TrimSpace(host)
	if host == "" {
		return "", ErrHostRequired
	}
	return host, nil
}

func AddUser(q *sql_queries.Queries) (string, string, error) {
	bigUserid, err := rand.Int(rand.Reader, big.NewInt(int64(math.MaxInt64)))
	if err != nil {
		return "", "", errors.Join(ErrUserGeneration, err)
	}
	userid := bigUserid.Int64()
	q.AddUser(context.Background(), sql_queries.AddUserParams{
		UserID: userid,
		Role:   0,
	})
	loginUuid := uuid.New()
	q.AddLogin(context.Background(), sql_queries.AddLoginParams{
		LoginUuid:         loginUuid.String(),
		UserID:            userid,
		DeviceDescription: "Initial Device",
		Created:           time.Now().UTC().Unix(),
		IndexID:           int64(mathrand.Uint32()),
	})
	mnemonic, err := uuidUtils.NewMnemonicFromUuid(loginUuid)
	if err != nil {
		slog.Error("Cannot generate mnemonic", "err", err)
		return "", "", errors.Join(ErrMnemonicGeneration, err)
	}
	return mnemonic, loginUuid.String(), nil
}

func AddLogin(ctx context.Context, q *sql_queries.Queries, id int64) (string, string, error) {

	newLoginUuid := uuid.New()
	q.AddLogin(context.Background(), sql_queries.AddLoginParams{
		LoginUuid:         newLoginUuid.String(),
		UserID:            id,
		DeviceDescription: fmt.Sprintf("Added on %s", time.Now().UTC().Format(time.DateTime)),
		Created:           time.Now().UTC().Unix(),
		IndexID:           int64(mathrand.Uint32()),
	})
	m, err := uuidUtils.NewMnemonicFromUuid(newLoginUuid)
	if err != nil {
		slog.Error("cannot create a 12 word mnemonic", "err", err)
		return "", "", errors.Join(ErrMnemonicGeneration, err)
	}
	return m, newLoginUuid.String(), nil
}

func PrintQRCode(value string) error {
	qr, err := qrcode.New(value, qrcode.Medium)
	if err != nil {
		return err
	}
	bitmap := qr.Bitmap()
	fmt.Println("Scan QR:")
	for _, row := range bitmap {
		for _, isBlack := range row {
			if isBlack {
				fmt.Print("##")
			} else {
				fmt.Print("  ")
			}
		}
		fmt.Println()
	}
	return nil
}
