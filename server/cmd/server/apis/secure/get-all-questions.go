package secure

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func GetAllQuestions(userid int64, s Server, w http.ResponseWriter, r *http.Request) {
	res, err := s.Db.GetAllQuestions(r.Context(), userid)
	if err != nil {
		slog.Error("can't get all questions", "userid", userid, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	flashcards := []Flashcard{}
	for _, q := range res {
		tags, err := s.Db.GetTagsByQuestion(r.Context(), q.QuestionID)
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
