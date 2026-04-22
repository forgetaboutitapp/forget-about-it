package do

import (
	"context"
	"log/slog"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	v1 "github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/client_server/v1"
)

func CheckNewToken(_ context.Context, user sql_queries.User, _ *Server, _ *v1.CheckNewTokenRequest) *v1.CheckNewTokenResponse {
	_, found := server.UsersWaiting[user.UserID]

	slog.Info("token found", "found", found, "UsersWaiting", server.UsersWaiting)
	return &v1.CheckNewTokenResponse{
		Result: &v1.CheckNewTokenResponse_Ok{
			Ok: &v1.CheckNewToken{
				Done: !found,
			},
		},
	}
}
