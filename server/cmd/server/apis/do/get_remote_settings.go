package do

import (
	"context"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/client_to_server"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/server_to_client"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log/slog"
	"slices"
	"strconv"
)

func GetRemoteSettings(ctx context.Context, userid int64, token string, s Server, arg *client_to_server.GetRemoteSettings) *server_to_client.Message {
	if userid == 0 {
		panic("userid is empty")
	}

	slog.Info("userid", "userid", userid)
	userSettingsRows, err := s.Db.FindLoginIDByUser(ctx, userid)

	if err != nil || len(userSettingsRows) == 0 {
		slog.Error("can't find login by userid", "uuid", userid, "err", err)
		return makeError("Cannot find user settings")
	}
	algorithmsRows, err := s.Db.GetSpacingAlgorithms(ctx)
	if err != nil {
		slog.Error("can't get algorithms", "err", err)
		return makeError("Cannot get algorithms")
	}
	defaultAlgorithm, err := s.Db.GetDefaultAlgorithm(ctx, userid)
	if err != nil {
		slog.Error("can't get default algorithm", "err", err)
		return makeError("Cannot get default algorithms")
	}

	settings := server_to_client.GetRemoteSettings{}
	indexOfCurrentDevice := -1

	for index, row := range userSettingsRows {
		var lastUsedString int64 = 0
		if val, ok := row.Lastused.(int64); ok {
			lastUsedString = val
		}
		if token == row.LoginUuid {
			indexOfCurrentDevice = index
		}
		slog.Info("getting row settings", "desc", row.DeviceDescription, "id", row.IndexID)
		settings.RemoteDevices = append(settings.RemoteDevices, &server_to_client.Device{Title: row.DeviceDescription, LastUsed: &timestamppb.Timestamp{Seconds: lastUsedString}, DateAdded: &timestamppb.Timestamp{Seconds: row.Created}, LoginId: strconv.Itoa(int(row.IndexID))})
	}
	if indexOfCurrentDevice == -1 {
		slog.Error("row not found", "userSettingsRows", userSettingsRows, "token", token)
		return makeError("Internal Server Error")
	}
	settings.RemoteDevices[indexOfCurrentDevice], settings.RemoteDevices[0] = settings.RemoteDevices[0], settings.RemoteDevices[indexOfCurrentDevice]
	for _, row := range algorithmsRows {
		settings.Algorithms = append(settings.Algorithms, &server_to_client.Algorithm{
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
		foundId := slices.IndexFunc(settings.Algorithms, func(a *server_to_client.Algorithm) bool { return a.AlgorithmId == uint32(defaultAlgorithm.Int64) })
		if foundId != -1 {
			settings.Algorithms[0], settings.Algorithms[foundId] = settings.Algorithms[foundId], settings.Algorithms[0]
		}
		settings.DefaultAlgorithm = uint32(defaultAlgorithm.Int64)
	}
	return &server_to_client.Message{ReturnMessage: &server_to_client.Message_OkMessage{OkMessage: &server_to_client.OkMessage{OkMessage: &server_to_client.OkMessage_GetRemoteSettings{GetRemoteSettings: &settings}}}}

}
