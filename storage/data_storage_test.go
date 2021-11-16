package storage

import (
	dfedata "data-feature-engineer/data"
	"github.com/shopspring/decimal"
	"reflect"
	"testing"
)

func bootstrap_linked_list() InputDataStorage {
	result := &LinkedListDataStorage{}
	return result
}

func TestLinkedListDataStorage_Append(t *testing.T) {
	list := bootstrap_linked_list()

	data := []*dfedata.InputData{
		{DecimalCost: decimal.NewFromInt(10), Timestamp: 1},
		{DecimalCost: decimal.NewFromInt(10), Timestamp: 50},
		{DecimalCost: decimal.NewFromInt(10), Timestamp: 200},
		{DecimalCost: decimal.NewFromInt(10), Timestamp: 215},
		{DecimalCost: decimal.NewFromInt(10), Timestamp: 250},
	}

	newList := list.Append(data)

	if len(list.Iterate()) != 0 {
		t.Errorf("LinkedListDataStorage.Append(%#v), modified original struct, length = %d", data, len(list.Iterate()))
	}

	result := newList.Iterate()

	if !reflect.DeepEqual(result, data) {
		t.Errorf("LinkedListDataStorage.Append(%#v), incorrect append %#v", data, newList.Iterate())
	}
}

func TestLinkedListDataStorage_Clone(t *testing.T) {
	list := bootstrap_linked_list()

	data := []*dfedata.InputData{{DecimalCost: decimal.NewFromInt(10), Timestamp: 1} }

	list = list.Append(data)

	newList := list.Clone()

	if list == newList {
		t.Errorf("LinkedListDataStorage.Clone(), pointers are same %#v = %#v", list, newList)
	}

	if list.(*LinkedListDataStorage).item == newList.(*LinkedListDataStorage).item {
		t.Errorf("LinkedListDataStorage.Clone(), item pointers are same %#v = %#v", list.(*LinkedListDataStorage).item, newList.(*LinkedListDataStorage).item)
	}

	if !reflect.DeepEqual(list.Iterate(), newList.Iterate()) {
		t.Errorf("LinkedListDataStorage.Clone(), item.data pointers are not same %#v = %#v", list.Iterate(), newList.Iterate())
	}
}

func TestLinkedListDataStorage_Remove_Single(t *testing.T) {
	list := bootstrap_linked_list()

	data := []*dfedata.InputData{{DecimalCost: decimal.NewFromInt(10), Timestamp: 1} }

	list = list.Append(data)

	newList, err := list.Remove(list.(*LinkedListDataStorage).item)

	if err != nil {
		t.Error(err)
	}

	if newList.(*LinkedListDataStorage).length == list.(*LinkedListDataStorage).length {
		t.Errorf("LinkedListDataStorage.Remove(), length not updated %d should be 0", newList.(*LinkedListDataStorage).length)
	}

	if newList.(*LinkedListDataStorage).item != nil {
		t.Errorf("LinkedListDataStorage.Remove(), item not removed %#v", newList.(*LinkedListDataStorage).item)
	}
}

func TestLinkedListDataStorage_Remove_WhenManyPresent(t *testing.T) {
	list := bootstrap_linked_list()

	data := []*dfedata.InputData{
		{DecimalCost: decimal.NewFromInt(10), Timestamp: 1},
		{DecimalCost: decimal.NewFromInt(10), Timestamp: 50},
		{DecimalCost: decimal.NewFromInt(10), Timestamp: 200},
		{DecimalCost: decimal.NewFromInt(10), Timestamp: 215},
		{DecimalCost: decimal.NewFromInt(10), Timestamp: 250},
	}

	list = list.Append(data)

	var ptrToRemove = list.(*LinkedListDataStorage).head.next

	newList, err := list.Remove(ptrToRemove)

	if err != nil {
		t.Error(err)
	}

	if newList.(*LinkedListDataStorage).length == list.(*LinkedListDataStorage).length {
		t.Errorf("LinkedListDataStorage.Remove(), length not updated %d should be 4", newList.(*LinkedListDataStorage).length)
	}

	if newList.(*LinkedListDataStorage).head.next == ptrToRemove {
		t.Errorf("LinkedListDataStorage.Remove(), item not removed %#v, %#v", newList.(*LinkedListDataStorage).head.next, ptrToRemove)
	}

	if reflect.DeepEqual(newList.Iterate(), list.Iterate()) {
		t.Errorf("LinkedListDataStorage.Remove(), item not removed %#v, %#v", newList.Iterate(), list.Iterate())
	}
}

func TestLinkedListDataStorage_InvalidateDataBeforeTimestamp(t *testing.T) {
	list := bootstrap_linked_list()

	data := []*dfedata.InputData{
		{DecimalCost: decimal.NewFromInt(10), Timestamp: 1},
		{DecimalCost: decimal.NewFromInt(10), Timestamp: 50},
		{DecimalCost: decimal.NewFromInt(10), Timestamp: 200},
		{DecimalCost: decimal.NewFromInt(10), Timestamp: 215},
		{DecimalCost: decimal.NewFromInt(10), Timestamp: 250},
	}

	list = list.Append(data)

	newList := list.InvalidateDataBeforeTimestamp(100)

	if newList.(*LinkedListDataStorage).length == list.(*LinkedListDataStorage).length {
		t.Errorf("LinkedListDataStorage.InvalidateDataBeforeTimestamp(100), length not updated %d should be 3", newList.(*LinkedListDataStorage).length)
	}

	if reflect.DeepEqual(newList.Iterate(), list.Iterate()) {
		t.Errorf("LinkedListDataStorage.InvalidateDataBeforeTimestamp(100), original list was mutated %#v == %#v", newList.Iterate(), list.Iterate())
	}

	if !reflect.DeepEqual(newList.Iterate(), data[2:5]) {
		t.Errorf("LinkedListDataStorage.InvalidateDataBeforeTimestamp(100), items were not invalidated %#v != %#v", newList.Iterate(), data[2:5])
	}
}