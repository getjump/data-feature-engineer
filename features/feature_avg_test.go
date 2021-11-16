package features

import (
	dfedata "data-feature-engineer/data"
	"data-feature-engineer/storage"
	"github.com/shopspring/decimal"
	"testing"
)

func bootstrapAvgFeature(WindowSeconds uint64) *AvgFeature {
	dataStorage := &storage.LinkedListDataStorage{}
	feature := &AvgFeature{}
	return feature.New(WindowSeconds, storage.InputDataStorage(dataStorage))
}

func TestAvgFeature_Update(t *testing.T) {
	f1, f2 := bootstrapAvgFeature(100), bootstrapAvgFeature(100)

	var tests = []struct {
		input []*dfedata.InputData
		expected decimal.Decimal
		TimeCurrent uint64
		Feature *AvgFeature
	}{
		{ // 1
			[]*dfedata.InputData{{DecimalCost: decimal.NewFromInt(10), Timestamp: 1} },
			decimal.NewFromInt(10),
			1,
			f1,
		},
		{ // 2
			[]*dfedata.InputData{{DecimalCost: decimal.NewFromInt(10), Timestamp: 2} },
			decimal.NewFromInt(10),
			2,
			f1,
		},
		// If there were no data in period we should store previous value
		{ // 3
			[]*dfedata.InputData{},
			decimal.NewFromInt(10),
			300,
			f1,
		},
		// But only for window size via storage invalidation
		{ // 4
			[]*dfedata.InputData{{DecimalCost: decimal.NewFromInt(300), Timestamp: 301}},
			decimal.NewFromInt(300),
			301,
			f1,
		},
		{ // 5
			[]*dfedata.InputData{{DecimalCost: decimal.NewFromInt(500), Timestamp: 399}},
			decimal.NewFromInt(400),
			400,
			f1,
		},
		// We should also properly handle window
		{ // 6
			[]*dfedata.InputData{
				{DecimalCost: decimal.NewFromInt(15), Timestamp: 200},
				{DecimalCost: decimal.NewFromInt(15), Timestamp: 250},
				{DecimalCost: decimal.NewFromInt(9), Timestamp: 300},
				{DecimalCost: decimal.NewFromInt(15), Timestamp: 350},
			},
			decimal.NewFromInt(13),
			350,
			f2,
		},
	}

	for i, tt := range tests {
		tt.Feature.Update(tt.TimeCurrent, tt.input)

		actual := tt.Feature.GetValue()

		if !actual.Equal(tt.expected) {
			t.Errorf("AvgFeature_Update(%#v): expected %s, actual %s, test=%d", tt.input, tt.expected.String(), actual.String(), i+1)
		}
	}
}
