package secure

import (
	"net/http"

	"github.com/forgetaboutitapp/forget-about-it/server"
	"github.com/google/uuid"
	"golang.org/x/exp/slog"
)

func DeleteNewToken(userid uuid.UUID, s Server, w http.ResponseWriter, r *http.Request) {
	slog.Info("Deleting Token of user", "userid", userid)
	slog.Info("Tokens present", "users waiting for a login", server.UsersWaiting)
	func() {
		server.MutexUsersWaiting.Lock()
		defer server.MutexUsersWaiting.Unlock()

		delete(server.UsersWaiting, userid)

	}()
}
