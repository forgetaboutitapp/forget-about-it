package secure

import (
	"context"
	"log/slog"

	"github.com/forgetaboutitapp/forget-about-it/server"
)

func DeleteNewToken(ctx context.Context, userid int64, s Server, m map[string]any) (map[string]any, error) {
	slog.Info("Deleting Token of user", "userid", userid)
	slog.Info("Tokens present", "users waiting for a login", server.UsersWaiting)
	func() {
		server.MutexUsersWaiting.Lock()
		defer server.MutexUsersWaiting.Unlock()

		delete(server.UsersWaiting, userid)

	}()
	return nil, nil
}
