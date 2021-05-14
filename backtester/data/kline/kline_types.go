package kline

import (
	"errors"
	"time"

	"github.com/openware/gocryptotrader/backtester/data"
	gctkline "github.com/openware/pkg/kline"
)

var errNoCandleData = errors.New("no candle data provided")

// DataFromKline is a struct which implements the data.Streamer interface
// It holds candle data for a specified range with helper functions
type DataFromKline struct {
	Item gctkline.Item
	data.Base
	Range gctkline.IntervalRangeHolder

	addedTimes map[time.Time]bool
}
