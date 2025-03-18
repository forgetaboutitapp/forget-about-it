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

func addUser(ctx context.Context, t *testing.T, q *sql_queries.Queries, userid int64) {
	err := q.AddUser(ctx, sql_queries.AddUserParams{UserID: userid, Role: 0, Created: 123})
	if err != nil {
		t.Fatal("cant add new user", err)
	}
}
func addQuestion(ctx context.Context, t *testing.T, q *sql_queries.Queries, userid int64, questionid int64) {
	err := q.AddNewQuestion(ctx, sql_queries.AddNewQuestionParams{QuestionID: questionid, UserID: userid, Question: "q1", Answer: "a1", Enabled: 1})
	if err != nil {
		t.Fatal("cant add new question", err)
	}
}
func addTag(ctx context.Context, t *testing.T, q *sql_queries.Queries, questionid int64, tags []string) {
	for _, tag := range tags {
		err := q.AddNewTag(ctx, sql_queries.AddNewTagParams{QuestionID: questionid, Tag: tag})
		if err != nil {
			t.Fatal("cant add new question", err)
		}
	}
}

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
