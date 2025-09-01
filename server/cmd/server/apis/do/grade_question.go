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
	initTime := time.Now().UTC().Unix()
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
		Timestamp:  initTime,
	}
	err := s.Db.GradeQuestion(ctx, params)
	if err != nil {
		slog.Error("unable to save grades to the database", "params", params, "err", err)
		return makeError("Internal Server Error")
	}

	defaultAlgo, err := s.Db.GetDefaultAlgorithm(ctx, userid)

	if err != nil {
		slog.Error("can't get default algo", "uuid", userid, "err", err)
		return makeError("Can't get the default spacing algorithm")
	}
	algos, err := s.Db.GetSpacingAlgorithms(ctx)

	var algo sql_queries.SpacingAlgorithm
	if defaultAlgo.Valid {
		var algosIds []int
		found := false
		for _, curAlgo := range algos {
			algosIds = append(algosIds, int(curAlgo.AlgorithmID))
			if curAlgo.AlgorithmID == defaultAlgo.Int64 {
				algo = curAlgo
				found = true
				break
			}
		}
		if !found {
			slog.Error("default algo does not match a valid algorithm id", "algoId", algosIds, "defaultAlgo", defaultAlgo.Int64)
			return makeError("Invalid default algorithm")
		}
	} else if len(algos) > 0 {
		algo = algos[0]
	} else {
		slog.Error("There are no algorithms available", "algo", defaultAlgo, "algos", algo)
		return makeError("Internal Server Error")
	}

	allGrades, err := s.Db.GetAllGrades(ctx, userid)
	if err != nil {
		slog.Error("can't get all grades", "err", err)
		return makeError("Internal Server Error")
	}
	allQuestions, err := s.Db.GetAllQuestions(ctx, userid)

	if err != nil {
		slog.Error("can't get all questions", "err", err)
		return makeError("Internal Server Error")
	}

	tagsByQuestion := make(map[uint32][]string)
	for _, question := range allQuestions {
		tags, err := s.Db.GetTagsByQuestion(ctx, question.QuestionID)
		if err != nil {
			slog.Error("can't get tag for grade", "questionid", question.QuestionID, "err", err)
			return makeError("Internal Server Error")
		}
		tagsByQuestion[uint32(question.QuestionID)] = tags
	}

	algoArgs := RunAlgorithm{
		algo: AlgorithmStruct{
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
		},
		allGrades:      allGrades,
		tagsByQuestion: tagsByQuestion,
	}
	ret, err, displayError := runAlgorithm(ctx, algoArgs, false)
	if displayError != "" {
		slog.Error("error message from scheduler", "displayError", displayError)
		return makeError(displayError)
	}

	if err != nil {
		slog.Error("cannot run wasm", "algoname", algo.AlgorithmName, "err", err)
		return makeError("Unable to run wasm")
	}
	res := ret.FutureCardHeatmap[arg.Questionid]
	slog.Info("Question availability", "questionid", arg.Questionid, "res", res)

	if int64(res) < initTime {
		slog.Error("due time cannot be in the past", "res", res, "now", initTime)
		return makeError("Internal Server Error")
	}
	return makeOk(&server_to_client.OkMessage{OkMessage: &server_to_client.OkMessage_GradeQuestion{GradeQuestion: &server_to_client.GradeQuestion{
		NextDue: res - uint64(initTime),
	}}})
}
