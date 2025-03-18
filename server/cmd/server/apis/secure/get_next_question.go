package secure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	"github.com/hashicorp/go-set"
	"github.com/tetratelabs/wazero"
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

func GetNextQuestion(ctx context.Context, userid int64, s Server, m map[string]any) (map[string]any, error) {
	slog.Info("Getting next question")
	tagsSent := m["tagSent"].([]string)
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
	algos, err := func() ([]sql_queries.GetSpacingAlgorithmsRow, error) {
		server.DbLock.RLock()
		defer server.DbLock.RUnlock()
		a, b := s.Db.GetSpacingAlgorithms(ctx)
		return a, b
	}()
	if err != nil {
		slog.Error("can't get spacing algorithm", "uuid", userid, "err", err)
		return nil, errors.Join(ErrCantGetSpacingAlgo, err)

	}
	var algo sql_queries.GetSpacingAlgorithmsRow
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

	runtime := wazero.NewRuntime(ctx)

	defer runtime.Close(ctx) // This closes everything this Runtime created.
	_, err = runtime.NewHostModuleBuilder(algo.ModuleName).
		Instantiate(ctx)
	if err != nil {
		slog.Error("can't make a new module builder", "algo", algo.AlgorithmID, "err", err)
		return nil, errors.Join(ErrCantGetNewModBuilder, err)
	}
	initializationFunctions := strings.Split(algo.InitializationFunctions, ",")
	modConfig := wazero.NewModuleConfig().WithStartFunctions()
	if len(initializationFunctions) != 0 {
		modConfig = modConfig.WithStartFunctions(initializationFunctions...)
	}
	mod, err := runtime.InstantiateWithConfig(ctx, algo.Algorithm, modConfig)
	if err != nil {
		slog.Error("can't instantiate", "algo", algo.AlgorithmID, "modConfig", modConfig, "err", err)
		return nil, errors.Join(ErrCantInstantiate, err)
	}
	addCard := mod.ExportedFunction("add-card")
	gradeCard := mod.ExportedFunction("grade-card")
	getCard := mod.ExportedFunction("get-cards")
	malloc := mod.ExportedFunction("alloc")
	free := mod.ExportedFunction("dealloc")
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

	for _, question := range allQuestionsWithRightTags {
		_, err = addCard.Call(ctx, uint64(question.QuestionID))
		if err != nil {
			slog.Error("Unable to call add-card function", "algo", algo.AlgorithmID, "err", err)
			return nil, errors.Join(ErrCantGetAddCard, err)
		}
	}
	for _, grades := range allGradesWithRightTags {
		_, err := gradeCard.Call(ctx, uint64(grades.QuestionID), uint64(grades.Timestamp), uint64(grades.Result))
		if err != nil {
			slog.Error("Unable to call grade-card function", "algo", algo.AlgorithmID, "err", err)
			return nil, errors.Join(ErrCantGetGradeCard, err)
		}
	}
	addr, err := malloc.Call(ctx, 24)
	if err != nil {
		slog.Error("Unable to allocate", "algo", algo.AlgorithmID, "err", err)
		return nil, errors.Join(ErrCantAllocate, err)
	}
	_, err = getCard.Call(ctx, addr[0], uint64(time.Now().Unix()))
	if err != nil {
		slog.Error("Unable to get next card", "algo", algo.AlgorithmID, "err", err)
		return nil, errors.Join(ErrCantCallNextCard, err)
	}

	lenDueCards, inRange := mod.Memory().ReadUint32Le(uint32(addr[0]))
	if !inRange {
		slog.Error("Cannot read 4 bytes in getting amount of due cards", "algo", algo.AlgorithmID)
		return nil, ErrCantReadBytes
	}

	lenNonDueCards, inRange := mod.Memory().ReadUint32Le(uint32(addr[0] + 4))
	if !inRange {
		slog.Error("Cannot read 4 bytes in getting amount of non due cards", "algo", algo.AlgorithmID)
		return nil, ErrCantReadBytes
	}

	lenNewCards, inRange := mod.Memory().ReadUint32Le(uint32(addr[0] + 8))
	if !inRange {
		slog.Error("Cannot read 4 bytes in getting amount of new cards", "algo", algo.AlgorithmID)
		return nil, ErrCantReadBytes
	}
	nextCard, inRange := mod.Memory().ReadUint64Le(uint32(addr[0] + 16))
	if !inRange {
		slog.Error("Cannot read 8 bytes in getting due card", "algo", algo.AlgorithmID)
		return nil, ErrCantReadBytes
	}

	_, err = free.Call(ctx, addr[0])
	if err != nil {
		slog.Error("Cannot call free", "algo", algo.AlgorithmID, "err", err)
		return nil, ErrCantReadBytes
	}
	question := ""
	answer := ""
	found := false
	for _, questionGot := range allQuestions {
		if questionGot.QuestionID == int64(nextCard) {
			question = questionGot.Question
			answer = questionGot.Answer
			found = true
		}
	}
	if !found {
		slog.Error("Cannot find question id", "nextCard", nextCard)
		return nil, ErrCanGetQuestions

	}
	resultMap := map[string]any{"amount-due-cards": lenDueCards, "amount-new-cards": lenNewCards, "amount-non-due-cards": lenNonDueCards, "id": nextCard, "question": question, "answer": answer}
	return resultMap, nil
}
