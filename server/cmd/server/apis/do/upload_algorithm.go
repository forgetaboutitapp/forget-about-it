package do

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/client_to_server"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/server_to_client"
	"log/slog"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
)

func UploadAlgorithm(ctx context.Context, userid int64, token string, s Server, arg *client_to_server.UploadAlgorithm) *server_to_client.Message {
	data := arg.Algorithm
	var algorithm AlgorithmStruct
	err := json.Unmarshal([]byte(data), &algorithm)
	if err != nil {
		slog.Error("Unable to parse data", "data", data, "err", err)
		return makeError("Internal Server Error")
	}
	wasmBinary, err := base64.StdEncoding.DecodeString(algorithm.WasmString)
	if err != nil {
		slog.Error("Unable to decode wasm", "data", data, "err", err)
		return makeError("Internal Server Error")
	}
	algorithm.WasmBytes = wasmBinary
	algorithm.Author = strings.TrimSpace(algorithm.Author)
	algorithm.DownloadUrl = strings.TrimSpace(algorithm.DownloadUrl)
	algorithm.License = strings.TrimSpace(algorithm.License)
	algorithm.AlgorithmName = strings.TrimSpace(algorithm.AlgorithmName)
	algorithm.RemoteURL = strings.TrimSpace(algorithm.RemoteURL)

	if algorithm.Alloc == "" {
		slog.Error("algorithm alloc cannot be null")
		return makeError("Algorithm file is not valid - alloc must not be empty")
	}
	if algorithm.ApiVersion == 0 {
		slog.Error("algorithm ApiVersion cannot be null")
		return makeError("Algorithm file is not valid - Api Version must not be 0")
	}
	if algorithm.Author == "" {
		slog.Error("algorithm Author cannot be null")
		return makeError("Algorithm file is not valid - Author must not be empty")
	}
	if algorithm.Dealloc == "" {
		slog.Error("algorithm Dealloc cannot be null")
		return makeError("Algorithm file is not valid - Dealloc must not be empty")
	}
	if algorithm.DownloadUrl == "" {
		slog.Error("algorithm DownloadUrl cannot be null")
		return makeError("Algorithm file is not valid - DownloadUrl must not be empty")
	}
	if algorithm.Init == "" {
		slog.Error("algorithm Init cannot be null")
		return makeError("Algorithm file is not valid - Init must not be empty")
	}
	if algorithm.License == "" {
		slog.Error("algorithm License cannot be null")
		return makeError("Algorithm file is not valid - License must not be empty")
	}
	if algorithm.ModuleName == "" {
		slog.Error("algorithm ModuleName cannot be null")
		return makeError("Algorithm file is not valid - ModuleName must not be empty")
	}
	if algorithm.AlgorithmName == "" {
		slog.Error("algorithm AlgorithmName cannot be null")
		return makeError("Algorithm file is not valid - AlgorithmName must not be empty")
	}
	if algorithm.RemoteURL == "" {
		slog.Error("algorithm RemoteURL cannot be null")
		return makeError("Algorithm file is not valid - RemoteURL must not be empty")
	}
	if len(algorithm.WasmBytes) == 0 {
		slog.Error("algorithm WasmBytes cannot be null")
		return makeError("Algorithm file is not valid - Wasm must not be empty")
	}
	if algorithm.Version == 0 {
		slog.Error("algorithm Version cannot be null")
		return makeError("Algorithm file is not valid - Version must not be empty")
	}

	var grades []sql_queries.QuestionsLog
	runArg := RunAlgorithm{
		algo:           algorithm,
		allGrades:      grades,
		tagsByQuestion: map[uint32][]string{},
		tagsToAsk:      []string{},
	}
	_, err, displayError := runAlgorithm(ctx, runArg, false)
	if displayError != "" {
		slog.Error("cannot run wasm", "name", algorithm.AlgorithmName, "err", displayError)
		return makeError("cannot run wasm: " + displayError)
	}
	if err != nil {
		slog.Error("cannot run wasm", "data", algorithm.AlgorithmName, "err", err)
		return makeError("Unable to run wasm")
	}

	id := rand.Int32()
	slog.Info("about to start tx")
	_, err = s.OrigDB.Exec("PRAGMA defer_foreign_keys = on")
	if err != nil {
		slog.Error("cannot defer foreign keys", "err", err)
		return makeError("Internal Server Error")
	}
	defer s.OrigDB.Exec("PRAGMA defer_foreign_keys = off")
	tx, err := s.OrigDB.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("cannot start transaction", "err", err)
		return makeError("Internal Server Error")
	}
	defer func(tx *sql.Tx) {
		tx.Rollback()

	}(tx)
	qtx := s.Db.WithTx(tx)
	slog.Info("about to get spacing algo")
	algos, err := qtx.GetSpacingAlgorithms(ctx)
	if err != nil {
		slog.Error("cannot get algos", "err", err)
		return makeError("Internal Server Error")
	}
	isFirstAlgorithm := len(algos) > 0
	newId := int32(-1)
	slog.Info("bf is first algo", "isFirstAlgorithm", isFirstAlgorithm)
	for _, algo := range algos {
		slog.Info("compare", "algo", algo.AlgorithmName, "algorithm", algorithm.AlgorithmName)
		if algo.AlgorithmName == algorithm.AlgorithmName {
			newId = int32(algo.AlgorithmID)
		}
	}

	if newId != -1 {
		slog.Info("about to delete algo")
		err = qtx.DeleteAlgorithmByName(ctx, algorithm.AlgorithmName)
		if err != nil {
			slog.Error("cannot delete algorithm name", "err", err)
			return makeError("Internal Server Error")
		}
		id = newId
	}
	slog.Info("about to add algo")
	params := sql_queries.AddSpacingAlgorithmParams{
		AlgorithmID:   int64(id),
		Alloc:         algorithm.Alloc,
		ApiVersion:    int64(algorithm.ApiVersion),
		Author:        algorithm.Author,
		Dealloc:       algorithm.Dealloc,
		Desc:          sql.NullString{Valid: true, String: algorithm.Desc},
		DownloadUrl:   algorithm.DownloadUrl,
		Init:          algorithm.Init,
		License:       algorithm.License,
		ModuleName:    algorithm.ModuleName,
		AlgorithmName: algorithm.AlgorithmName,
		RemoteUrl:     algorithm.RemoteURL,
		Version:       int64(algorithm.Version),
		Timestamp:     time.Now().UTC().Unix(),
		Wasm:          algorithm.WasmBytes,
	}
	server.DbLock.Lock()
	defer server.DbLock.Unlock()
	err = qtx.AddSpacingAlgorithm(ctx, params)
	if err != nil {
		slog.Error("Cannot save data to the db", "err", err)
		return makeError("Internal Server Error")
	}
	slog.Info("about to set default")
	if isFirstAlgorithm {
		err = qtx.SetDefaultAlgorithm(ctx, sql_queries.SetDefaultAlgorithmParams{DefaultAlgorithm: sql.NullInt64{Valid: true, Int64: int64(id)}})
		if err != nil {
			slog.Error("cannot set default algo", "id", id, "err", err)
			return makeError("Internal Server Error")
		}
	}
	err = tx.Commit()
	if err != nil {
		slog.Error("cannot commit tx", "err", err)
		return makeError("Internal Server Error")
	}
	return makeOk(&server_to_client.OkMessage{OkMessage: &server_to_client.OkMessage_GradeQuestion{GradeQuestion: &server_to_client.GradeQuestion{}}})
}
