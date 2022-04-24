package utility

import "testing"

func TestIsIn(t *testing.T) {
	if !IsIn([]string{"1", "100", "20"}, "20") {
		t.Errorf("20 should be in the list")
	}

	if IsIn([]string{"1", "100", "20"}, "200") {
		t.Errorf("200 should not be in the list")
	}
}

func TestIsInGeneric(t *testing.T) {
	if !IsInGeneric([]string{"1", "100", "20"}, "20") {
		t.Errorf("20 should be in the list")
	}

	if IsInGeneric([]string{"1", "100", "20"}, "200") {
		t.Errorf("200 should not be in the list")
	}

	if !IsInGeneric([]int{1, 2, 3, 4, 5}, 5) {
		t.Errorf("5 should be in the list")
	}

	if IsInGeneric([]int{1, 2, 3, 4, 5}, 6) {
		t.Errorf("5 should be in the list")
	}

	if !IsInGeneric([]uint{1, 2, 3, 4, 5}, 5) {
		t.Errorf("5 should be in the list")
	}

	if IsInGeneric([]uint{1, 2, 3, 4, 5}, 6) {
		t.Errorf("5 should be in the list")
	}

	if !IsInGeneric([]float64{1, 2, 3, 4, 5}, 5) {
		t.Errorf("5 should be in the list")
	}

	if IsInGeneric([]float64{1, 2, 3, 4, 5}, 6) {
		t.Errorf("5 should be in the list")
	}
}
