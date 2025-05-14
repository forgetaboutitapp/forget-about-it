package do

import (
	"cmp"
	"context"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/client_to_server"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/server_to_client"
	"log/slog"
	"slices"
)

func GetAllTags(ctx context.Context, userid int64, s Server, _ *client_to_server.GetAllTags) *server_to_client.Message {
	questions, err := s.Db.GetAllQuestions(ctx, userid)

	if err != nil {
		slog.Error("can't get questions by userid", "uuid", userid, "err", err)
		return makeError("Can't get all questions")
	}
	tagsMap := make(map[string][]int64)
	for _, question := range questions {

		tags, err := s.Db.GetTagsByQuestion(ctx, question.QuestionID)
		if err != nil {
			slog.Error("can't get tags by questionsid", "questionid", question.QuestionID, "err", err)
			return makeError("Can't get all questions")
		}
		for _, tag := range tags {
			tagsMap[tag] = append(tagsMap[tag], question.QuestionID)
		}
	}
	var tagSet []*server_to_client.Tag
	for tag, questions := range tagsMap {
		tagSet = append(tagSet, &server_to_client.Tag{Tag: tag, TotalQuestions: uint32(len(questions))})
	}
	slices.SortFunc(tagSet, func(a, b *server_to_client.Tag) int {
		return cmp.Compare(a.Tag, b.Tag)
	})
	allAlgos, err := s.Db.GetSpacingAlgorithms(ctx)
	if err != nil {
		slog.Error("can't get algos", "err", err)
		return makeError("Can't get spacing algorithm")
	}
	slog.Info("can run", "len(algos)", len(allAlgos))
	return &server_to_client.Message{ReturnMessage: &server_to_client.Message_OkMessage{OkMessage: &server_to_client.OkMessage{OkMessage: &server_to_client.OkMessage_GetAllTags{GetAllTags: &server_to_client.GetAllTags{Tags: tagSet, CanRun: len(allAlgos) > 0}}}}}
}
