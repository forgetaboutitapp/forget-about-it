package do

import (
	"context"
	"log/slog"

	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	v1 "github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/client_server/v1"
)

func RemoveAlgorithm(ctx context.Context, user sql_queries.User, _ string, s *Server, req *v1.RemoveAlgorithmRequest) *v1.RemoveAlgorithmResponse {
	slog.Info("Deleting Algorithm", "userid", user.UserID, "algorithm_name", req.AlgorithmId)
	algoId := req.AlgorithmId
	err := s.Db.DeleteAlgorithmByName(ctx, algoId)
	if err != nil {
		slog.Error("cannot delete Algorithm", "id", algoId, "err", err)
		return &v1.RemoveAlgorithmResponse{
			Result: &v1.RemoveAlgorithmResponse_Error{
				Error: &v1.ErrorMessage{Error: "Internal Server Error"},
			},
		}
	}
	return &v1.RemoveAlgorithmResponse{
		Result: &v1.RemoveAlgorithmResponse_Ok{
			Ok: &v1.RemoveAlgorithm{},
		},
	}
}
