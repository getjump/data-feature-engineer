package features

import (
	dfedata "data-feature-engineer/data"
	"github.com/gammazero/deque"
	"github.com/shopspring/decimal"
	"math"
)

type DecimalMinMaxComparable interface {
	Compare(decimal1 decimal.Decimal, decimal2 decimal.Decimal) bool
}

type BasicMinMaxFeature struct {
	// This deck implementation is mutable which is not bad here
	dq deque.Deque
	Comparator DecimalMinMaxComparable

	BasicFeature
}

// Update To Efficiently implement min/max for moving window we can invalidate data based on a timestamp
// We can memoize timestamp of maximal element if it less than invalidation timestamp, we should recalculate maximal on
// data inside our window? Can we do that more efficiently? Sliding Window Minimum Algorithm with Deque?
// Do we need DataStorage then? Or probably we can implement Deque for that use case?
// TODO We can further optimize our several window seconds implementation by utilizing already computed data like:
// Feature_window_3600(Feature_window_1800(... etc)) do we need to...?
// So we compute from lowest to widest then we just skip already computed fragment and take its value as LastValue from lower
// Also deque for both min and max?
func (f *BasicMinMaxFeature) Update(TimeCurrent uint64, data []*dfedata.InputData, connectionChannel ...chan ConnectionChannelData)  {
	// We can escape there since no updates will be made to LastValue deque will be cleared when there will be data available
	// Minimal value will stay the same as in the last period
	if !dfedata.IsThereAreAnyDataToProcess(TimeCurrent, f.WindowSeconds, data) {
		return
	}

	for f.dq.Len() > 0 && int64(f.dq.Front().(*dfedata.InputData).Timestamp) <= int64(TimeCurrent) - int64(f.WindowSeconds) {
		f.dq.PopFront()
	}

	for _, log := range data {
		if !log.IsInWindow(TimeCurrent, f.WindowSeconds) {
			continue
		}

		// Here we do basic part of sliding maximum window algorithm
		for f.dq.Len() > 0 && f.Comparator.Compare(f.dq.Back().(*dfedata.InputData).DecimalCost, log.DecimalCost) {
			f.dq.PopBack()
		}

		f.dq.PushBack(log)
	}

	f.LastValue = f.dq.Front().(*dfedata.InputData).DecimalCost
	f.OnUpdated(TimeCurrent, data)
}

type MaxFeature struct {
	BasicMinMaxFeature
}

func (f *MaxFeature) New(WindowSeconds uint64) *MaxFeature {
	f.LastValue = decimal.NewFromInt(0)
	f.WindowSeconds = WindowSeconds
	f.Comparator = f
	return f
}

func (f *MaxFeature) Compare(decimal1 decimal.Decimal, decimal2 decimal.Decimal) bool {
	return decimal1.LessThanOrEqual(decimal2)
}

type MinFeature struct {
	BasicMinMaxFeature
}

func (f *MinFeature) New(WindowSeconds uint64) *MinFeature {
	f.LastValue = decimal.NewFromInt(math.MaxInt)
	f.WindowSeconds = WindowSeconds
	f.Comparator = f
	return f
}

func (f *MinFeature) Compare(decimal1 decimal.Decimal, decimal2 decimal.Decimal) bool {
	return decimal1.GreaterThanOrEqual(decimal2)
}