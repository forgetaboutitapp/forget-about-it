package do

import (
	"context"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/client_to_server"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/server_to_client"
	"log/slog"

	"github.com/forgetaboutitapp/forget-about-it/server"
)

func CheckNewToken(_ context.Context, userid int64, _ Server, _ *client_to_server.CheckNewToken) *server_to_client.Message {
	_, found := server.UsersWaiting[userid]

	slog.Info("token found", "found", found, "UsersWaiting", server.UsersWaiting)
	if !found {
		return &server_to_client.Message{ReturnMessage: &server_to_client.Message_OkMessage{OkMessage: &server_to_client.OkMessage{OkMessage: &server_to_client.OkMessage_CheckNewToken{CheckNewToken: &server_to_client.CheckNewToken{Done: true}}}}}
	} else {
		return &server_to_client.Message{ReturnMessage: &server_to_client.Message_OkMessage{OkMessage: &server_to_client.OkMessage{OkMessage: &server_to_client.OkMessage_CheckNewToken{CheckNewToken: &server_to_client.CheckNewToken{Done: false}}}}}
	}
}
