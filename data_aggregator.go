package main

import (
	dfeData "data-feature-engineer/data"
	"errors"
	"github.com/bits-and-blooms/bitset"
	"sort"
)

// DataAggregatorInterface should aggregate data in time manner with respect to available window sizes
type DataAggregatorInterface interface {
	Update(TimeCurrent uint64, data []*dfeData.InputData)
	GetDataBatch() (result *[][]*dfeData.InputData, err error)
	GetDataForWindow(WindowSeconds uint64) (result []*dfeData.InputData, err error)
}

type DataAggregator struct {
	// WindowSeconds are sorted to process minimal first
	WindowSeconds []uint64
	WindowSecondsMap map[uint64]int
	windowDataSlices [][]*dfeData.InputData
}

func (da *DataAggregator) New(WindowSeconds[] uint64) *DataAggregator {
	da.WindowSeconds = WindowSeconds
	sort.SliceStable(da.WindowSeconds, func (i, j int) bool {
		return WindowSeconds[i] < WindowSeconds[j]
	})

	da.windowDataSlices = make([][]*dfeData.InputData, len(da.WindowSeconds))
	da.WindowSecondsMap = make(map[uint64]int)

	for windowIndex, WindowSeconds := range da.WindowSeconds {
		da.WindowSecondsMap[WindowSeconds] = windowIndex
	}

	return da
}

// Update A Bunch of data comes into this method, then it should be split for available window sizes
// Data are sorted from the lowest timestamp to highest
// To lower amount of processing we probably can go through slice in inverse way, but then we must sort(?)
// So the lowest window can process fewer data
func (da *DataAggregator) Update(TimeCurrent uint64, data []*dfeData.InputData)  {
	isDataProcessedBitset := bitset.New(uint(len(data)))

	// This way the widest computations only for widest window size, in best case where are huge gaps
	// In data, we're processing much lower amount of data
	for windowIndex, WindowSeconds := range da.WindowSeconds {
		for i := len(data) - 1; i >= 0; i-- {
			// If element is in future we don't care and continue processing,
			// but if it is before we can stop processing
			if (data)[i].IsBeforeWindow(TimeCurrent, WindowSeconds) {
				break
			}

			if !(data)[i].IsInWindow(TimeCurrent, WindowSeconds) {
				continue
			}

			// It is already processed by another window
			if isDataProcessedBitset.Test(uint(i)) {
				continue
			}

			isDataProcessedBitset.Set(uint(i))
			// Does it always guaranteed for lower window to be with higher timestamps?
			// Yes because if it fits in window, it is highest, like FIFO^
			da.windowDataSlices[windowIndex] = append(da.windowDataSlices[windowIndex], (data)[i])

			// Maybe just join slices when supplying proper windowDataSlices, because data duplication
			// Is very bad for cheap instances because of limited RAM
			//for appendixWindow := windowIndex + 1; appendixWindow < len(da.WindowSeconds); i++ {
			//	// If lowest window contains that element, it means that all highest will also contain this
			//	// Data are inverse sorted
			//	// How do we effectively inverse slice in Go?
			//	windowDataSlices[windowIndex] = append(windowDataSlices[windowIndex], data[i])
			//}
		}
	}
}

func (da *DataAggregator) GetDataBatch() (result *[][]*dfeData.InputData, err error) {
	data := make([][]*dfeData.InputData, len(da.WindowSeconds))

	for _, WindowSeconds := range da.WindowSeconds {
		resultInner, errInner := da.GetDataForWindow(WindowSeconds)

		if errInner != nil {
			err = errInner
			return
		}

		data = append(data, resultInner)
	}

	result = &data
	return
}

func (da *DataAggregator) GetDataForWindow(WindowSeconds uint64) (result []*dfeData.InputData, err error) {
	var windowIndex int
	var ok bool

	err = nil

	if windowIndex, ok = da.WindowSecondsMap[WindowSeconds]; !ok {
		err = errors.New("that window is not prepared")
		return
	}

	resultLen := 0

	for windowIndexSub := windowIndex; windowIndexSub < len(da.WindowSeconds); windowIndexSub++ {
		resultLen += len(da.windowDataSlices[windowIndexSub])
	}

	data := make([]*dfeData.InputData, 0, resultLen)

	for windowIndexSub := windowIndex; windowIndexSub >= 0; windowIndexSub-- {
		data = append(data, da.windowDataSlices[windowIndexSub]...)
	}

	// There data are actually inverse and they are sorted
	InverseWindowDataSlice(&data)
	result = data

	return
}

// InverseWindowDataSlice inverses window data slices in place to be RAM effective
func InverseWindowDataSlice(windowDataSlice *[]*dfeData.InputData) {
	for i, j := 0, len(*windowDataSlice) - 1; i < j; i, j = i + 1, j-1 {
		(*windowDataSlice)[i], (*windowDataSlice)[j] = (*windowDataSlice)[j], (*windowDataSlice)[i]
	}
}