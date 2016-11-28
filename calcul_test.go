package chmoufrack

import (
	"fmt"
	"testing"
)

func TestVMADistance(t *testing.T) {
	c, err := calcul_vma_distance(12, 100, 1000)
	if err != nil {
		t.Error(err)
	}
	result := "5'00"
	if c != result {
		t.Error(c + " should equal" + result)
	}

	result, err = calcul_vma_distance(11, 100, 1000)
	if err != nil {
		t.Error(err)
	}
	if result != "5'27" {
		t.Error(c + " should equal" + result)
	}
}

func TestVMAVitesse(t *testing.T) {
	c := calcul_vma_speed(20, 50)
	result := 10.0
	if c != result {
		t.Error(fmt.Sprintf("%.0f should equal %s", result, c))
	}
}

func TestVMAPace(t *testing.T) {
	c := calcul_pace(14)
	result := "4'17"
	if c != result {
		t.Error(c + " should equal " + result)
	}
	// should add a 0 at the begin
	c = calcul_pace(11.9)
	result = "5'03"
	if c != result {
		t.Error(c + " should equal " + result)
	}

	// should have just the minutes
	c = calcul_pace(5)
	result = "12'"
	if c != result {
		t.Error(c + " should equal " + result)
	}

	// should have only the seconds
	c = calcul_pace(100)
	result = "36\""
	if c != result {
		t.Error(c + " should equal " + result)
	}

	// this should equal to 4'60 which mean 5 if our thing work well
	c = calcul_pace(12.02)
	result = "5'"
	if c != result {
		t.Error(c + " should equal " + result)
	}

}
