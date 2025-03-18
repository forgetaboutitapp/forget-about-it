package secure

import (
	"context"
	"errors"
	"log/slog"

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

var ErrCantFindUser = errors.New("can't find user")

func GetRemoteSettings(ctx context.Context, userid int64, s Server, _ map[string]any) (map[string]any, error) {
	if userid == 0 {
		panic("userid is empty")
	}

	slog.Info("userid", "userid", userid)
	rows, err := func() ([]sql_queries.FindLoginIDByUserRow, error) {
		return s.Db.FindLoginIDByUser(ctx, userid)
	}()
	if err != nil || len(rows) == 0 {
		slog.Error("can't find login by userid", "uuid", userid, "err", err)
		return nil, errors.Join(ErrCantFindUser, err)
	}
	settings := RemoteSettings{}
	for _, row := range rows {
		var lastUsedString *int64 = nil
		if val, ok := row.Lastused.(int64); ok {
			lastUsedString = &val
		}
		settings.RemoteDevices = append(settings.RemoteDevices, RemoteDevice{Title: row.DeviceDescription, LastUsed: lastUsedString, DateAdded: row.Created})
	}
	return map[string]any{"settings": settings}, nil

}
