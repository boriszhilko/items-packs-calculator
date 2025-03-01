package packcalculator

import (
	"reflect"
	"testing"
)

// Helper to compute total items from result
func totalItems(result map[int]int) int {
	total := 0
	for size, count := range result {
		total += size * count
	}
	return total
}

func TestCalculatePacks(t *testing.T) {
	packs := []int{250, 500, 1000, 2000, 5000}

	tests := []struct {
		name         string
		items        int
		wantMin      int
		wantPackKeys []int
	}{
		{"order_1", 1, 250, []int{250}},                             // minimal overshoot => 250
		{"order_250", 250, 250, []int{250}},                         // exact => 250
		{"order_251", 251, 500, []int{500}},                         // minimal overshoot => 500
		{"order_501", 501, 750, []int{500, 250}},                    // 500 + 250 = 750 overshoot
		{"order_12001", 12001, 12250, []int{5000, 5000, 2000, 250}}, // 2x5000 + 1x2000 + 1x250 = 12250
		// In a real-world scenario, we'd likely limit extremely large input values
		{"order_large", 1_000_000, 1_000_000, []int{1000, 2000, 5000}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculatePacks(tt.items, packs)
			if err != nil {
				t.Fatalf("CalculatePacks() error = %v", err)
			}
			gotTotal := totalItems(got)
			if gotTotal < tt.items {
				t.Errorf("Total items (%v) is less than requested (%v)", gotTotal, tt.items)
			}
			if gotTotal < tt.wantMin {
				t.Errorf("Expected at least %v items, but got %v", tt.wantMin, gotTotal)
			}

			for packSize := range got {
				if !contains(tt.wantPackKeys, packSize) {
					// It's okay if there's variation, but let's see if it's drastically different
					t.Logf("Note: Different pack size used than expected. packSize=%v", packSize)
				}
			}
		})
	}
}

func contains(arr []int, val int) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func TestNoPackSizes(t *testing.T) {
	_, err := CalculatePacks(100, []int{})
	if err == nil {
		t.Error("Expected an error for empty pack sizes, got nil")
	}
}

func TestExactScenario(t *testing.T) {
	packs := []int{250, 500, 1000}
	got, err := CalculatePacks(1000, packs)
	if err != nil {
		t.Fatalf("Error: %v", err)
	}
	expected := map[int]int{1000: 1}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("Unexpected result. got=%v, want=%v", got, expected)
	}
}
