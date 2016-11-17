package main

import (
	"reflect"
	"testing"
)

func TestGetVMA(t *testing.T) {
	vma := []int{13}
	c := getVMA(vma)
	result := vma
	if reflect.DeepEqual(c, result) {
		t.Errorf("%v != %v", c, result)
	}

	vma = []int{13, 17}
	c = getVMA(vma)
	result = []int{13, 14, 15, 16}
	if reflect.DeepEqual(c, result) {
		t.Errorf("%v != %v", c, result)
	}

}
