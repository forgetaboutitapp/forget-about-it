package secure

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
)

type RemoteDevice struct {
	LastUsed  *int64 `json:"last-used,omitempty"`
	Title     string `json:"title"`
	DateAdded int64  `json:"date-added"`
}

type RemoteSettings struct {
	RemoteDevices []RemoteDevice `json:"remote-devices,omitempty"`
}

func GetRemoteSettings(userid int64, s Server, w http.ResponseWriter, r *http.Request) {
	if userid == 0 {
		panic("userid is empty")
	}
	timeoutContext, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()
	slog.Info("userid", "userid", userid)
	rows, err := func() ([]sql_queries.FindLoginIDByUserRow, error) {
		return s.Db.FindLoginIDByUser(timeoutContext, userid)
	}()
	if err != nil {
		slog.Error("can't find login by userid", "uuid", userid, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	settings := RemoteSettings{}
	for _, row := range rows {
		var lastUsedString *int64 = nil
		if val, ok := row.Lastused.(int64); ok {
			lastUsedString = &val
		}
		settings.RemoteDevices = append(settings.RemoteDevices, RemoteDevice{Title: row.DeviceDescription, LastUsed: lastUsedString, DateAdded: row.Created})
	}
	jsonVal, err := json.Marshal(settings)
	if err != nil {
		slog.Error("can't find logi by userid", "uuid", userid, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	slog.Info("Returning json", "json", string(jsonVal))
	w.Write(jsonVal)

}
