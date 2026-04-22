package do

import (
	"context"
	"log/slog"
	"slices"
	"strconv"

	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	v1 "github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/client_server/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func GetRemoteSettings(ctx context.Context, user sql_queries.User, token string, s *Server, _ *v1.GetRemoteSettingsRequest) *v1.GetRemoteSettingsResponse {
	slog.Info("userid", "userid", user.UserID)
	userSettingsRows, err := s.Db.FindLoginIDByUser(ctx, user.UserID)

	if err != nil || len(userSettingsRows) == 0 {
		slog.Error("can't find login by userid", "userid", user.UserID, "err", err)
		return &v1.GetRemoteSettingsResponse{
			Result: &v1.GetRemoteSettingsResponse_Error{
				Error: &v1.ErrorMessage{Error: "Cannot find user settings"},
			},
		}
	}
	algorithmsRows, err := s.Db.GetSpacingAlgorithms(ctx)
	if err != nil {
		slog.Error("can't get algorithms", "err", err)
		return &v1.GetRemoteSettingsResponse{
			Result: &v1.GetRemoteSettingsResponse_Error{
				Error: &v1.ErrorMessage{Error: "Cannot get algorithms"},
			},
		}
	}
	defaultAlgorithm, err := s.Db.GetDefaultAlgorithm(ctx, user.UserID)
	if err != nil {
		slog.Error("can't get default algorithm", "err", err)
		return &v1.GetRemoteSettingsResponse{
			Result: &v1.GetRemoteSettingsResponse_Error{
				Error: &v1.ErrorMessage{Error: "Cannot get default algorithms"},
			},
		}
	}

	settings := v1.GetRemoteSettings{}
	indexOfCurrentDevice := -1

	for index, row := range userSettingsRows {
		var lastUsedSeconds int64 = 0
		if val, ok := row.Lastused.(int64); ok {
			lastUsedSeconds = val
		}
		if token == row.LoginUuid {
			indexOfCurrentDevice = index
		}
		slog.Info("getting row settings", "desc", row.DeviceDescription, "id", row.IndexID)
		settings.RemoteDevices = append(settings.RemoteDevices, &v1.Device{
			Title:     row.DeviceDescription,
			LastUsed:  &timestamppb.Timestamp{Seconds: lastUsedSeconds},
			DateAdded: &timestamppb.Timestamp{Seconds: row.Created},
			LoginId:   strconv.Itoa(int(row.IndexID)),
		})
	}
	if indexOfCurrentDevice == -1 {
		slog.Error("row not found", "userSettingsRows", userSettingsRows, "token", token)
		return &v1.GetRemoteSettingsResponse{
			Result: &v1.GetRemoteSettingsResponse_Error{
				Error: &v1.ErrorMessage{Error: "Internal Server Error"},
			},
		}
	}
	settings.RemoteDevices[indexOfCurrentDevice], settings.RemoteDevices[0] = settings.RemoteDevices[0], settings.RemoteDevices[indexOfCurrentDevice]
	for _, row := range algorithmsRows {
		settings.Algorithms = append(settings.Algorithms, &v1.Algorithm{
			AlgorithmId:   uint32(row.AlgorithmID),
			AuthorName:    row.Author,
			License:       row.License,
			RemoteUrl:     row.RemoteUrl,
			DownloadUrl:   row.DownloadUrl,
			DateAdded:     &timestamppb.Timestamp{Seconds: row.Timestamp},
			Version:       uint32(row.Version),
			AlgorithmName: row.AlgorithmName,
		})
	}
	if defaultAlgorithm.Valid {
		foundId := slices.IndexFunc(settings.Algorithms, func(a *v1.Algorithm) bool { return a.AlgorithmId == uint32(defaultAlgorithm.Int64) })
		if foundId != -1 {
			settings.Algorithms[0], settings.Algorithms[foundId] = settings.Algorithms[foundId], settings.Algorithms[0]
		}
		settings.DefaultAlgorithm = uint32(defaultAlgorithm.Int64)
	}
	return &v1.GetRemoteSettingsResponse{
		Result: &v1.GetRemoteSettingsResponse_Ok{
			Ok: &settings,
		},
	}
}
