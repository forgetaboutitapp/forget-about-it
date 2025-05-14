package do

import (
	"context"
	"database/sql"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/client_to_server"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/server_to_client"
	"log/slog"
	"strconv"
)

func RemoveLogin(ctx context.Context, userid int64, token string, s Server, arg *client_to_server.RemoveLogin) *server_to_client.Message {
	slog.Info("Deleting Login", "userid", userid, "login-id", arg.LoginId)
	loginId := arg.LoginId
	realId, err := strconv.Atoi(loginId)
	if err != nil {
		slog.Error("cannot convert loginid to int", "id", loginId, "err", err)
		return makeError("Question ID is not valid")
	}
	tx, err := s.OrigDB.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("cannot start transaction", "err", err)
		return makeError("Internal Server Error")
	}
	defer func(tx *sql.Tx) {
		tx.Rollback()

	}(tx)
	qtx := s.Db.WithTx(tx)
	loginUUID, err := qtx.GetLoginUUIDFromIndexId(ctx, int64(realId))
	if err != nil {
		slog.Error("cannot get uuid", "id", loginId, "err", err)
		return makeError("Internal Server Error")
	}
	err = qtx.DeleteLoginsFromLogs(ctx, loginUUID.LoginUuid)
	if err != nil {
		slog.Error("cannot delete from logs", "id", loginId, "err", err)
		return makeError("Internal Server Error")
	}
	err = qtx.DeleteLoginByIndexId(ctx, int64(realId))
	if err != nil {
		slog.Error("cannot delete login", "loginUUID", loginUUID, "id", loginId, "realId", realId, "err", err)
		return makeError("Internal Server Error")
	}
	err = tx.Commit()
	if err != nil {
		slog.Error("cannot commit transaction", "err", err)
		return makeError("Internal Server Error")
	}
	return makeOk(&server_to_client.OkMessage{OkMessage: &server_to_client.OkMessage_RemoveLogin{RemoveLogin: &server_to_client.RemoveLogin{}}})
}
