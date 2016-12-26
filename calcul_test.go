package chmoufrack

import (
	"reflect"
	"testing"
)

func TestGetVMAs(t *testing.T) {
	actual := getVmas("14:18")
	expected := []int{14, 15, 16, 17, 18}
	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("GETVmas failed: '%v', got '%v", expected, actual)
	}
}

func TestVMADistance(t *testing.T) {
	c, err := calculVmaDistance(12, 100, 1000)
	if err != nil {
		t.Error(err)
	}
	result := "5'00"
	if c != result {
		t.Error(c + " should equal" + result)
	}

	result, err = calculVmaDistance(11, 100, 1000)
	if err != nil {
		t.Error(err)
	}
	if result != "5'27" {
		t.Error(c + " should equal" + result)
	}
}

func TestVMAVitesse(t *testing.T) {
	c := calculVmaSpeed(20, 50)
	result := 10.0
	if c != result {
		t.Errorf("%.0f should equal %.0f", result, c)
	}
}

func TestVMAPace(t *testing.T) {
	c := calculPace(14)
	result := "4'17"
	if c != result {
		t.Error(c + " should equal " + result)
	}
	// should add a 0 at the begin
	c = calculPace(11.9)
	result = "5'03"
	if c != result {
		t.Error(c + " should equal " + result)
	}

	// should have just the minutes
	c = calculPace(5)
	result = "12'"
	if c != result {
		t.Error(c + " should equal " + result)
	}

	// should have only the seconds
	c = calculPace(100)
	result = "36\""
	if c != result {
		t.Error(c + " should equal " + result)
	}

	// this should equal to 4'60 which mean 5 if our thing work well
	c = calculPace(12.02)
	result = "5'"
	if c != result {
		t.Error(c + " should equal " + result)
	}

}
