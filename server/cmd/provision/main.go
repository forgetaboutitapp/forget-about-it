package main

import (
	"context"
	"fmt"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/dbUtils"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	uuidUtils "github.com/forgetaboutitapp/forget-about-it/server/pkg/uuid_utils"
	"github.com/google/uuid"
)

func main() {
	fmt.Println(server.DBFilename)
	db, err := dbUtils.OpenDatabase(context.Background())
	if err != nil {
		panic(err)
	}
	q := sql_queries.New(db)

	adminUser, err := q.GetUser(context.Background(), 0)
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
	} else if len(adminUser) == 1 {
		mnemonic, id, err := AddLogin(context.TODO(), q, adminUser[0])
		if err != nil {
			panic(err)
		}
		fmt.Printf("Your 12 word mnemonic is:\n%s", mnemonic)
		fmt.Printf("Your copyable login is: %s\n", id)
	}
}

func AddUser(q *sql_queries.Queries) (string, string, error) {
	userUUID := uuid.New()

	q.AddUser(context.Background(), sql_queries.AddUserParams{
		UserUuid: userUUID.String(),
		Role:     0,
	})

	loginUuid := uuid.New()
	q.AddLogin(context.Background(), sql_queries.AddLoginParams{
		LoginUuid:         loginUuid.String(),
		UserUuid:          userUUID.String(),
		DeviceDescription: "Initial Device",
		Created:           time.Now().Unix(),
	})
	mnemonic, err := uuidUtils.NewMnemonicFromUuid(loginUuid)
	if err != nil {
		return "", "", fmt.Errorf("cannot create mnemonic: %w", err)
	}
	return mnemonic, loginUuid.String(), nil
}

func AddLogin(ctx context.Context, q *sql_queries.Queries, id string) (string, string, error) {
	userUUID, err := uuid.Parse(id)
	if err != nil {
		return "", "", fmt.Errorf("user id (%s) is not a valid uuid: %w", id, err)
	}
	newLoginUuid := uuid.New()
	q.AddLogin(context.Background(), sql_queries.AddLoginParams{
		LoginUuid:         newLoginUuid.String(),
		UserUuid:          userUUID.String(),
		DeviceDescription: fmt.Sprintf("Added on %s", time.Now().UTC().Format(time.DateTime)),
		Created:           time.Now().Unix(),
	})
	m, err := uuidUtils.NewMnemonicFromUuid(newLoginUuid)
	if err != nil {
		return "", "", fmt.Errorf("cannot create a 12 word mnemonic from %s: %w", id, err)
	}
	return m, newLoginUuid.String(), nil
}
