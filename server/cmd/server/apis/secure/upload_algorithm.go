package secure

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"log/slog"
	"math/rand/v2"
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
	questions := []sql_queries.GetAllQuestionsRow{}
	grades := []sql_queries.QuestionsLog{}
	_, err = runAlgorithm(ctx, algorithm, questions, grades)
	if err != nil {
		slog.Error("cannot run wasm", "data", algorithm.AlgorithmName, "err", err)
		return map[string]any{"error": "Unable to run wasm"}, nil
	}
	algos, err := s.Db.GetSpacingAlgorithms(ctx)
	if err != nil {
		slog.Error("cannot get algos", "err", err)
		return nil, errors.Join(err, ErrCantGetSpacingAlgo)
	}
	id := rand.Int32()
	tx, err := s.OrigDB.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("cannot start transaction", "err", err)
		return nil, err
	}
	defer tx.Rollback()
	qtx := s.Db.WithTx(tx)
	params := sql_queries.AddSpacingAlgorithmParams{
		AlgorithmID:   id,
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
	if len(algos) == 0 {
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
