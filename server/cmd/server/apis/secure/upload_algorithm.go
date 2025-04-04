package secure

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log/slog"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
)

func UploadAlgorithm(ctx context.Context, userid int64, s Server, m map[string]any) (map[string]any, error) {
	data := m["data"].(string)
	var algorithm AlgorithmStruct
	err := json.Unmarshal([]byte(data), &algorithm)
	if err != nil {
		slog.Error("Unable to parse data", "data", data, "err", err)
		return map[string]any{"error": "Unable to parse data"}, nil
	}
	wasmBinary, err := base64.StdEncoding.DecodeString(algorithm.WasmString)
	if err != nil {
		slog.Error("Unable to decode wasm", "data", data, "err", err)
		return map[string]any{"error": "Unable to decode wasm"}, nil
	}
	algorithm.WasmBytes = wasmBinary
	algorithm.Author = strings.TrimSpace(algorithm.Author)
	algorithm.DownloadUrl = strings.TrimSpace(algorithm.DownloadUrl)
	algorithm.License = strings.TrimSpace(algorithm.License)
	algorithm.AlgorithmName = strings.TrimSpace(algorithm.AlgorithmName)
	algorithm.RemoteURL = strings.TrimSpace(algorithm.RemoteURL)

	if algorithm.Alloc == "" {
		slog.Error("algorithm alloc cannot be null")
		return map[string]any{"error": "Algorithm file is not valid - alloc must not be empty"}, nil
	}
	if algorithm.ApiVersion == 0 {
		slog.Error("algorithm ApiVersion cannot be null")
		return map[string]any{"error": "Algorithm file is not valid - Api Version must not be 0"}, nil
	}
	if algorithm.Author == "" {
		slog.Error("algorithm Author cannot be null")
		return map[string]any{"error": "Algorithm file is not valid - Author must not be empty"}, nil
	}
	if algorithm.Dealloc == "" {
		slog.Error("algorithm Dealloc cannot be null")
		return map[string]any{"error": "Algorithm file is not valid - Dealloc must not be empty"}, nil
	}
	if algorithm.DownloadUrl == "" {
		slog.Error("algorithm DownloadUrl cannot be null")
		return map[string]any{"error": "Algorithm file is not valid - DownloadUrl must not be empty"}, nil
	}
	if algorithm.Init == "" {
		slog.Error("algorithm Init cannot be null")
		return map[string]any{"error": "Algorithm file is not valid - Init must not be empty"}, nil
	}
	if algorithm.License == "" {
		slog.Error("algorithm License cannot be null")
		return map[string]any{"error": "Algorithm file is not valid - License must not be empty"}, nil
	}
	if algorithm.ModuleName == "" {
		slog.Error("algorithm ModuleName cannot be null")
		return map[string]any{"error": "Algorithm file is not valid - ModuleName must not be empty"}, nil
	}
	if algorithm.AlgorithmName == "" {
		slog.Error("algorithm AlgorithmName cannot be null")
		return map[string]any{"error": "Algorithm file is not valid - AlgorithmName must not be empty"}, nil
	}
	if algorithm.RemoteURL == "" {
		slog.Error("algorithm RemoteURL cannot be null")
		return map[string]any{"error": "Algorithm file is not valid - RemoteURL must not be empty"}, nil
	}
	if len(algorithm.WasmBytes) == 0 {
		slog.Error("algorithm WasmBytes cannot be null")
		return map[string]any{"error": "Algorithm file is not valid - Wasm must not be empty"}, nil
	}
	if algorithm.Version == 0 {
		slog.Error("algorithm Version cannot be null")
		return map[string]any{"error": "Algorithm file is not valid - Version must not be empty"}, nil
	}

	questions := []sql_queries.GetAllQuestionsRow{}
	grades := []sql_queries.QuestionsLog{}
	_, err = runAlgorithm(ctx, algorithm, questions, grades)
	if err != nil {
		slog.Error("cannot run wasm", "data", algorithm.AlgorithmName, "err", err)
		return map[string]any{"error": "Unable to run wasm"}, nil
	}

	id := rand.Int32()
	slog.Info("about to start tx")
	_, err = s.OrigDB.Exec("PRAGMA defer_foreign_keys = on")
	if err != nil {
		slog.Error("cannot defer foreign keys", "err", err)
		return nil, err
	}
	defer s.OrigDB.Exec("PRAGMA defer_foreign_keys = off")
	tx, err := s.OrigDB.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("cannot start transaction", "err", err)
		return nil, err
	}
	defer tx.Rollback()
	qtx := s.Db.WithTx(tx)
	slog.Info("about to get spacing algo")
	algos, err := qtx.GetSpacingAlgorithms(ctx)
	if err != nil {
		slog.Error("cannot get algos", "err", err)
		return nil, errors.Join(err, ErrCantGetSpacingAlgo)
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
			return nil, err
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
		Timestamp:     time.Now().Unix(),
		Wasm:          algorithm.WasmBytes,
	}
	server.DbLock.Lock()
	defer server.DbLock.Unlock()
	err = qtx.AddSpacingAlgorithm(ctx, params)
	if err != nil {
		slog.Error("Cannot save data to the db", "err", err)
		return map[string]any{"error": "Unable to save to the db"}, nil
	}
	slog.Info("about to set default")
	if isFirstAlgorithm {
		err = qtx.SetDefaultAlgorithm(ctx, sql_queries.SetDefaultAlgorithmParams{DefaultAlgorithm: sql.NullInt64{Valid: true, Int64: int64(id)}})
		if err != nil {
			slog.Error("cannot set default algo", "id", id, "err", err)
			return nil, errors.Join(err, ErrCantGetSpacingAlgo)
		}
	}
	err = tx.Commit()
	if err != nil {
		slog.Error("cannot commit tx", "err", err)
		return nil, err
	}
	return map[string]any{"ok": "ok"}, nil
}
