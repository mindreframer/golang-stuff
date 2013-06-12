package common

import (
	"reflect"
	"testing"
)

func TestMergeItems(t *testing.T) {
	i1 := []Item{Item{Key: []byte{4}, Timestamp: 44}}
	i2 := []Item{Item{Key: []byte{5}, Timestamp: 44}}
	ary := []*[]Item{&i1, &i2}
	expected := []Item{Item{Key: []byte{4}, Timestamp: 44}, Item{Key: []byte{5}, Timestamp: 44}}
	found := MergeItems(ary, true)
	if !reflect.DeepEqual(expected, found) {
		t.Errorf("%v should be %v", found, expected)
	}
	i1 = []Item{Item{Key: []byte{5}, Timestamp: 45}}
	ary = []*[]Item{&i1, &i2}
	expected = []Item{Item{Key: []byte{5}, Timestamp: 45}}
	found = MergeItems(ary, true)
	if !reflect.DeepEqual(expected, found) {
		t.Errorf("%v should be %v", found, expected)
	}
	i1 = []Item{Item{Key: []byte{5}, Timestamp: 45}, Item{Key: []byte{6}, Timestamp: 45}, Item{Key: []byte{7}, Timestamp: 45}, Item{Key: []byte{8}, Timestamp: 45}}
	i2 = []Item{Item{Key: []byte{5}, Timestamp: 45}, Item{Key: []byte{7}, Timestamp: 45}, Item{Key: []byte{8}, Timestamp: 45}}
	ary = []*[]Item{&i1, &i2}
	expected = []Item{Item{Key: []byte{5}, Timestamp: 45}, Item{Key: []byte{6}, Timestamp: 45}, Item{Key: []byte{7}, Timestamp: 45}, Item{Key: []byte{8}, Timestamp: 45}}
	found = MergeItems(ary, true)
	if !reflect.DeepEqual(expected, found) {
		t.Errorf("%v should be %v", found, expected)
	}
}
