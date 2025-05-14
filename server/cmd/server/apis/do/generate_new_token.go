package do

import (
	"context"
	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	uuidUtils "github.com/forgetaboutitapp/forget-about-it/server/pkg/uuid_utils"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/client_to_server"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/server_to_client"
	"github.com/google/uuid"
	"log/slog"
	"math/rand/v2"
	"strings"
	"time"
)

func GenerateNewToken(ctx context.Context, userid int64, s Server, _ *client_to_server.DeleteNewToken) *server_to_client.Message {
	newUUID := uuid.New()
	slog.Info("About to get new mnemonic")
	mnenmonic, err := uuidUtils.NewMnemonicFromUuid(newUUID)
	if err != nil {
		slog.Error("can't get mnemonic from uuid", "uuid", newUUID, "err", err)
		return makeError("can't generate mnemonic")
	}

	slog.Info("about to write to db")
	params := sql_queries.AddLoginParams{
		LoginUuid:         newUUID.String(),
		UserID:            userid,
		DeviceDescription: "",
		Created:           time.Now().Unix(),
		IndexID:           int64(rand.Uint32()),
	}

	err = s.Db.AddLogin(ctx, params)
	if err != nil {
		slog.Error("cannot create new login", "params", params, "err", err)
		return makeError("Internal Server Error")
	}

	slog.Info("About to generate json")

	returnVal := makeOk(&server_to_client.OkMessage{OkMessage: &server_to_client.OkMessage_GenerateNewToken{GenerateNewToken: &server_to_client.GenerateNewToken{NewUuid: newUUID.String(), Mnemonic: strings.Split(mnenmonic, " ")}}})
	slog.Info("About to done")
	func() {
		server.MutexUsersWaiting.Lock()
		defer server.MutexUsersWaiting.Unlock()
		server.UsersWaiting[userid] = struct{}{}
	}()
	return returnVal
}
