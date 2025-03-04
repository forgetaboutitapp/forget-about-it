package secure

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	uuidUtils "github.com/forgetaboutitapp/forget-about-it/server/pkg/uuid_utils"
	"github.com/google/uuid"
)

func GenerateNewToken(userid uuid.UUID, s Server, w http.ResponseWriter, r *http.Request) {
	newUUID := uuid.New()
	slog.Info("About to get new mnemonic")
	mnenmonic, err := uuidUtils.NewMnemonicFromUuid(newUUID)
	if err != nil {
		slog.Error("can't get mnemonic from uuid", "uuid", newUUID, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	slog.Info("About to write to channel")

	slog.Info("About to generate json")
	jsonVal, err := json.Marshal(map[string]any{"type": "ok", "newUUID": newUUID.String(), "mnemonic": strings.Split(mnenmonic, " ")})
	if err != nil {
		panic(err)
	}
	slog.Info("about to write to db")
	timeoutContext, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	params := sql_queries.AddLoginParams{
		LoginUuid:         newUUID.String(),
		UserUuid:          userid.String(),
		DeviceDescription: "",
		Created:           time.Now().Unix(),
	}
	err = s.Db.AddLogin(timeoutContext, params)
	if err != nil {
		slog.Error("cannot create new login", "params", params, "err", err)
		jsonVal, err := json.Marshal(map[string]string{"type": "error", "message": "Internal Server Error"})
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(w, "%s\n\n", jsonVal)
		return
	}
	slog.Info("about to write to client")

	_, err = fmt.Fprintf(w, "%s\n\n", jsonVal)

	if err != nil {
		slog.Error("write to client", "newUUID", newUUID, "mnemonic", mnenmonic, "err", err)
		return
	}
	slog.Info("About to done")
	func() {
		server.MutexUsersWaiting.Lock()
		defer server.MutexUsersWaiting.Unlock()
		server.UsersWaiting[userid] = struct{}{}
	}()
}
