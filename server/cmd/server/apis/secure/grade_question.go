package secure

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
)

var ErrCantSaveGrades = errors.New("can't save grades to the database")
var ErrCorrectValIsNotBoolean = errors.New("is-correct val is not boolean")
var ErrMapIsInvalidType = errors.New("map is invalid type")

func GradeQuestion(ctx context.Context, userid int64, s Server, m map[string]any) (map[string]any, error) {
	slog.Info("Grading question")

	questionId, correctType := m["question-id"].(float64)
	if !correctType {
		slog.Error("map is invalid type", "question-id", m["question-id"], "m", m)
		return nil, ErrMapIsInvalidType
	}
	correct, correctType := (m["correct"]).(bool)
	if !correctType {
		slog.Error("map is invalid type", "correct", m["correct"], "m", m)
		return nil, ErrMapIsInvalidType
	}
	result := 0
	if correct {
		result = 1
	} else {
		result = 0
	}
	err := func() error {
		server.DbLock.Lock()
		defer server.DbLock.Unlock()
		return s.Db.GradeQuestion(ctx, sql_queries.GradeQuestionParams{
			QuestionID: int64(questionId),
			Result:     int64(result),
			Timestamp:  time.Now().Unix(),
		})
	}()
	if err != nil {
		slog.Error("unable to save grades to the database", "params", sql_queries.GradeQuestionParams{
			QuestionID: int64(questionId),
			Result:     int64(result),
			Timestamp:  time.Now().Unix(),
		}, "err", err)
		return nil, errors.Join(ErrCantSaveGrades, err)
	}
	return nil, nil
}
