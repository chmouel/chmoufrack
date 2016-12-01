package chmoufrack

import (
	"reflect"
	"testing"
)

func TestGetVMAs(t *testing.T) {
	actual := getVMAS("14:18")
	expected := []int{14, 15, 16, 17, 18}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("GETVmas failed: '%v', got '%v", expected, actual)
	}
}
