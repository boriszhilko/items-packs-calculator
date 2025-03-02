package packcalculator

import (
	"errors"
	"sort"
)

// dp[sum] holds the best solution for "sum" items:
//
//	packCount   -> how many packs used
//	overshoot   -> sum - items (could be negative if sum < items, or ≥0 if sum >= items)
//	combination -> map[packSize]int distribution
type dpVal struct {
	packCount   int
	overshoot   int
	combination map[int]int
}

// CalculatePacks determines the best way to fulfill an order respecting these rules:
//  1. Only whole packs can be sent
//  2. Minimize leftover (overshoot) over the requested number of items
//  3. Among solutions that achieve the same leftover, minimize the pack count
//
// It returns a map of packSize -> count of packs, or an error if packSizes is empty.
func CalculatePacks(items int, packSizes []int) (map[int]int, error) {
	if len(packSizes) == 0 {
		return nil, errors.New("no pack sizes provided")
	}

	// Sort ascending for consistency.
	sort.Ints(packSizes)

	// We allow sums up to items + largest pack size to find minimal overshoot.
	maxRange := items + packSizes[len(packSizes)-1]

	dp := make([]*dpVal, maxRange+1)

	// Base case: at sum=0, overshoot=-items, packCount=0, no packs used.
	dp[0] = &dpVal{
		packCount:   0,
		overshoot:   -items,
		combination: make(map[int]int),
	}

	// Fill dp array
	for sum := 0; sum <= maxRange; sum++ {
		if dp[sum] == nil {
			continue // no known way to get this sum
		}
		for _, pack := range packSizes {
			newSum := sum + pack
			if newSum > maxRange {
				break
			}
			current := dp[sum]

			// Build a candidate newVal
			newVal := &dpVal{
				packCount:   current.packCount + 1,
				overshoot:   newSum - items,
				combination: copyMap(current.combination), // function from existing code
			}
			newVal.combination[pack]++

			// If dp[newSum] is unset or newVal is "better," update dp[newSum].
			if dp[newSum] == nil {
				dp[newSum] = newVal
			} else if betterSolution(newVal, dp[newSum]) {
				dp[newSum] = newVal
			}
		}
	}

	// Find best solution among sums >= items
	var best *dpVal
	for s := items; s <= maxRange; s++ {
		if dp[s] == nil {
			continue
		}
		// dp[s].overshoot is s - items, guaranteed ≥0 here
		if best == nil || betterSolution(dp[s], best) {
			best = dp[s]
		}
	}
	if best == nil {
		return nil, errors.New("could not find a valid combination")
	}

	return best.combination, nil
}

// betterSolution returns true if lhs is better than rhs under these rules:
//  1. smaller overshoot is better
//  2. if overshoot is the same, fewer packCount is better
func betterSolution(lhs, rhs *dpVal) bool {
	if lhs.overshoot < rhs.overshoot {
		return true
	} else if lhs.overshoot > rhs.overshoot {
		return false
	}
	// overshoot is equal, pick fewer packs
	return lhs.packCount < rhs.packCount
}

func copyMap(m map[int]int) map[int]int {
	newMap := make(map[int]int, len(m))
	for k, v := range m {
		newMap[k] = v
	}
	return newMap
}
