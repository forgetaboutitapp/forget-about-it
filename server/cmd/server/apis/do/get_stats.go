package do

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/forgetaboutitapp/forget-about-it/server/pkg/sql_queries"
	v1 "github.com/forgetaboutitapp/forget-about-it/server/protobufs-build/client_server/v1"
)

func GetStats(ctx context.Context, user sql_queries.User, _ string, s *Server, req *v1.GetStatsRequest) *v1.GetStatsResponse {
	slog.Info("Getting stats of user", "userid", user.UserID)

	logs, err := s.Db.GetAllGrades(ctx, user.UserID)
	if err != nil {
		slog.Error("Error getting all grades", "err", err)
		return &v1.GetStatsResponse{
			Result: &v1.GetStatsResponse_Error{
				Error: &v1.ErrorMessage{Error: "Internal Server Error"},
			},
		}
	}

	heatmapData := map[uint64]uint32{}
	for _, log := range logs {
		ts := uint64(toStartOfDay(time.Unix(log.Timestamp, 0), int(req.TzOffset)).Unix())
		slog.Info("Adding heatmap data", "timestamp", log.Timestamp, "start_of_day", ts, "tzoffset", req.TzOffset)
		heatmapData[ts] += 1
	}

	// NOTE: futureResults logic from original code seemed unused in the return, keeping it out for now.

	pastResults := map[uint64]float64{}
	pastPositive := map[uint64]int{}
	pastTotal := map[uint64]int{}
	for _, log := range logs {
		ts := uint64(toStartOfDay(time.Unix(log.Timestamp, 0), int(req.TzOffset)).Unix())

		slog.Info("logs", "time", strconv.FormatUint(ts, 10), "res", log.Result, "pp", pastPositive[ts])

		if log.Result == 1 {
			pastPositive[ts]++
		}
		pastTotal[ts]++
	}

	for k, v := range pastTotal {
		pastResults[k] = float64(pastPositive[k]) / float64(v)
	}

	slog.Info("sending data", "past_results_count", len(pastResults), "heatmap_count", len(heatmapData))
	return &v1.GetStatsResponse{
		Result: &v1.GetStatsResponse_Ok{
			Ok: &v1.GetStats{
				PastResults: pastResults,
				PastUsage:   heatmapData,
			},
		},
	}
}

func toStartOfDay(t time.Time, offset int) time.Time {
	location := time.FixedZone("WHATEVER_LOCAL", offset)
	year, month, day := t.In(location).Date()
	dayStartTime := time.Date(year, month, day, 0, 0, 0, 0, location)
	return dayStartTime
}
