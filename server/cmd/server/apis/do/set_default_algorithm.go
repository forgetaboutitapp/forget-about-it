package do

import (
	"context"
	"database/sql"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/client_to_server"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/server_to_client"
	"log/slog"

	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
)

func SetDefaultAlgorithm(ctx context.Context, userid int64, token string, s Server, arg *client_to_server.SetDefaultAlgorithm) *server_to_client.Message {
	if userid == 0 {
		panic("userid is empty")
	}

	slog.Info("userid", "userid", userid)
	err := s.Db.SetDefaultAlgorithm(ctx, sql_queries.SetDefaultAlgorithmParams{UserID: userid, DefaultAlgorithm: sql.NullInt64{
		Valid: true,
		Int64: int64(arg.AlgorithmId),
	}})
	if err != nil {
		slog.Error("can't configure algorithm", "id", int64(arg.AlgorithmId), "err", err)
		return makeError("Internal Server Error")
	}
	return makeOk(&server_to_client.OkMessage{OkMessage: &server_to_client.OkMessage_SetDefaultAlgorithm{SetDefaultAlgorithm: &server_to_client.SetDefaultAlgorithm{}}})

}
