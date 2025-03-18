package secure

import (
	"cmp"
	"context"
	"errors"
	"log/slog"
	"slices"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
)

type TagSet struct {
	Tag          string `json:"tag"`
	NumQuestions int    `json:"num-questions"`
}

func GetAllTags(ctx context.Context, userid int64, s Server, m map[string]any) (map[string]any, error) {
	questions, err := func() ([]sql_queries.GetAllQuestionsRow, error) {
		server.DbLock.RLock()
		defer server.DbLock.RUnlock()
		res, err := s.Db.GetAllQuestions(ctx, userid)
		return res, err
	}()
	if err != nil {
		slog.Error("can't get questions by userid", "uuid", userid, "err", err)
		return nil, errors.Join(ErrCanGetQuestions, err)
	}
	tagsMap := make(map[string][]int64)
	for _, question := range questions {

		tags, err := func() ([]string, error) {
			server.DbLock.RLock()
			defer server.DbLock.RUnlock()
			tags, err := s.Db.GetTagsByQuestion(ctx, question.QuestionID)
			return tags, err
		}()
		if err != nil {
			slog.Error("can't get questions by userid", "uuid", userid, "err", err)
			return nil, errors.Join(ErrCanGetQuestions, err)
		}
		for _, tag := range tags {
			tagsMap[tag] = append(tagsMap[tag], question.QuestionID)
		}
	}
	tagSet := []TagSet{}
	for tag, questions := range tagsMap {
		tagSet = append(tagSet, TagSet{Tag: tag, NumQuestions: len(questions)})
	}
	slices.SortFunc(tagSet, func(a, b TagSet) int {
		return cmp.Compare(a.Tag, b.Tag)
	})
	return map[string]any{"tag-set": tagSet}, nil
}
