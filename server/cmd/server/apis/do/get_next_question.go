package do

import (
	"context"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/client_to_server"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/server_to_client"
	"log/slog"

	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
)

func GetNextQuestion(ctx context.Context, userid int64, s Server, arg *client_to_server.GetNextQuestion) *server_to_client.Message {
	slog.Info("Getting next question")
	tagsToAsk := arg.Tags
	defaultAlgo, err := s.Db.GetDefaultAlgorithm(ctx, userid)

	if err != nil {
		slog.Error("can't get default algo", "uuid", userid, "err", err)
		return makeError("Can't get the default spacing algorithm")
	}
	algos, err := s.Db.GetSpacingAlgorithms(ctx)

	if err != nil {
		slog.Error("can't get spacing algorithm", "uuid", userid, "err", err)
		return makeError("Can't get the default spacing algorithm")

	}
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

	slog.Info("Running algo", "name", algo.AlgorithmName)
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
		tagsToAsk:      tagsToAsk,
	}
	slog.Info("Running algo", "allGrades", allGrades, "tagsToAsk", tagsToAsk, "tagsByQuestion", tagsByQuestion, "getNewQuestion", arg.GetNewQuestion)
	ret, err, displayError := runAlgorithm(ctx, algoArgs, arg.GetNewQuestion)
	if displayError != "" {
		slog.Error("error message from scheduler", "displayError", displayError)
		return makeError(displayError)
	}
	if err != nil {
		slog.Error("cannot run wasm", "algoname", algo.AlgorithmName, "err", err)
		return makeError("Unable to run wasm")
	}
	question := ""
	answer := ""
	explanation := ""
	memoHint := ""
	found := false
	var nextCardId uint32
	var toq server_to_client.GetNextQuestion_TypeOfQuestion
	if c := ret.GetNewCard(); c != nil {
		nextCardId = c.Id
		toq = server_to_client.GetNextQuestion_TYPE_OF_QUESTION_NEW
	} else if c := ret.GetDueCard(); c != nil {
		nextCardId = c.Id
		toq = server_to_client.GetNextQuestion_TYPE_OF_QUESTION_DUE
	} else if c := ret.GetNonDueCard(); c != nil {
		nextCardId = c.Id
		toq = server_to_client.GetNextQuestion_TYPE_OF_QUESTION_NON_DUE
	} else if c := ret.GetNoCard(); c != nil {
		slog.Error("There should be at least one card", "ret", ret)
		return makeError("Internal Server Error")
	} else {
		slog.Error("wrong type", "ret", ret)
		return makeError("Internal Server Error")
	}
	for _, questionGot := range allQuestions {
		if questionGot.QuestionID == int64(nextCardId) {
			question = questionGot.Question
			answer = questionGot.Answer
			memoHint = questionGot.MemoHint
			explanation = questionGot.Explanation
			found = true
		}
	}
	if !found {
		slog.Error("Cannot find question id", "nextCard", ret)
		return makeError("Cannot find question id")

	}

	return &server_to_client.Message{
		ReturnMessage: &server_to_client.Message_OkMessage{
			OkMessage: &server_to_client.OkMessage{
				OkMessage: &server_to_client.OkMessage_GetNextQuestion{
					GetNextQuestion: &server_to_client.GetNextQuestion{
						Flashcard: &server_to_client.Flashcard{
							Id:          nextCardId,
							Question:    question,
							Answer:      answer,
							MemoHint:    memoHint,
							Explanation: explanation,
						},
						NewQuestions:    uint32(ret.AmntNewCards),
						DueQuestions:    uint32(ret.AmntDueCards),
						NonDueQuestions: uint32(ret.AmntNonDueCards),
						TypeOfQuestion:  toq,
					},
				},
			},
		},
	}
}
