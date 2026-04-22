package do

import (
	"context"
	"log/slog"

	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	v1 "github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/client_server/v1"
)

func LogData(_ context.Context, user sql_queries.User, _ *Server, req *v1.LogDataRequest) *v1.LogDataResponse {
	slog.Info("Client log received", "userid", user.UserID, "log", req.Log)
	return &v1.LogDataResponse{
		Result: &v1.LogDataResponse_Ok{
			Ok: &v1.LogData{},
		},
	}
}
