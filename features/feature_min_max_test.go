package features

import (
	dfedata "data-feature-engineer/data"
	"github.com/shopspring/decimal"
	"testing"
)

func TestMaxFeature_GetValue(t *testing.T) {
	var tests = []struct {
		input []*dfedata.InputData
		expected decimal.Decimal
		WindowSeconds uint64
		TimeCurrent uint64
	}{
		{
			[]*dfedata.InputData{{DecimalCost: decimal.NewFromInt(15), Timestamp: 132312}, {DecimalCost: decimal.NewFromInt(1231), Timestamp: 134312}},
			decimal.NewFromInt(1231),
			5000,
			134312,
		},
		{
			[]*dfedata.InputData{{DecimalCost: decimal.NewFromInt(15000), Timestamp: 132312}, {DecimalCost: decimal.NewFromInt(1231), Timestamp: 134312}},
			decimal.NewFromInt(15000),
			5000,
			134313,
		},
	}

	f := &MaxFeature{}

	for i, tt := range tests {
		f = f.New(tt.WindowSeconds)

		f.Update(tt.TimeCurrent, tt.input)

		actual := f.GetValue()

		if !actual.Equal(tt.expected) {
			t.Errorf("TestMaxFeature_GetValue(%#v): expected %s, actual %s, %d", tt.input, tt.expected, actual, i)
		}
	}
}

func TestMinFeature_GetValue(t *testing.T) {
	f, f1 := &MinFeature{}, &MinFeature{}
	f.New(5000)
	f1.New(1)

	var tests = []struct {
		input []*dfedata.InputData
		expected decimal.Decimal
		Feature *MinFeature
		TimeCurrent uint64
	}{
		{
			[]*dfedata.InputData{{DecimalCost: decimal.NewFromInt(15), Timestamp: 132312}, {DecimalCost: decimal.NewFromInt(1231), Timestamp: 134312}},
			decimal.NewFromInt(15),
			f,
			134312,
		},
		{
			[]*dfedata.InputData{{DecimalCost: decimal.NewFromInt(14), Timestamp: 136312}},
			decimal.NewFromInt(14),
			f,
			137312,
		},
		{
			[]*dfedata.InputData{{DecimalCost: decimal.NewFromInt(15), Timestamp: 137313}},
			decimal.NewFromInt(14),
			f,
			138312,
		},
		{
			[]*dfedata.InputData{{DecimalCost: decimal.NewFromInt(20), Timestamp: 137313}},
			decimal.NewFromInt(20),
			f1,
			137313,
		},
		{
			[]*dfedata.InputData{},
			decimal.NewFromInt(14),
			f,
			137313,
		},
	}

	for i, tt := range tests {
		tt.Feature.Update(tt.TimeCurrent, tt.input)

		actual := tt.Feature.GetValue()

		if !actual.Equal(tt.expected) {
			t.Errorf("TestMinFeature_GetValue(%#v): expected %s, actual %s, %d", tt.input, tt.expected, actual, i)
		}
	}
}