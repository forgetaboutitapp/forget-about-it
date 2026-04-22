package do

import (
	"cmp"
	"context"
	"log/slog"
	"slices"

	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	v1 "github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/client_server/v1"
)

func GetAllTags(ctx context.Context, user sql_queries.User, s *Server, _ *v1.GetAllTagsRequest) *v1.GetAllTagsResponse {
	questions, err := s.Db.GetAllQuestions(ctx, user.UserID)

	if err != nil {
		slog.Error("can't get questions by userid", "userid", user.UserID, "err", err)
		return &v1.GetAllTagsResponse{
			Result: &v1.GetAllTagsResponse_Error{
				Error: &v1.ErrorMessage{Error: "Can't get all questions"},
			},
		}
	}
	tagsMap := make(map[string][]int64)
	for _, question := range questions {
		tags, err := s.Db.GetTagsByQuestion(ctx, question.QuestionID)
		if err != nil {
			slog.Error("can't get tags by questionsid", "questionid", question.QuestionID, "err", err)
			return &v1.GetAllTagsResponse{
				Result: &v1.GetAllTagsResponse_Error{
					Error: &v1.ErrorMessage{Error: "Can't get all questions"},
				},
			}
		}
		for _, tag := range tags {
			tagsMap[tag] = append(tagsMap[tag], question.QuestionID)
		}
	}
	var tagSet []*v1.Tag
	for tag, questions := range tagsMap {
		tagSet = append(tagSet, &v1.Tag{Tag: tag, TotalQuestions: uint32(len(questions))})
	}
	slices.SortFunc(tagSet, func(a, b *v1.Tag) int {
		return cmp.Compare(a.Tag, b.Tag)
	})
	allAlgos, err := s.Db.GetSpacingAlgorithms(ctx)
	if err != nil {
		slog.Error("can't get algos", "err", err)
		return &v1.GetAllTagsResponse{
			Result: &v1.GetAllTagsResponse_Error{
				Error: &v1.ErrorMessage{Error: "Can't get spacing algorithm"},
			},
		}
	}
	slog.Info("can run", "len(algos)", len(allAlgos))
	return &v1.GetAllTagsResponse{
		Result: &v1.GetAllTagsResponse_Ok{
			Ok: &v1.GetAllTags{
				Tags:   tagSet,
				CanRun: len(allAlgos) > 0,
			},
		},
	}
}
