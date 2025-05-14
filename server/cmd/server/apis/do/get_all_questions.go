package do

import (
	"context"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/client_to_server"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/server_to_client"
	"log/slog"
)

func GetAllQuestions(ctx context.Context, userid int64, s Server, _ *client_to_server.GetAllQuestions) *server_to_client.Message {
	res, err := s.Db.GetAllQuestions(ctx, userid)
	if err != nil {
		slog.Error("can't get all questions", "userid", userid, "err", err)
		return makeError("Can't get all questions")
	}

	var flashcards []*server_to_client.Flashcard

	for _, q := range res {
		tags, err := s.Db.GetTagsByQuestion(ctx, q.QuestionID)
		if err != nil {
			slog.Error("can't get tag for question", "questionid", q.QuestionID, "err", err)
			return makeError("Can't get all questions")
		}
		flashcards = append(flashcards, &server_to_client.Flashcard{
			Id:          uint32(q.QuestionID),
			Answer:      q.Answer,
			Question:    q.Question,
			Tags:        tags,
			MemoHint:    q.MemoHint,
			Explanation: q.Explanation,
		})
	}
	return &server_to_client.Message{ReturnMessage: &server_to_client.Message_OkMessage{OkMessage: &server_to_client.OkMessage{OkMessage: &server_to_client.OkMessage_GetAllQuestions{GetAllQuestions: &server_to_client.GetAllQuestions{Flashcards: flashcards}}}}}
}
