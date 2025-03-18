package dbUtils

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"path"
	"slices"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/cmd/server/apis/secure"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
)

var ErrCantOpen = errors.New("cannot open sqlite file")
var ErrCantEnableForeignKeys = errors.New("cannot enable foreign keys")
var ErrCantGetDBVersion = errors.New("cannot get db version")
var ErrCantDoMigrations = errors.New("cannot do migrations")

func OpenDatabase(ctx context.Context) (*sql.DB, error) {
	slog.Info("Opening DB", "db file", server.DBFilename)
	db, err := sql.Open("sqlite", server.DBFilename)
	if err != nil {
		return nil, errors.Join(ErrCantOpen, err)
	}
	err = enableForeignKeys(ctx, db)
	if err != nil {
		return nil, errors.Join(ErrCantEnableForeignKeys, err)
	}

	v, err := GetDbVersion(ctx, db)
	if err != nil {
		return nil, errors.Join(ErrCantGetDBVersion, err)
	}
	_, err = DoMigrations(ctx, db, v)
	if err != nil {
		return nil, errors.Join(ErrCantDoMigrations, err)
	}
	return db, nil
}
func enableForeignKeys(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, "PRAGMA foreign_keys=ON")
	return err
}

type executor interface {
	ExecContext(context context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(context context.Context, query string, args ...any) *sql.Row
}

func GetDbVersion(ctx context.Context, db executor) (int, error) {
	r := db.QueryRowContext(ctx, "PRAGMA user_version")
	v := 0
	err := r.Scan(&v)
	return v, err
}
func SetDbVersion(ctx context.Context, db executor, v int) error {
	_, err := db.ExecContext(ctx, "PRAGMA user_version = "+strconv.Itoa(v))
	return err
}

func DoMigrations(ctx context.Context, db *sql.DB, originalVersion int) (int, error) {
	p, err := server.DDL.ReadDir("sql/migrations")
	if err != nil {
		log.Panicf("cannot read dir(%s): %s", p, err.Error())
	}
	versionsToMigrate := []int{}
	for _, file := range p {
		fileNum, err := strconv.Atoi(strings.Split(file.Name(), ".")[0])
		if err != nil {
			log.Panicf("cannot convert first of (%s) to int: %s", strings.Split(file.Name(), "."), err.Error())
		}
		if fileNum > originalVersion {
			versionsToMigrate = append(versionsToMigrate, fileNum)
		}
	}
	slices.Sort(versionsToMigrate)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, errors.Join(secure.ErrCantInitTransaction, err)
	}
	defer tx.Rollback()
	v := 0
	for _, curVersion := range versionsToMigrate {
		pathOfFile := path.Join("sql/migrations", fmt.Sprintf("%03d.sql", curVersion))
		p, err := server.DDL.ReadFile(pathOfFile)
		if err != nil {
			log.Panicf("cannot read file (%s): %s", pathOfFile, err.Error())
		}
		if _, err := tx.ExecContext(ctx, string(p)); err != nil {
			log.Panicf("cannot do transaction: %s", err.Error())
		}

		err = SetDbVersion(ctx, tx, curVersion)
		if err != nil {
			log.Panicf("cannot set db version (migrating to: %d): %s", curVersion, err.Error())
		}
		v, err = GetDbVersion(ctx, tx)
		if err != nil {
			log.Panicf("cannot get db version (migrating to: %d): %s", curVersion, err.Error())
		}
	}
	err = tx.Commit()
	if err != nil {
		return 0, errors.Join(secure.ErrCantCommit, err)

	}
	return v, nil

}

func GetDB() (*sql_queries.Queries, *sql.DB) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	defer cancel()
	db, err := OpenDatabase(ctx)
	db.SetMaxOpenConns(1)
	if err != nil {
		panic(err)
	}
	q := sql_queries.New(db)

	return q, db
}
