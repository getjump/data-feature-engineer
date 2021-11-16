package main

import (
	dfeData "data-feature-engineer/data"
	"github.com/shopspring/decimal"
	"reflect"
	"testing"
)

func isSliceUint64Sorted(slice []uint64) bool {
	for i := 0; i < len(slice) - 1; i++ {
		if slice[i] > slice[i+1] {
			return false
		}
	}

	return true
}

func bootstrapDataAggregator(WindowSeconds []uint64) *DataAggregator {
	dataAggregator := DataAggregator{}
	return dataAggregator.New(WindowSeconds)
}

func TestDataAggregator_New(t *testing.T) {
	dataAggregator := DataAggregator{}

	dataAggregator.New([]uint64 { 150, 50, 4, 2 })

	if !isSliceUint64Sorted(dataAggregator.WindowSeconds) {
		t.Errorf("TestDataAggregator.New WindowSeconds are not sorted %#v", dataAggregator.WindowSeconds)
	}
}

func TestDataAggregator_Update(t *testing.T) {
	data := []*dfeData.InputData {
		{DecimalCost: decimal.NewFromInt(15), Timestamp: 200},
		{DecimalCost: decimal.NewFromInt(15), Timestamp: 250},
		{DecimalCost: decimal.NewFromInt(9), Timestamp: 300},
		{DecimalCost: decimal.NewFromInt(15), Timestamp: 350},
	}

	aggregator := bootstrapDataAggregator([]uint64 {5, 15, 100, 3600})

	aggregator.Update(201, data)

	if len(aggregator.windowDataSlices[0]) != 1 {
		t.Errorf("TestDataAggregator.Update Wrong amount of elements in window slice [0] should be 1, got %d",
			len(aggregator.windowDataSlices[0]))
	}

	await := []*dfeData.InputData { data[0] }

	if !reflect.DeepEqual(await, aggregator.windowDataSlices[0]) {
		t.Errorf("TestDataAggregator.Update Window Data Slices [0] is wrong should be %s, got %s",
			await, aggregator.windowDataSlices[0])
	}

	if len(aggregator.windowDataSlices[2]) != 2 {
		t.Errorf("TestDataAggregator.Update Wrong amount of elements in window slice [2] should be 3, got %d",
			len(aggregator.windowDataSlices[2]))
	}

	await = []*dfeData.InputData { data[1], data[2] }
	result := aggregator.windowDataSlices[2]
	InverseWindowDataSlice(&result)

	if !reflect.DeepEqual(await, result) {
		t.Errorf("TestDataAggregator.Update Window Data Slices [2] is wrong should be %s, got %s",
			await, aggregator.windowDataSlices[2])
	}

	await = []*dfeData.InputData { data[3] }
	result = aggregator.windowDataSlices[3]
	InverseWindowDataSlice(&result)

	if !reflect.DeepEqual(await, result) {
		t.Errorf("TestDataAggregator.Update Window Data Slices [3] is wrong should be %s, got %s",
			await, aggregator.windowDataSlices[3])
	}
}

func TestDataAggregator_GetDataForWindow(t *testing.T) {
	data := []*dfeData.InputData {
		{DecimalCost: decimal.NewFromInt(15), Timestamp: 200},
		{DecimalCost: decimal.NewFromInt(15), Timestamp: 250},
		{DecimalCost: decimal.NewFromInt(9), Timestamp: 300},
		{DecimalCost: decimal.NewFromInt(15), Timestamp: 350},
	}

	aggregator := bootstrapDataAggregator([]uint64 {5, 15, 100, 3600})

	aggregator.Update(201, data)

	await := []*dfeData.InputData { data[0], data[1], data[2] }
	result, err := aggregator.GetDataForWindow(100)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(await, result) {
		t.Errorf("TestDataAggregator.GetDataForWindow(100) is wrong should be %s, got %s",
			await, result)
	}
}