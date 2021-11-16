package features

import (
	dfedata "data-feature-engineer/data"
	"data-feature-engineer/storage"
	"github.com/shopspring/decimal"
)

// AvgFeature calculates moving average based on window seconds, it doesn't recalculate whole batch of data everytime,
// It just makes use of Mean math properties.
type AvgFeature struct {
	BasicRunningFeature
}

func (f *AvgFeature) New(WindowSeconds uint64, dataStorage storage.InputDataStorage) *AvgFeature {
	f.DataStorage = dataStorage
	f.LastValue = decimal.NewFromInt(0)
	f.RunningFeature = f
	f.WindowSeconds = WindowSeconds
	return f
}

func (f *AvgFeature) InvalidateData(data *dfedata.InputData) {
	n := decimal.NewFromInt(int64(f.LastAmount))
	f.LastValue = f.LastValue.Mul(n).Sub(data.DecimalCost).Div(n.Sub(decimal.NewFromInt(1)))
	f.LastAmount -= 1
}

func (f *AvgFeature) CalculateData(data *dfedata.InputData) {
	// We're calculating moving average continuously
	averaged := data.DecimalCost.Sub(f.LastValue).Div(decimal.NewFromInt(int64(f.LastAmount + 1)))

	f.LastValue = f.LastValue.Add(averaged)
	f.LastAmount += 1
}