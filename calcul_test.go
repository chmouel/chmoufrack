package main

import "testing"

func TestVMADistance(t *testing.T) {
	c := calcul_vma_distance(12, 100, 1000)
	result := "5'"
	if c != result {
		t.Error(c + " should equal" + result)
	}

	result = calcul_vma_distance(11, 100, 1000)
	if result != "5'27\"" {
		t.Error(c + " should equal" + result)
	}
}

func TestVMAVitesse(t *testing.T) {
	c := calcul_vma_vitesse(20, 50)
	result := "10"
	if c != result {
		t.Error(c + " should equal " + result)
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
