package do

import (
	"context"
	"database/sql"
	"log/slog"
	"strconv"

	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	v1 "github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/client_server/v1"
)

func RemoveLogin(ctx context.Context, user sql_queries.User, _ string, s *Server, req *v1.RemoveLoginRequest) *v1.RemoveLoginResponse {
	slog.Info("Deleting Login", "userid", user.UserID, "login-id", req.LoginId)
	loginId := req.LoginId
	realId, err := strconv.Atoi(loginId)
	if err != nil {
		slog.Error("cannot convert loginid to int", "id", loginId, "err", err)
		return &v1.RemoveLoginResponse{
			Result: &v1.RemoveLoginResponse_Error{
				Error: &v1.ErrorMessage{Error: "Login ID is not valid"},
			},
		}
	}
	tx, err := s.OrigDB.BeginTx(ctx, nil)
	if err != nil {
		slog.Error("cannot start transaction", "err", err)
		return &v1.RemoveLoginResponse{
			Result: &v1.RemoveLoginResponse_Error{
				Error: &v1.ErrorMessage{Error: "Internal Server Error"},
			},
		}
	}
	defer func(tx *sql.Tx) {
		tx.Rollback()
	}(tx)
	qtx := s.Db.WithTx(tx)
	loginUUID, err := qtx.GetLoginUUIDFromIndexId(ctx, int64(realId))
	if err != nil {
		slog.Error("cannot get uuid", "id", loginId, "err", err)
		return &v1.RemoveLoginResponse{
			Result: &v1.RemoveLoginResponse_Error{
				Error: &v1.ErrorMessage{Error: "Internal Server Error"},
			},
		}
	}
	err = qtx.DeleteLoginsFromLogs(ctx, loginUUID.LoginUuid)
	if err != nil {
		slog.Error("cannot delete from logs", "id", loginId, "err", err)
		return &v1.RemoveLoginResponse{
			Result: &v1.RemoveLoginResponse_Error{
				Error: &v1.ErrorMessage{Error: "Internal Server Error"},
			},
		}
	}
	err = qtx.DeleteLoginByIndexId(ctx, int64(realId))
	if err != nil {
		slog.Error("cannot delete login", "loginUUID", loginUUID, "id", loginId, "realId", realId, "err", err)
		return &v1.RemoveLoginResponse{
			Result: &v1.RemoveLoginResponse_Error{
				Error: &v1.ErrorMessage{Error: "Internal Server Error"},
			},
		}
	}
	err = tx.Commit()
	if err != nil {
		slog.Error("cannot commit transaction", "err", err)
		return &v1.RemoveLoginResponse{
			Result: &v1.RemoveLoginResponse_Error{
				Error: &v1.ErrorMessage{Error: "Internal Server Error"},
			},
		}
	}
	return &v1.RemoveLoginResponse{
		Result: &v1.RemoveLoginResponse_Ok{
			Ok: &v1.RemoveLogin{},
		},
	}
}
