package do

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	v1 "github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/client_server/v1"
)

func SetDefaultAlgorithm(ctx context.Context, user sql_queries.User, _ string, s *Server, req *v1.SetDefaultAlgorithmRequest) *v1.SetDefaultAlgorithmResponse {
	slog.Info("userid", "userid", user.UserID)
	err := s.Db.SetDefaultAlgorithm(ctx, sql_queries.SetDefaultAlgorithmParams{
		UserID: user.UserID,
		DefaultAlgorithm: sql.NullInt64{
			Valid: true,
			Int64: int64(req.AlgorithmId),
		},
	})
	if err != nil {
		slog.Error("can't configure algorithm", "id", int64(req.AlgorithmId), "err", err)
		return &v1.SetDefaultAlgorithmResponse{
			Result: &v1.SetDefaultAlgorithmResponse_Error{
				Error: &v1.ErrorMessage{Error: "Internal Server Error"},
			},
		}
	}
	return &v1.SetDefaultAlgorithmResponse{
		Result: &v1.SetDefaultAlgorithmResponse_Ok{
			Ok: &v1.SetDefaultAlgorithm{},
		},
	}
}
