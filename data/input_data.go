package data

import (
	"fmt"
	"github.com/shopspring/decimal"
	"math"
)

// InputData It is using decimal.Decimal, to deal with property of Money, probably we should rely on own implementation
// For example, there is no need to deal with more than two decimal places for USD or RUB, but we must deal with more for BTC for example
type InputData struct {
	DecimalCost decimal.Decimal
	Timestamp uint64
}

func IsThereAreAnyDataToProcess(TimeCurrent uint64, WindowSeconds uint64, data []*InputData) bool {
	for _, log := range data {
		// We are skipping data not in window
		if !log.IsInWindow(TimeCurrent, WindowSeconds) {
			continue
		}

		return true
	}

	return false
}

// IsInWindow Conversion to float64 is crap, but will do, probably TODO something better
func (d *InputData) IsInWindow(TimeCurrent uint64, WindowSeconds uint64) bool {
	return math.Abs(float64(d.Timestamp) - float64(TimeCurrent)) <= float64(WindowSeconds)
}

func (d *InputData) IsBeforeWindow(TimeCurrent uint64, WindowSeconds uint64) bool {
	return float64(d.Timestamp) - float64(TimeCurrent) + float64(WindowSeconds) < 0
}

func (d *InputData) String() string {
	return fmt.Sprintf("%d: %s", d.Timestamp, d.DecimalCost.String())
}