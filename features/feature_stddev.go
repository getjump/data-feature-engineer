package features

import (
	dfedata "data-feature-engineer/data"
	"data-feature-engineer/storage"
	"github.com/shopspring/decimal"
	"math"
)

// StdDevFeature Update looks like AvgFeature.Update, refactor probably?
// Well we mostly Can't reuse AvgFeature result because of the structure of running StdDev algorithm
// StdDevFeature implements Welford's online algorithm for continuous computation of standard deviation
type StdDevFeature struct {
	LastMean decimal.Decimal
	LastS decimal.Decimal

	BasicRunningFeature
}

func (f *StdDevFeature) New(WindowSeconds uint64, dataStorage storage.InputDataStorage) *StdDevFeature {
	// We already have mean in AvgFeature
	f.LastMean = decimal.NewFromInt(0)
	f.LastS = decimal.NewFromInt(0)
	f.WindowSeconds = WindowSeconds
	f.RunningFeature = f
	f.DataStorage = dataStorage
	return f
}

//func (f *StdDevFeature) calculateMeanVarianceS(data *dfedata.InputData, amount uint64) (mean decimal.Decimal, variance decimal.Decimal, S decimal.Decimal) {
//	n := decimal.NewFromInt(int64(amount))
//	PrevMean := f.LastMean
//	deltaPrevMean := data.DecimalCost.Sub(PrevMean)
//	mean = PrevMean.Add(deltaPrevMean.Div(n))
//
//	deltaLastMean := data.DecimalCost.Sub(f.LastMean)
//	S = f.LastS.Add(deltaPrevMean.Mul(deltaLastMean))
//
//	subValue, _ := S.Div(decimal.NewFromInt(int64(amount - 1))).Float64()
//	sqrtValue := math.Sqrt(subValue)
//	variance = decimal.NewFromFloat(sqrtValue)
//	return
//}

func (f *StdDevFeature) InvalidateData(data *dfedata.InputData) {
	//if f.LastAmount - 1 < 2 {
	//	f.LastValue = decimal.NewFromInt(0)
	//	return
	//}
	if f.LastAmount <= 2 {
		f.LastAmount -= 1
		return
	}

	n := decimal.NewFromInt(int64(f.LastAmount))
	PrevMean := f.LastMean
	deltaPrevMean := data.DecimalCost.Sub(PrevMean)
	f.LastMean = PrevMean.Add(deltaPrevMean.Div(n))

	deltaLastMean := data.DecimalCost.Sub(f.LastMean)
	f.LastS = f.LastS.Add(deltaPrevMean.Mul(deltaLastMean))

	subValue, _ := f.LastS.Div(decimal.NewFromInt(int64(f.LastAmount - 1))).Float64()
	sqrtValue := math.Sqrt(subValue)
	f.LastValue = f.LastValue.Sub(decimal.NewFromFloat(sqrtValue))

	f.LastAmount -= 1
}

func (f *StdDevFeature) CalculateData(data *dfedata.InputData) {
	n := decimal.NewFromInt(int64(f.LastAmount + 1))
	if f.LastAmount == 0 {
		f.LastMean = data.DecimalCost
		f.LastS = decimal.NewFromInt(0)
	} else {
		PrevMean := f.LastMean
		deltaPrevMean := data.DecimalCost.Sub(PrevMean)
		f.LastMean = PrevMean.Add(deltaPrevMean.Div(n))

		deltaLastMean := data.DecimalCost.Sub(f.LastMean)
		f.LastS = f.LastS.Add(deltaPrevMean.Mul(deltaLastMean))
	}

	f.LastAmount += 1

	// Implementing square root for Decimal is not an easy task
	// Maybe converting to float taking math.Sqrt and returning LastValue will do
	if f.LastAmount >= 2 {
		subValue, _ := f.LastS.Div(decimal.NewFromInt(int64(f.LastAmount - 1))).Float64()
		sqrtValue := math.Sqrt(subValue)
		f.LastValue = decimal.NewFromFloat(sqrtValue)
	} else {
		f.LastValue = decimal.NewFromInt(0)
	}
}

func (f *StdDevFeature) GetAmount() uint64 {
	return f.LastAmount
}