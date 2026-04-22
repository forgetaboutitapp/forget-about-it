package do

import (
	"context"
	"log/slog"
	"math/rand/v2"
	"strings"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	uuidUtils "github.com/forgetaboutitapp/forget-about-it/server/pkg/uuid_utils"
	v1 "github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/client_server/v1"
	"github.com/google/uuid"
)

func GenerateNewToken(ctx context.Context, user sql_queries.User, s *Server, _ *v1.GenerateNewTokenRequest) *v1.GenerateNewTokenResponse {
	newUUID := uuid.New()
	slog.Info("About to get new mnemonic")
	mnenmonic, err := uuidUtils.NewMnemonicFromUuid(newUUID)
	if err != nil {
		slog.Error("can't get mnemonic from uuid", "uuid", newUUID, "err", err)
		return &v1.GenerateNewTokenResponse{
			Result: &v1.GenerateNewTokenResponse_Error{
				Error: &v1.ErrorMessage{Error: "can't generate mnemonic"},
			},
		}
	}

	slog.Info("about to write to db")
	params := sql_queries.AddLoginParams{
		LoginUuid:         newUUID.String(),
		UserID:            user.UserID,
		DeviceDescription: "",
		Created:           time.Now().UTC().Unix(),
		IndexID:           int64(rand.Uint32()),
	}

	err = s.Db.AddLogin(ctx, params)
	if err != nil {
		slog.Error("cannot create new login", "params", params, "err", err)
		return &v1.GenerateNewTokenResponse{
			Result: &v1.GenerateNewTokenResponse_Error{
				Error: &v1.ErrorMessage{Error: "Internal Server Error"},
			},
		}
	}

	slog.Info("About to generate response")
	response := &v1.GenerateNewTokenResponse{
		Result: &v1.GenerateNewTokenResponse_Ok{
			Ok: &v1.GenerateNewToken{
				NewUuid:  newUUID.String(),
				Mnemonic: strings.Split(mnenmonic, " "),
			},
		},
	}

	func() {
		server.MutexUsersWaiting.Lock()
		defer server.MutexUsersWaiting.Unlock()
		server.UsersWaiting[user.UserID] = struct{}{}
	}()
	return response
}
