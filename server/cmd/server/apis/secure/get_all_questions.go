package secure

import (
	"context"
	"errors"
	"log/slog"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
)

var ErrCanGetQuestions = errors.New("can't get questions")
var ErrCantGetTag = errors.New("can't get tag")

func GetAllQuestions(ctx context.Context, userid int64, s Server, m map[string]any) (map[string]any, error) {
	res, err := func() ([]sql_queries.GetAllQuestionsRow, error) {
		server.DbLock.RLock()
		defer server.DbLock.RUnlock()
		val, err := s.Db.GetAllQuestions(ctx, userid)
		return val, err
	}()
	if err != nil {
		slog.Error("can't get all questions", "userid", userid, "err", err)
		return nil, errors.Join(ErrCanGetQuestions, err)
	}

	flashcards := []Flashcard{}

	for _, q := range res {

		tags, err := func() ([]string, error) {
			server.DbLock.RLock()
			defer server.DbLock.RUnlock()
			tags, err := s.Db.GetTagsByQuestion(ctx, q.QuestionID)
			return tags, err
		}()
		if err != nil {
			slog.Error("can't get tag for question", "questionid", q.QuestionID, "err", err)
			return nil, errors.Join(ErrCantGetTag, err)
		}
		flashcards = append(flashcards, Flashcard{Id: q.QuestionID, Answer: q.Answer, Question: q.Question, Tags: tags})
	}
	return map[string]any{"flashcards": flashcards}, nil
}
