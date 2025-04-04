package secure

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
)

type RemoteDevice struct {
	LastUsed  *int64 `json:"last-used,omitempty"`
	Title     string `json:"title"`
	DateAdded int64  `json:"date-added"`
}

type RemoteAlgorithm struct {
	AlgorithmID    uint32 `json:"id,omitempty"`
	AuthorName     string `json:"author-name,omitempty"`
	License        string `json:"license,omitempty"`
	RemoteURL      string `json:"remote-url,omitempty"`
	DownloadURL    string `json:"download-url,omitempty"`
	TimestampAdded string `json:"time-added,omitempty"`
	Version        int64  `json:"version,omitempty"`
	AlgorithmName  string `json:"algorithm-name,omitempty"`
}

type RemoteSettings struct {
	RemoteDevices    []RemoteDevice    `json:"remote-devices,omitempty"`
	RemoteAlgorithms []RemoteAlgorithm `json:"remote-algorithms,omitempty"`
}

var ErrCantFindUser = errors.New("can't find user")

func GetRemoteSettings(ctx context.Context, userid int64, s Server, _ map[string]any) (map[string]any, error) {
	if userid == 0 {
		panic("userid is empty")
	}

	slog.Info("userid", "userid", userid)
	userSettingsRows, err := func() ([]sql_queries.FindLoginIDByUserRow, error) {
		return s.Db.FindLoginIDByUser(ctx, userid)
	}()

	if err != nil || len(userSettingsRows) == 0 {
		slog.Error("can't find login by userid", "uuid", userid, "err", err)
		return nil, errors.Join(ErrCantFindUser, err)
	}
	algorithmsRows, err := func() ([]sql_queries.SpacingAlgorithm, error) {
		return s.Db.GetSpacingAlgorithms(ctx)
	}()
	if err != nil {
		slog.Error("can't get algorithms", "err", err)
		return nil, errors.Join(ErrCantGetSpacingAlgo, err)
	}
	defaultAlgorithm, err := s.Db.GetDefaultAlgorithm(ctx, userid)
	if err != nil {
		slog.Error("can't get default algorithm", "err", err)
		return nil, errors.Join(ErrCantGetSpacingAlgo, err)
	}

	settings := RemoteSettings{}
	for _, row := range userSettingsRows {
		var lastUsedString *int64 = nil
		if val, ok := row.Lastused.(int64); ok {
			lastUsedString = &val
		}
		settings.RemoteDevices = append(settings.RemoteDevices, RemoteDevice{Title: row.DeviceDescription, LastUsed: lastUsedString, DateAdded: row.Created})
	}
	for _, row := range algorithmsRows {
		settings.RemoteAlgorithms = append(settings.RemoteAlgorithms, RemoteAlgorithm{
			AlgorithmID:    uint32(row.AlgorithmID),
			AuthorName:     row.Author,
			License:        row.License,
			RemoteURL:      row.RemoteUrl,
			DownloadURL:    row.DownloadUrl,
			TimestampAdded: time.Unix(row.Timestamp, 0).Format(time.RFC1123),
			Version:        row.Version,
			AlgorithmName:  row.AlgorithmName,
		})
	}
	if defaultAlgorithm.Valid {
		foundId := slices.IndexFunc(settings.RemoteAlgorithms, func(a RemoteAlgorithm) bool { return a.AlgorithmID == uint32(defaultAlgorithm.Int64) })
		if foundId != -1 {
			settings.RemoteAlgorithms[0], settings.RemoteAlgorithms[foundId] = settings.RemoteAlgorithms[foundId], settings.RemoteAlgorithms[0]
		}
	}
	return map[string]any{"settings": settings}, nil

}
