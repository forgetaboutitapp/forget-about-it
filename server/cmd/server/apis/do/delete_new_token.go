package do

import (
	"context"
	"log/slog"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	v1 "github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/client_server/v1"
)

func DeleteNewToken(_ context.Context, user sql_queries.User, _ *Server, _ *v1.DeleteNewTokenRequest) *v1.DeleteNewTokenResponse {
	slog.Info("Deleting Token of user", "userid", user.UserID)
	slog.Info("Tokens present", "users waiting for a login", server.UsersWaiting)

	delete(server.UsersWaiting, user.UserID)
	return &v1.DeleteNewTokenResponse{
		Result: &v1.DeleteNewTokenResponse_Ok{
			Ok: &v1.DeleteNewToken{},
		},
	}
}
