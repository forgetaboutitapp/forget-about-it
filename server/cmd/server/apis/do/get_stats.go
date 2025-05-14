package do

import (
	"context"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/client_to_server"
	"github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/protobufs/server_to_client"
	"log/slog"
	"strconv"
	"time"
)

func GetStats(ctx context.Context, userid int64, token string, s Server, arg *client_to_server.GetStats) *server_to_client.Message {
	slog.Info("Getting stats of user", "userid", userid)

	logs, err := s.Db.GetAllGrades(ctx, userid)
	if err != nil {
		slog.Error("Error getting all grades", "err", err)
		return makeError("Internal Server Error")
	}

	heatmapData := map[uint64]uint32{}
	for _, log := range logs {
		ts := uint64(toStartOfDay(time.Unix(log.Timestamp, 0), int(arg.TzOffset)).Unix())
		slog.Info("Adding heatmap data", "timestamp", log.Timestamp-7*60*60, "timestamp", ts, "tzoffset", arg)
		heatmapData[ts] += 1
	}

	futureResults := map[string]int{}
	for i := range 100 {
		futureResults[strconv.FormatInt(time.Now().Unix()+int64(i*60*60*24), 10)] = 100 - i
	}

	pastResults := map[uint64]float64{}
	pastPositive := map[uint64]int{}
	pastTotal := map[uint64]int{}
	for _, log := range logs {
		ts := uint64(toStartOfDay(time.Unix(log.Timestamp, 0), int(arg.TzOffset)).Unix())

		slog.Info("logs", "time", strconv.FormatUint(ts, 10), "res", log.Result, "pp", pastPositive[ts])

		if log.Result == 1 {
			pastPositive[ts]++
		}
		pastTotal[ts]++

	}
	for k, v := range pastTotal {
		pastResults[k] = float64(pastPositive[k]) / float64(v)

	}

	returnVal := map[string]any{"heatmap-data": heatmapData, "past-results": pastResults}
	slog.Info("sending data", "returnVal", returnVal)
	return makeOk(&server_to_client.OkMessage{OkMessage: &server_to_client.OkMessage_GetStats{GetStats: &server_to_client.GetStats{PastResults: pastResults, PastUsage: heatmapData}}})
}

func toStartOfDay(t time.Time, offset int) time.Time {
	location := time.FixedZone("WHATEVER_LOCAL", offset)
	year, month, day := t.In(location).Date()
	dayStartTime := time.Date(year, month, day, 0, 0, 0, 0, location)
	return dayStartTime
}
