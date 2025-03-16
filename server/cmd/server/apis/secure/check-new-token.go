package secure

import (
	"log/slog"
	"net/http"

	"github.com/forgetaboutitapp/forget-about-it/server"
)

func CheckNewToken(userid int64, s Server, w http.ResponseWriter, r *http.Request) {
	server.MutexUsersWaiting.Lock()
	defer server.MutexUsersWaiting.Unlock()
	_, found := server.UsersWaiting[userid]

	slog.Info("token found", "found", found, "UsersWaiting", server.UsersWaiting)
	if !found {
		w.Write([]byte("done"))
		return
	} else {
		w.Write([]byte("waiting"))
		return
	}
}
