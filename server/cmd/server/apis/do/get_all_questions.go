package do

import (
	"context"
	"log/slog"

	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	v1 "github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/client_server/v1"
)

func GetAllQuestions(ctx context.Context, user sql_queries.User, s *Server, _ *v1.GetAllQuestionsRequest) *v1.GetAllQuestionsResponse {
	res, err := s.Db.GetAllQuestions(ctx, user.UserID)
	if err != nil {
		slog.Error("can't get all questions", "userid", user.UserID, "err", err)
		return &v1.GetAllQuestionsResponse{
			Result: &v1.GetAllQuestionsResponse_Error{
				Error: &v1.ErrorMessage{Error: "Can't get all questions"},
			},
		}
	}

	var flashcards []*v1.Flashcard

	for _, q := range res {
		tags, err := s.Db.GetTagsByQuestion(ctx, q.QuestionID)
		if err != nil {
			slog.Error("can't get tag for question", "questionid", q.QuestionID, "err", err)
			return &v1.GetAllQuestionsResponse{
				Result: &v1.GetAllQuestionsResponse_Error{
					Error: &v1.ErrorMessage{Error: "Can't get all questions"},
				},
			}
		}
		flashcards = append(flashcards, &v1.Flashcard{
			Id:          uint32(q.QuestionID),
			Answer:      q.Answer,
			Question:    q.Question,
			Tags:        tags,
			MemoHint:    q.MemoHint,
			Explanation: q.Explanation,
		})
	}
	return &v1.GetAllQuestionsResponse{
		Result: &v1.GetAllQuestionsResponse_Ok{
			Ok: &v1.GetAllQuestions{Flashcards: flashcards},
		},
	}
}
