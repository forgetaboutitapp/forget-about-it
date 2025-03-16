package secure

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
)

func GradeQuestion(userid int64, s Server, w http.ResponseWriter, r *http.Request) {
	slog.Info("Grading question")
	d, err := io.ReadAll(r.Body)
	if err != nil {
		slog.Error("Could not read body", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var mapData map[string]any
	err = json.Unmarshal(d, &mapData)
	if err != nil {
		slog.Error("unable to unmarshal mapData", "data", string(d), "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	questionId := int64((mapData["question-id"]).(float64))
	correct := int64((mapData["question-id"]).(float64))
	result := 0
	if correct != 0 {
		result = 1
	}
	err = func() error {
		server.DbLock.Lock()
		defer server.DbLock.Unlock()
		return s.Db.GradeQuestion(r.Context(), sql_queries.GradeQuestionParams{
			QuestionID: questionId,
			Result:     int64(result),
			Timestamp:  time.Now().UnixMicro(),
		})
	}()
	if err != nil {
		slog.Error("unable to save grades to the database", "params", sql_queries.GradeQuestionParams{
			QuestionID: questionId,
			Result:     int64(result),
			Timestamp:  time.Now().UnixMicro(),
		}, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
