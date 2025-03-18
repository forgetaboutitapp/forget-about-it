package secure

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	uuidUtils "github.com/forgetaboutitapp/forget-about-it/server/pkg/uuid_utils"
	"github.com/google/uuid"
)

var ErrMnemonicFromUUID = errors.New("can't get mnemonic from uuid")
var ErrNewLogin = errors.New("can't create a new login")

func GenerateNewToken(ctx context.Context, userid int64, s Server, m map[string]any) (map[string]any, error) {
	newUUID := uuid.New()
	slog.Info("About to get new mnemonic")
	mnenmonic, err := uuidUtils.NewMnemonicFromUuid(newUUID)
	if err != nil {
		slog.Error("can't get mnemonic from uuid", "uuid", newUUID, "err", err)
		return nil, errors.Join(ErrMnemonicFromUUID, err)
	}

	slog.Info("about to write to db")
	params := sql_queries.AddLoginParams{
		LoginUuid:         newUUID.String(),
		UserID:            userid,
		DeviceDescription: "",
		Created:           time.Now().Unix(),
	}

	err = func() error {
		server.DbLock.Lock()
		defer server.DbLock.Unlock()
		return s.Db.AddLogin(ctx, params)
	}()
	if err != nil {
		slog.Error("cannot create new login", "params", params, "err", err)
		return nil, errors.Join(ErrNewLogin, err)
	}

	slog.Info("About to generate json")
	jsonVal := map[string]any{"type": "ok", "newUUID": newUUID.String(), "mnemonic": strings.Split(mnenmonic, " ")}

	slog.Info("About to done")
	func() {
		server.MutexUsersWaiting.Lock()
		defer server.MutexUsersWaiting.Unlock()
		server.UsersWaiting[userid] = struct{}{}
	}()
	return jsonVal, nil
}
