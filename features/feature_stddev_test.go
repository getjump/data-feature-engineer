package features

import (
	dfedata "data-feature-engineer/data"
	"data-feature-engineer/storage"
	"github.com/shopspring/decimal"
	"testing"
)

func bootstrapStdDevFeature(WindowSeconds uint64) *StdDevFeature {
	feature := &StdDevFeature{}
	dataStorage := &storage.LinkedListDataStorage{}
	return feature.New(WindowSeconds, dataStorage)
}

func TestStdDevFeature_Update(t *testing.T) {
	f1, f2 := bootstrapStdDevFeature(100), bootstrapStdDevFeature(100)

	var tests = []struct {
		input []*dfedata.InputData
		expected decimal.Decimal
		TimeCurrent uint64
		Feature *StdDevFeature
	}{
		{ // 1: Should return 0 because standard deviation is defined for N >= 2
			[]*dfedata.InputData{{DecimalCost: decimal.NewFromInt(10), Timestamp: 1} },
			decimal.NewFromInt(0),
			1,
			f1,
		},
		{ // 2: Should return actual value because there 2 elements
			[]*dfedata.InputData{{DecimalCost: decimal.NewFromInt(15), Timestamp: 2} },
			decimal.NewFromFloat(3.5355339059327378),
			2,
			f1,
		},
		// If there were no data in period we should store previous value
		{ // 3
			[]*dfedata.InputData{},
			decimal.NewFromFloat(3.5355339059327378),
			300,
			f1,
		},
		// But only for window size via storage invalidation, since in this period only one data value is 0
		{ // 4
			[]*dfedata.InputData{{DecimalCost: decimal.NewFromInt(300), Timestamp: 301}},
			decimal.NewFromInt(0),
			301,
			f1,
		},
		{ // 5
			[]*dfedata.InputData{{DecimalCost: decimal.NewFromInt(500), Timestamp: 399}},
			decimal.NewFromFloat(141.4213562373095),
			400,
			f1,
		},
		// We should also properly handle window
		{ // 6
			[]*dfedata.InputData{
				{DecimalCost: decimal.NewFromInt(3), Timestamp: 200},
				{DecimalCost: decimal.NewFromInt(19), Timestamp: 250},
				{DecimalCost: decimal.NewFromInt(21), Timestamp: 300},
				{DecimalCost: decimal.NewFromInt(17), Timestamp: 350},
			},
			decimal.NewFromFloat(2),
			350,
			f2,
		},
	}

	for i, tt := range tests {
		tt.Feature.Update(tt.TimeCurrent, tt.input)

		actual := tt.Feature.GetValue()

		if actual.String() != tt.expected.String() {
			t.Errorf("StdDevFeature_Update(%#v): expected %s, actual %s, %d", tt.input, tt.expected, actual, i)
		}
	}
}