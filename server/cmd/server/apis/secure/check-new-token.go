package secure

import (
	"context"
	"log/slog"

	"github.com/forgetaboutitapp/forget-about-it/server"
)

func CheckNewToken(ctx context.Context, userid int64, s Server, m map[string]any) (map[string]any, error) {
	server.MutexUsersWaiting.Lock()
	defer server.MutexUsersWaiting.Unlock()
	_, found := server.UsersWaiting[userid]

	slog.Info("token found", "found", found, "UsersWaiting", server.UsersWaiting)
	if !found {
		return map[string]any{"result": "done"}, nil
	} else {
		return map[string]any{"result": "waiting"}, nil
	}
}
