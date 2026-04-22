package do

import (
	"context"
	"log/slog"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	v1 "github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/client_server/v1"
)

func GradeQuestion(ctx context.Context, user sql_queries.User, _ string, s *Server, req *v1.GradeQuestionRequest) *v1.GradeQuestionResponse {
	slog.Info("Grading question", "req", req)
	initTime := time.Now().UTC().Unix()
	questionId := req.Questionid

	correct := req.Correct
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
		return &v1.GradeQuestionResponse{
			Result: &v1.GradeQuestionResponse_Error{
				Error: &v1.ErrorMessage{Error: "Internal Server Error"},
			},
		}
	}

	defaultAlgo, err := s.Db.GetDefaultAlgorithm(ctx, user.UserID)

	if err != nil {
		slog.Error("can't get default algo", "userid", user.UserID, "err", err)
		return &v1.GradeQuestionResponse{
			Result: &v1.GradeQuestionResponse_Error{
				Error: &v1.ErrorMessage{Error: "Can't get the default spacing algorithm"},
			},
		}
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
			return &v1.GradeQuestionResponse{
				Result: &v1.GradeQuestionResponse_Error{
					Error: &v1.ErrorMessage{Error: "Invalid default algorithm"},
				},
			}
		}
	} else if len(algos) > 0 {
		algo = algos[0]
	} else {
		slog.Error("There are no algorithms available", "userid", user.UserID)
		return &v1.GradeQuestionResponse{
			Result: &v1.GradeQuestionResponse_Error{
				Error: &v1.ErrorMessage{Error: "Internal Server Error"},
			},
		}
	}

	allGrades, err := s.Db.GetAllGrades(ctx, user.UserID)
	if err != nil {
		slog.Error("can't get all grades", "err", err)
		return &v1.GradeQuestionResponse{
			Result: &v1.GradeQuestionResponse_Error{
				Error: &v1.ErrorMessage{Error: "Internal Server Error"},
			},
		}
	}
	allQuestions, err := s.Db.GetAllQuestions(ctx, user.UserID)

	if err != nil {
		slog.Error("can't get all questions", "err", err)
		return &v1.GradeQuestionResponse{
			Result: &v1.GradeQuestionResponse_Error{
				Error: &v1.ErrorMessage{Error: "Internal Server Error"},
			},
		}
	}

	tagsByQuestion := make(map[uint32][]string)
	for _, question := range allQuestions {
		tags, err := s.Db.GetTagsByQuestion(ctx, question.QuestionID)
		if err != nil {
			slog.Error("can't get tag for grade", "questionid", question.QuestionID, "err", err)
			return &v1.GradeQuestionResponse{
				Result: &v1.GradeQuestionResponse_Error{
					Error: &v1.ErrorMessage{Error: "Internal Server Error"},
				},
			}
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
		return &v1.GradeQuestionResponse{
			Result: &v1.GradeQuestionResponse_Error{
				Error: &v1.ErrorMessage{Error: displayError},
			},
		}
	}

	if err != nil {
		slog.Error("cannot run wasm", "algoname", algo.AlgorithmName, "err", err)
		return &v1.GradeQuestionResponse{
			Result: &v1.GradeQuestionResponse_Error{
				Error: &v1.ErrorMessage{Error: "Unable to run wasm"},
			},
		}
	}
	res := ret.FutureCardHeatmap[req.Questionid]
	slog.Info("Question availability", "questionid", req.Questionid, "res", res)

	if int64(res) < initTime {
		slog.Error("due time cannot be in the past", "res", res, "now", initTime)
		return &v1.GradeQuestionResponse{
			Result: &v1.GradeQuestionResponse_Error{
				Error: &v1.ErrorMessage{Error: "Internal Server Error"},
			},
		}
	}
	return &v1.GradeQuestionResponse{
		Result: &v1.GradeQuestionResponse_Ok{
			Ok: &v1.GradeQuestion{
				NextDue: res - uint64(initTime),
			},
		},
	}
}
