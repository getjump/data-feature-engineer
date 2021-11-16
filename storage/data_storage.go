package storage

import (
	dfedata "data-feature-engineer/data"
	"errors"
)

type InputDataStorage interface {
	Iterate() []*dfedata.InputData
	Append(data []*dfedata.InputData) InputDataStorage
	Clone() InputDataStorage
	Remove(item interface{}) (InputDataStorage, error)
	RemoveInputData(data []*dfedata.InputData) (InputDataStorage, error)
	InvalidateDataBeforeTimestamp(beforeTimestamp uint64) InputDataStorage
}

// LinkedListDataStorage Should be mostly immutable, we clone linked list just before operations and return pointer to new after
// Worst case scenario: We hold pointers to `len(WINDOWS_SIZES) * AMOUNT_OF_FEATURES` LinkedListDataStorage instances
// Which are all different because they all rely on different data
// Next time I will probably implement adapter for interface for better testability ¯\_(ツ)_/¯
type LinkedListDataStorage struct {
	length uint64
	item *LinkedListDataStorageItem

	head *LinkedListDataStorageItem
}

type LinkedListDataStorageItem struct {
	prev *LinkedListDataStorageItem
	next *LinkedListDataStorageItem

	data *dfedata.InputData
}

func (item *LinkedListDataStorageItem) Clone() (result *LinkedListDataStorageItem) {
	return &LinkedListDataStorageItem{prev: item.prev, next: item.next, data: item.data}
}

// Iterate should not be immutable, since that is kinda syntactic sugar to use with go range
func (storage *LinkedListDataStorage) Iterate() []*dfedata.InputData {
	result := make([]*dfedata.InputData, 0, storage.length)

	// We first reset pointer to start of the list
	item := storage.head

	// Then we are building slice, since result is preallocated (capacity), and we are using pointers there shouldn't be much overhead
	for item != nil {
		//fmt.Printf("%#v\n", item)
		result = append(result, item.data)
		item = item.next
	}

	return result
}

func (storage *LinkedListDataStorage) Clone() InputDataStorage {
	// Empty list
	if storage.item == nil {
		return &LinkedListDataStorage{item: nil, length: 0, head: nil}
	}

	item := storage.item.Clone()
	resultItem := item

	for item.prev != nil {
		newItem := item.prev.Clone()
		item.prev = newItem
		item = item.prev
	}

	return &LinkedListDataStorage{item: resultItem, length: storage.length, head: storage.head.Clone()}
}

// CloneWithLookup is O(N)
func (storage *LinkedListDataStorage) CloneWithLookup(lookupValues ...*LinkedListDataStorageItem) (result *LinkedListDataStorage, resultLookupValues []*LinkedListDataStorageItem) {
	mapLookupValues := make(map[*LinkedListDataStorageItem]int, len(lookupValues))

	for position, value := range lookupValues {
		mapLookupValues[value] = position
	}
	resultLookupValues = make([]*LinkedListDataStorageItem, len(lookupValues))

	// Empty list
	if storage.item == nil {
		result = &LinkedListDataStorage{item: nil, length: 0, head: nil}
		return
	}
	item := storage.item.Clone()
	resultItem := item

	if position, ok := mapLookupValues[storage.item]; ok {
		resultLookupValues[position] = item
	}

	for item.prev != nil {
		newItem := item.prev.Clone()
		if position, ok := mapLookupValues[item.prev]; ok {
			resultLookupValues[position] = newItem
		}
		item.prev = newItem
		item = item.prev
	}

	result = &LinkedListDataStorage{item: resultItem, length: storage.length, head: item }

	return
}

func (storage *LinkedListDataStorage) Append(data []*dfedata.InputData) InputDataStorage {
	result, _ := storage.CloneWithLookup()

	for index, item := range data {
		// Should be a first element inserted
		// No need to store end, because item is guaranteed to point to end
		if result.item == nil {
			result.item = &LinkedListDataStorageItem{data: item}
			// Probably no need to check for that
			if result.length == 0 && index == 0 {
				result.head = result.item
			}
		} else {
			result.item.next = &LinkedListDataStorageItem{prev: result.item, data: item}
			result.item = result.item.next
		}
	}

	result.length += uint64(len(data))

	return result
}

// Remove we want to lookup pointer in cloned linked list
func (storage *LinkedListDataStorage) Remove(item interface{}) (result InputDataStorage, err error) {
	linkedItem := item.(*LinkedListDataStorageItem)
	cloned, lookupValues := storage.CloneWithLookup(linkedItem)

	var newLinkedItem *LinkedListDataStorageItem

	if len(lookupValues) == 1 {
		newLinkedItem = lookupValues[0]
	} else {
		err = errors.New("cloning went wrong, cant find pointer in cloned list")
		return
	}

	err = nil

	if cloned.length == 1 {
		cloned.item = nil
		cloned.head = nil
		cloned.length = 0
		result = cloned
		return
	}

	if newLinkedItem != nil {
		if newLinkedItem == cloned.head {
			if newLinkedItem.prev != nil {
				cloned.head = newLinkedItem.prev
			} else if newLinkedItem.next != nil {
				cloned.head = newLinkedItem.next
			}
		}

		if newLinkedItem.prev != nil {
			newLinkedItem.prev.next = newLinkedItem.next
		}

		if newLinkedItem.next != nil {
			newLinkedItem.next.prev = newLinkedItem.prev
		}

		cloned.length -= 1
	} else {
		err = errors.New("something was not right")
	}

	result = cloned

	return
}

// RemoveInputData We can do not so efficient input data slice removal which is worst case O(N) to delete from linked list, so to delete M data it is O(N*M) algorithm
// But to properly implement interface and in case we need to replace our storage with something else
// We can implement removal like so:
// 1. Reset to head of linked list which is guaranteed to be lowest timestamp available,
// 2. We iterate over linked list while comparing *InputData pointers which MUST BE GUARANTEED invariant in clones of linked list
// 3. On match, we remove element from linked list which is O(1)
// So we just need to go to linked list head, which we can actually store (and end also) - and traverse likely M elements, which can be deleted in O(1)
func (storage *LinkedListDataStorage) RemoveInputData(data []*dfedata.InputData) (result InputDataStorage, err error) {
	cloned, _ := storage.CloneWithLookup()

	i := 0
	item := cloned.head

	// data must be sorted
	for item != nil && i < len(data) {
		if item.data == data[i] {
			// remove
			if item.next != nil {
				item.next.prev = item.prev
			}

			if item.prev != nil {
				item.prev.next = item.next
			}

			i++
		}

		item = item.next
	}

	err = nil

	if i != len(data) {
		err = errors.New("not all elements were deleted")
	}

	result = cloned

	return
}

func (storage *LinkedListDataStorage) InvalidateDataBeforeTimestamp(beforeTimestamp uint64) InputDataStorage {
	result, _ := storage.CloneWithLookup()

	item := result.head

	// data must be sorted
	for item != nil {
		if item.data.Timestamp < beforeTimestamp {
			// remove

			if item.next != nil {
				if item == result.head {
					result.head = item.next
				}
				item.next.prev = item.prev
			}

			if item.prev != nil {
				item.prev = nil
			}

			if result.length != 0 {
				result.length -= 1
				result.item = item.next
			} else {
				result.item = nil
				result.head = nil
			}
		} else {
			// No need to continue on sorted list
			break
		}

		item = item.next
	}

	return result
}