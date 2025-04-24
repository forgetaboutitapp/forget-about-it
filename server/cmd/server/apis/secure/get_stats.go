package secure

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"time"
)

func GetStats(ctx context.Context, userid int64, s Server, data map[string]any) (map[string]any, error) {
	slog.Info("Getting stats of user", "userid", userid)

	logs, err := s.Db.GetAllGrades(ctx, userid)
	if err != nil {
		slog.Error("Error getting all grades", "err", err)
		return nil, errors.Join(err, ErrCantGetAllGrades)
	}

	heatmapData := map[string]int{}
	for _, log := range logs {
		heatmapData[strconv.FormatInt((log.Timestamp/24/60/60)*(24*60*60), 10)] += 1
	}

	futureResults := map[string]int{}
	for i := range 100 {
		futureResults[strconv.FormatInt(time.Now().Unix()+int64(i*60*60*24), 10)] = 100 - i
	}

	pastResults := map[string]int{}
	pastPositive := map[string]int{}
	pastTotal := map[string]int{}
	for _, log := range logs {
		slog.Info("logs", "time", strconv.FormatInt((log.Timestamp/24/60/60)*(24*60*60), 10), "res", log.Result, "pp", pastPositive[strconv.FormatInt((log.Timestamp/24/60/60)*(24*60*60), 10)])
		if log.Result == 1 {
			pastPositive[strconv.FormatInt((log.Timestamp/24/60/60)*(24*60*60), 10)]++
		}
		pastTotal[strconv.FormatInt((log.Timestamp/24/60/60)*(24*60*60), 10)]++

	}
	for k, v := range pastTotal {
		pastResults[k] = int(float64(pastPositive[k]) / float64(v) * 100)

	}

	returnVal := map[string]any{"heatmap-data": heatmapData, "past-results": pastResults}
	slog.Info("sending data", "returnVal", returnVal)
	return returnVal, nil
}
