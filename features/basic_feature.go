package features

import (
	dfedata "data-feature-engineer/data"
	decimal "github.com/shopspring/decimal"
)

type PairDataConnectionChannelFeature struct {
	Channel chan ConnectionChannelData
	Feature *Feature
}

type ConnectionChannelData struct {
	Value decimal.Decimal
	Amount uint64
	ProvidingFeature Feature
}

type Feature interface {
	Update(TimeCurrent uint64, data []*dfedata.InputData, connectionChannel ...chan ConnectionChannelData)

	Chain (f *Feature) Feature

	GetValue() decimal.Decimal
	GetAmount() uint64

	GetWindowSeconds() uint64
}

// BasicFeature TODO caching for calculations for same feature different window sizes
// Like for example for AvgFeature that will not work
// But for MinMax that should work, we can probably chain sub minmax calls? Like divide and conquer algorithm
type BasicFeature struct {
	LastValue decimal.Decimal
	ChainedFeatures []PairDataConnectionChannelFeature

	WindowSeconds uint64
}

// Update TODO probably we should provide also Feature via connection channel to determine which feature are sending value
func (f *BasicFeature) Update(TimeCurrent uint64, data []*dfedata.InputData, connectionChannel ...chan ConnectionChannelData) {
}

func (f *BasicFeature) GetAmount() uint64 {
	return 0
}

func (f *BasicFeature) GetValue() decimal.Decimal {
	return f.LastValue
}

func (f *BasicFeature) GetWindowSeconds() uint64 {
	return f.WindowSeconds
}

// Chain example AvgFeature.Chain(AvgSquaredFeature).Chain(StdDevFeature)
func (f *BasicFeature) Chain(f2 *Feature) Feature {
	connectionChannel := make(chan ConnectionChannelData)
	f.ChainedFeatures = append(f.ChainedFeatures, PairDataConnectionChannelFeature{connectionChannel, f2})
	return *f2
}

func (f *BasicFeature) OnUpdated(TimeCurrent uint64, data[]*dfedata.InputData) {
	channelData := ConnectionChannelData { f.GetValue(), f.GetAmount(), f }
	for _, pair := range f.ChainedFeatures {
		go (*pair.Feature).Update(TimeCurrent, data, pair.Channel)
		pair.Channel <- channelData
	}
}