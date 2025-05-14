package do

import (
	"context"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/client_to_server"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/server_to_client"
	"log/slog"

	"github.com/forgetaboutitapp/forget-about-it/server"
)

func DeleteNewToken(_ context.Context, userid int64, _ Server, _ *client_to_server.DeleteNewToken) *server_to_client.Message {
	slog.Info("Deleting Token of user", "userid", userid)
	slog.Info("Tokens present", "users waiting for a login", server.UsersWaiting)

	delete(server.UsersWaiting, userid)
	return &server_to_client.Message{ReturnMessage: &server_to_client.Message_OkMessage{}}
}
