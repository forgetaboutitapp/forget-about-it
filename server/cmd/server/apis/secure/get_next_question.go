package secure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	"github.com/hashicorp/go-set"
)

var ErrCantGetDefaultAlgo = errors.New("can't get default algo")
var ErrCantGetSpacingAlgo = errors.New("can't get spacing algo")
var ErrCantGetNewModBuilder = errors.New("can't get new mod builder")
var ErrCantInstantiate = errors.New("can't instantiate")
var ErrCantGetAllGrades = errors.New("can't get all grades")
var ErrCantGetTagForQuestion = errors.New("can't get tag for question")
var ErrCantGetTagForGrade = errors.New("can't get tag for grade")
var ErrCantGetAddCard = errors.New("can't get add-card function")
var ErrCantGetGradeCard = errors.New("can't get grade-card function")
var ErrCantAllocate = errors.New("can't allocate")
var ErrCantFree = errors.New("can't free")
var ErrCantCallNextCard = errors.New("can't call next card")
var ErrCantReadBytes = errors.New("can't read bytes")
var ErrCantRunWasm = errors.New("can't run wasm")

func GetNextQuestion(ctx context.Context, userid int64, s Server, m map[string]any) (map[string]any, error) {
	slog.Info("Getting next question")
	tagsSent := m["tag-sent"].([]string)
	tagsSentSet := set.From(tagsSent)
	defaultAlgo, err := func() (sql.NullInt64, error) {
		server.DbLock.RLock()
		defer server.DbLock.RUnlock()
		defaultAlgo, err := s.Db.GetDefaultAlgorithm(ctx, userid)
		return defaultAlgo, err
	}()
	if err != nil {
		slog.Error("can't get default algo", "uuid", userid, "err", err)
		return nil, errors.Join(ErrCantGetDefaultAlgo, err)
	}
	algos, err := func() ([]sql_queries.SpacingAlgorithm, error) {
		server.DbLock.RLock()
		defer server.DbLock.RUnlock()
		a, b := s.Db.GetSpacingAlgorithms(ctx)
		return a, b
	}()
	if err != nil {
		slog.Error("can't get spacing algorithm", "uuid", userid, "err", err)
		return nil, errors.Join(ErrCantGetSpacingAlgo, err)

	}
	var algo sql_queries.SpacingAlgorithm
	if defaultAlgo.Valid {

		for _, curAlgo := range algos {
			if curAlgo.AlgorithmID == defaultAlgo.Int64 {
				break
			}
		}
		panic(fmt.Sprintln("default algo does not match a valid algorithm id", algos, defaultAlgo))
	} else if len(algos) > 0 {
		algo = algos[0]
	} else {
		panic("There are no algorithms available")
	}

	allGrades, err := func() ([]sql_queries.QuestionsLog, error) {
		server.DbLock.RLock()
		defer server.DbLock.RUnlock()
		return s.Db.GetAllGrades(ctx, userid)
	}()
	if err != nil {
		slog.Error("can't get all grades", "err", err)
		return nil, errors.Join(ErrCantGetAllGrades, err)
	}
	allQuestions, err := func() ([]sql_queries.GetAllQuestionsRow, error) {
		server.DbLock.RLock()
		defer server.DbLock.RUnlock()
		return s.Db.GetAllQuestions(ctx, userid)
	}()
	if err != nil {
		slog.Error("can't get all questions", "err", err)
		return nil, errors.Join(ErrCanGetQuestions)
	}

	allGradesWithRightTags := []sql_queries.QuestionsLog{}
	allQuestionsWithRightTags := []sql_queries.GetAllQuestionsRow{}

	for _, question := range allQuestions {
		tags, err := func() ([]string, error) {
			server.DbLock.RLock()
			defer server.DbLock.RUnlock()
			tags, err := s.Db.GetTagsByQuestion(ctx, question.QuestionID)
			return tags, err
		}()
		if err != nil {
			slog.Error("can't get tag for question", "questionid", question.QuestionID, "err", err)
			return nil, errors.Join(ErrCantGetTagForQuestion, err)
		}
		tagsSet := set.From(tags)
		found := false
		for _, t := range tagsSentSet.Slice() {
			if tagsSet.Contains(t) {
				found = true
			}
		}
		if found {
			allQuestionsWithRightTags = append(allQuestionsWithRightTags, question)
		}
	}

	for _, grade := range allGrades {
		tags, err := func() ([]string, error) {
			server.DbLock.RLock()
			defer server.DbLock.RUnlock()
			tags, err := s.Db.GetTagsByQuestion(ctx, grade.QuestionID)
			return tags, err
		}()
		if err != nil {
			slog.Error("can't get tag for grade", "questionid", grade.QuestionID, "err", err)
			return nil, errors.Join(ErrCantGetTagForGrade, err)

		}
		tagsSet := set.From(tags)
		found := false
		for _, t := range tagsSentSet.Slice() {
			if tagsSet.Contains(t) {
				found = true
			}
		}
		if found {
			allGradesWithRightTags = append(allGradesWithRightTags, grade)
		}
	}

	ret, err := runAlgorithm(ctx, AlgorithmStruct{
		Alloc:         algo.Alloc,
		ApiVersion:    int(algo.ApiVersion),
		Author:        algo.Author,
		Dealloc:       algo.Dealloc,
		Desc:          algo.Desc.String,
		DownloadUrl:   algo.DownloadUrl,
		Init:          algo.Init,
		License:       algo.License,
		ModuleName:    algo.ModuleName,
		AlgorithmName: algo.AlgorithmName,
		RemoteURL:     algo.RemoteUrl,
		Version:       int(algo.Version),
		WasmBytes:     algo.Wasm,
	}, allQuestionsWithRightTags, allGradesWithRightTags)
	if err != nil {
		slog.Error("cannot run wasm", "algoname", algo.AlgorithmName, "err", err)
		return nil, errors.Join(err, ErrCantRunWasm)
	}
	question := ""
	answer := ""
	found := false
	for _, questionGot := range allQuestions {
		if questionGot.QuestionID == int64(ret.nextCard) {
			question = questionGot.Question
			answer = questionGot.Answer
			found = true
		}
	}
	if !found {
		slog.Error("Cannot find question id", "nextCard", ret.nextCard)
		return nil, ErrCanGetQuestions

	}

	resultMap := map[string]any{"amount-due-cards": ret.lenDueCards, "amount-new-cards": ret.lenNewCards, "amount-non-due-cards": ret.lenNonDueCards, "id": ret.nextCard, "question": question, "answer": answer, "card-type": ret.typeOfNextCard}
	return resultMap, nil
}
