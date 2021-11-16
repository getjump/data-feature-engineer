package features

import (
	dfedata "data-feature-engineer/data"
	"data-feature-engineer/storage"
	"github.com/shopspring/decimal"
)

// RunningFeatureInterface maybe make it immutable?
type RunningFeatureInterface interface {
	InvalidateData(data *dfedata.InputData)
	CalculateData(data *dfedata.InputData)
}

type BasicRunningFeature struct {
	LastAmount uint64

	DataStorage storage.InputDataStorage
	RunningFeature RunningFeatureInterface

	BasicFeature
}

// Update Data should go in sorted manner, probably linked list is an efficient underlying data storage for this use case
// Probably much code can be refactored for reuse in another features ¯\_(ツ)_/¯ (we do here)
func (f *BasicRunningFeature) Update(TimeCurrent uint64, data []*dfedata.InputData, connectionChannel ...chan ConnectionChannelData)  {
	// Deal with reallocation? Set to len of data, or precompute valid batch size
	var dataToAppend []*dfedata.InputData
	willAppend := dfedata.IsThereAreAnyDataToProcess(TimeCurrent, f.WindowSeconds, data)

	storageData := f.DataStorage.Iterate()

	// We are invalidating old results based on a current window size
	// \frac{1}{N} \Sum_{i}^{N} a_i - a_j = (\Sum_{i}^{N} a_i - a_j) 1/(N-1))
	for _, log := range storageData {
		// =
		//<[>....] <- our time window
		//        ^
		//        |TimeCurrent
		if log.IsInWindow(TimeCurrent, f.WindowSeconds) {
			continue
		}

		// We are invalidating data until there are no more than 1 element available, which we should preserve by TZ
		if f.LastAmount == 1 {
			// Check if data will be appended otherwise preserve last data
			// We should also find the latest available element if it is there
			if willAppend {
				f.LastValue = decimal.Zero
				f.LastAmount = 0
			}

			break
		}

		f.RunningFeature.InvalidateData(log)
	}

	if f.DataStorage != nil {
		diff := uint64(0)

		if int64(TimeCurrent) - int64(f.WindowSeconds) < 0 {
			diff = 0
		} else {
			diff = TimeCurrent - f.WindowSeconds
		}
		if willAppend {
			f.DataStorage = f.DataStorage.InvalidateDataBeforeTimestamp(diff)
		}
	}

	for _, log := range data {
		// We are skipping data not in window
		if !log.IsInWindow(TimeCurrent, f.WindowSeconds) {
			continue
		}

		dataToAppend = append(dataToAppend, log)

		f.RunningFeature.CalculateData(log)
	}


	// Append in one go for less overhead
	if f.DataStorage != nil {
		f.DataStorage = f.DataStorage.Append(dataToAppend)
	}

	f.OnUpdated(TimeCurrent, data)
}

func (f *BasicRunningFeature) GetAmount() uint64 {
	return f.LastAmount
}