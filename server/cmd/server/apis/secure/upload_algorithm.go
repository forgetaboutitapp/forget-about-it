package secure

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
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
	params := sql_queries.AddSpacingAlgorithmParams{
		AlgorithmID:   rand.Int64(),
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
	err = s.Db.AddSpacingAlgorithm(ctx, params)
	if err != nil {
		slog.Error("Cannot save data to the db", "err", err)
		return map[string]any{"error": "Unable to save to the db"}, nil
	}
	return map[string]any{"ok": "ok"}, nil
}
