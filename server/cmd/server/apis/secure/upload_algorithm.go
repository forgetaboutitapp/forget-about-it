package secure

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log/slog"

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
	return map[string]any{"ok": "ok"}, nil
}
