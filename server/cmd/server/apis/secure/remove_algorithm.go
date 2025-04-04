package secure

import (
	"context"
	"log/slog"
)

func RemoveAlgorithm(ctx context.Context, userid int64, s Server, args map[string]any) (map[string]any, error) {
	slog.Info("Deleting Login", "userid", userid)
	algoId := args["algorithm-name"].(string)
	err := s.Db.DeleteAlgorithmByName(ctx, algoId)
	if err != nil {
		slog.Error("cannot delete Algorithm", "id", algoId, "err", err)
		return map[string]any{"error": "Unable to delete algorithm"}, nil
	}
	return nil, nil
}
