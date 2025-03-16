package secure

import (
	"cmp"
	"encoding/json"
	"log/slog"
	"net/http"
	"slices"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
)

type TagSet struct {
	Tag          string `json:"tag"`
	NumQuestions int    `json:"num-questions"`
}

func GetAllTags(userid int64, s Server, w http.ResponseWriter, r *http.Request) {
	questions, err := func() ([]sql_queries.GetAllQuestionsRow, error) {
		server.DbLock.RLock()
		defer server.DbLock.RUnlock()
		res, err := s.Db.GetAllQuestions(r.Context(), userid)
		return res, err
	}()
	if err != nil {
		slog.Error("can't get questions by userid", "uuid", userid, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tagsMap := make(map[string][]int64)
	for _, question := range questions {

		tags, err := func() ([]string, error) {
			server.DbLock.RLock()
			defer server.DbLock.RUnlock()
			tags, err := s.Db.GetTagsByQuestion(r.Context(), question.QuestionID)
			return tags, err
		}()
		if err != nil {
			slog.Error("can't get questions by userid", "uuid", userid, "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
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
	data, err := json.Marshal(tagSet)
	if err != nil {
		panic(err)
	}
	w.Write(data)
}
