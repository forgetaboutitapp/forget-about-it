package secure

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
)

var ErrCantSetAlgorithm = errors.New("can' set default algorithm")

func SetDefaultAlgorithm(ctx context.Context, userid int64, s Server, idMap map[string]any) (map[string]any, error) {
	if userid == 0 {
		panic("userid is empty")
	}

	slog.Info("userid", "userid", userid)
	err := s.Db.SetDefaultAlgorithm(ctx, sql_queries.SetDefaultAlgorithmParams{UserID: userid, DefaultAlgorithm: sql.NullInt64{
		Valid: true,
		Int64: int64(idMap["algorithm-id"].(float64)),
	}})
	if err != nil {
		slog.Error("can't configure algorithm", "id", int64(idMap["algorithm-id"].(float64)), "err", err)
		return nil, errors.Join(ErrCantSetAlgorithm, err)
	}
	return map[string]any{}, nil

}
