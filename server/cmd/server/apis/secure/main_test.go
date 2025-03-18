package secure_test

import (
	"context"
	"database/sql"
	"log/slog"
	"testing"

	"github.com/forgetaboutitapp/forget-about-it/server"
	dbUtils "github.com/forgetaboutitapp/forget-about-it/server/pkg/db_utils"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
)

func Init(t *testing.T) (*sql_queries.Queries, *sql.DB) {
	slog.SetLogLoggerLevel(slog.LevelError)
	server.DBFilename = ":memory:"
	q, db := dbUtils.GetDB()
	err := q.AddUser(context.Background(), sql_queries.AddUserParams{UserID: 1, Role: 0, Created: 1000})
	if err != nil {
		t.Fatalf("adduser failed: %s", err.Error())
	}
	err = q.AddUser(context.Background(), sql_queries.AddUserParams{UserID: 2, Role: 0, Created: 1000})
	if err != nil {
		t.Fatalf("adduser failed: %s", err.Error())
	}

	err = q.AddUser(context.Background(), sql_queries.AddUserParams{UserID: 3, Role: 0, Created: 1000})
	if err != nil {
		t.Fatalf("adduser failed: %s", err.Error())
	}
	return q, db
}
