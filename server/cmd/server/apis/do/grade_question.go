package do

import (
	"context"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/client_to_server"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/server_to_client"
	"log/slog"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
)

func GradeQuestion(ctx context.Context, userid int64, token string, s Server, arg *client_to_server.GradeQuestion) *server_to_client.Message {
	slog.Info("Grading question", "arg", arg)

	questionId := arg.Questionid

	correct := arg.Correct
	result := 0
	if correct {
		result = 1
	} else {
		result = 0
	}
	params := sql_queries.GradeQuestionParams{
		QuestionID: int64(questionId),
		Result:     int64(result),
		Timestamp:  time.Now().UTC().Unix(),
	}
	err := s.Db.GradeQuestion(ctx, params)
	if err != nil {
		slog.Error("unable to save grades to the database", "params", params, "err", err)
		return makeError("Internal Server Error")
	}
	return makeOk(&server_to_client.OkMessage{OkMessage: &server_to_client.OkMessage_GradeQuestion{GradeQuestion: &server_to_client.GradeQuestion{}}})
}
