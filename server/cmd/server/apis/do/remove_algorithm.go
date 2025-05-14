package do

import (
	"context"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/client_to_server"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/server_to_client"
	"log/slog"
)

func RemoveAlgorithm(ctx context.Context, userid int64, token string, s Server, arg *client_to_server.RemoveAlgorithm) *server_to_client.Message {
	slog.Info("Deleting Login", "userid", userid)
	algoId := arg.AlgorithmId
	err := s.Db.DeleteAlgorithmByName(ctx, algoId)
	if err != nil {
		slog.Error("cannot delete Algorithm", "id", algoId, "err", err)
		return makeError("Internal Server Error")
	}
	return makeOk(&server_to_client.OkMessage{OkMessage: &server_to_client.OkMessage_RemoveAlgorithm{RemoveAlgorithm: &server_to_client.RemoveAlgorithm{}}})
}
