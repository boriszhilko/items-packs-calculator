package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

// LoadPackSizes reads pack sizes from a JSON file and returns them as a slice of ints.
func LoadPackSizes(path string) ([]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var packs []int
	if err := json.NewDecoder(file).Decode(&packs); err != nil {
		return nil, fmt.Errorf("failed to decode config file: %w", err)
	}

	// Ensure we have at least one pack size
	if len(packs) == 0 {
		return nil, errors.New("pack sizes list is empty")
	}

	// Ensure all pack sizes > 0
	for idx, size := range packs {
		if size <= 0 {
			return nil, fmt.Errorf("invalid pack size at index %d: %d", idx, size)
		}
	}

	return packs, nil
}
