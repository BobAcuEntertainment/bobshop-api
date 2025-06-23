package coin

import (
	"sort"
	"strconv"
	"time"
)

func parseDuration(timeframe string) time.Duration {
	multiplier := 1
	switch timeframe[len(timeframe)-1] {
	case 's':
		multiplier = 1
	case 'm':
		multiplier = 60
	case 'h':
		multiplier = 60 * 60
	case 'd':
		multiplier = 60 * 60 * 24
	case 'w':
		multiplier = 60 * 60 * 24 * 7
	case 'M':
		multiplier = 60 * 60 * 24 * 30 // Approximate month as 30 days
	case 'y':
		multiplier = 60 * 60 * 24 * 365
	default:
		return 0
	}

	value, _ := strconv.Atoi(timeframe[:len(timeframe)-1])
	return time.Duration(value*multiplier) * time.Second
}

func sortTimeframes(timeframes []string) []string {
	durations := make(map[string]time.Duration)

	for _, tf := range timeframes {
		durations[tf] = parseDuration(tf)
	}

	sort.Slice(timeframes, func(i, j int) bool {
		return durations[timeframes[i]] < durations[timeframes[j]]
	})

	return timeframes
}
