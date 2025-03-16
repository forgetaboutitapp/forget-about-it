package secure

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
)

func GetAllQuestions(userid int64, s Server, w http.ResponseWriter, r *http.Request) {
	res, err := func() ([]sql_queries.GetAllQuestionsRow, error) {
		server.DbLock.RLock()
		defer server.DbLock.RUnlock()
		val, err := s.Db.GetAllQuestions(r.Context(), userid)
		return val, err
	}()
	if err != nil {
		slog.Error("can't get all questions", "userid", userid, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	flashcards := []Flashcard{}

	for _, q := range res {

		tags, err := func() ([]string, error) {
			server.DbLock.RLock()
			defer server.DbLock.RUnlock()
			tags, err := s.Db.GetTagsByQuestion(r.Context(), q.QuestionID)
			return tags, err
		}()
		if err != nil {
			slog.Error("can't get tag for question", "questionid", q.QuestionID, "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return

		}
		flashcards = append(flashcards, Flashcard{Id: q.QuestionID, Answer: q.Answer, Question: q.Question, Tags: tags})
	}
	slog.Info("Getting all questions", "flashcards", flashcards)
	jsonText, err := json.Marshal(flashcards)
	if err != nil {
		panic(err)
	}
	w.Write(jsonText)
}
