package dbUtils

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"path"
	"slices"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"

	"github.com/forgetaboutitapp/forget-about-it/server"
)

func OpenDatabase(ctx context.Context) (*sql.DB, error) {
	slog.Info("Opening DB", "db file", server.DBFilename)
	db, err := sql.Open("sqlite", server.DBFilename)
	if err != nil {
		return nil, fmt.Errorf("cannot open sqlite file: %w", err)
	}
	err = enableForeignKeys(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("cannot enable foreign keys: %w", err)
	}

	v, err := GetDbVersion(ctx, db)
	if err != nil {
		return nil, fmt.Errorf("cannot get db version: %w", err)
	}
	_, err = DoMigrations(ctx, db, v)
	if err != nil {
		return nil, fmt.Errorf("cannot do migrations: %w", err)
	}
	return db, nil
}
func enableForeignKeys(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, "PRAGMA foreign_keys=ON")
	return err
}

type executor interface {
	ExecContext(context context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(context context.Context, query string, args ...interface{}) *sql.Row
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
		return 0, fmt.Errorf("cannot read dir(%s): %w", p, err)
	}
	versionsToMigrate := []int{}
	for _, file := range p {
		fileNum, err := strconv.Atoi(strings.Split(file.Name(), ".")[0])
		if err != nil {
			return 0, fmt.Errorf("cannot convert first of (%s) to int: %w", strings.Split(file.Name(), "."), err)
		}
		if fileNum > originalVersion {
			versionsToMigrate = append(versionsToMigrate, fileNum)
		}
	}
	slices.Sort(versionsToMigrate)

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return 0, fmt.Errorf("cannot begin transaction: %w", err)
	}
	defer tx.Rollback()
	v := 0
	for _, curVersion := range versionsToMigrate {
		pathOfFile := path.Join("sql/migrations", fmt.Sprintf("%03d.sql", curVersion))
		p, err := server.DDL.ReadFile(pathOfFile)
		if err != nil {
			return 0, fmt.Errorf("cannot read file (%s): %w", pathOfFile, err)
		}
		if _, err := tx.ExecContext(ctx, string(p)); err != nil {
			return 0, fmt.Errorf("cannot do transaction: %w", err)
		}

		err = SetDbVersion(ctx, tx, curVersion)
		if err != nil {
			return 0, fmt.Errorf("cannot set db version (migrating to: %d): %w", curVersion, err)
		}
		v, err = GetDbVersion(ctx, tx)
		if err != nil {
			return 0, fmt.Errorf("cannot get db version (migrating to: %d): %w", curVersion, err)
		}
	}
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("cannot commit: %w", err)

	}
	return v, nil

}
