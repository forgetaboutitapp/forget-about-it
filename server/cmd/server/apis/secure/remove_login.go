package secure

import (
	"context"
	"log/slog"
	"strconv"
)

func RemoveLogin(ctx context.Context, userid int64, s Server, args map[string]any) (map[string]any, error) {
	slog.Info("Deleting Login", "userid", userid, "login-id", args["login-id"].(string))
	loginId := args["login-id"].(string)
	realId, err := strconv.Atoi(loginId)
	if err != nil {
		slog.Error("cannot convert loginid to int", "id", loginId, "err", err)
		return nil, err
	}
	tx, err := s.OrigDB.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("cannot start transaction", "err", err)
		return nil, err
	}
	defer tx.Rollback()
	qtx := s.Db.WithTx(tx)
	loginUUID, err := qtx.GetLoginUUIDFromIndexId(ctx, int64(realId))
	if err != nil {
		slog.Error("cannot get uuid", "id", loginId, "err", err)
		return nil, err
	}
	err = qtx.DeleteLoginsFromLogs(ctx, loginUUID.LoginUuid)
	if err != nil {
		slog.Error("cannot delete from logs", "id", loginId, "err", err)
		return nil, err
	}
	err = qtx.DeleteLoginByIndexId(ctx, int64(realId))
	if err != nil {
		slog.Error("cannot delete login", "loginUUID", loginUUID, "id", loginId, "realId", realId, "err", err)
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		slog.Error("cannot commit transaction", "err", err)
		return nil, err
	}
	return nil, nil
}
