package main

import (
	dfedata "data-feature-engineer/data"
	"data-feature-engineer/features"
)

// FeatureEngineer Should also make data storage for features (should there are be N instances of data storage?)
// Probably we can make one instance and make use of immutability of Linked List implementation for example
// When feature mutate storage it mutates only copy available to itself, data are manipulated via pointers,
// So there are should not be much overhead
type FeatureEngineer struct {
	Features []features.Feature
	DataAggregator DataAggregatorInterface
}

func (f *FeatureEngineer) Update(TimeCurrent uint64, data []*dfedata.InputData) error {
	f.DataAggregator.Update(TimeCurrent, data)
	dataBatch, err := f.DataAggregator.GetDataBatch()

	if err != nil {
		return err
	}

	// We need to somehow reuse common data between features
	for _, feature := range f.Features {
		feature.Update(TimeCurrent, (*dataBatch)[feature.GetWindowSeconds()])
	}

	return nil
}

func (f *FeatureEngineer) AppendFeature(feature features.Feature) {
	f.Features = append(f.Features, feature)
}