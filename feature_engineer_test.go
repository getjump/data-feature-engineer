package main

import (
	"data-feature-engineer/features"
	"data-feature-engineer/storage"
	"testing"
)

func TestFeatureEngineer_AppendFeature(t *testing.T) {
	f := features.MaxFeature{}
	f.New(100)

	f2 := features.AvgFeature{}
	dataStorage := &storage.LinkedListDataStorage{}
	f2.New(100, dataStorage)

	fe := FeatureEngineer{}
	fe.AppendFeature(&f)
	fe.AppendFeature(&f2)
}