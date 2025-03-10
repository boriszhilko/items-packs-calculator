package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"items-packs-calculator/internal/config"
	"items-packs-calculator/internal/packcalculator"
)

const (
	headerContentType  = "Content-Type"
	headerAllowOrigin  = "Access-Control-Allow-Origin"
	headerAllowMethods = "Access-Control-Allow-Methods"
	headerAllowHeaders = "Access-Control-Allow-Headers"

	contentTypeJSON = "application/json"
)

// CalculationRequest is the JSON payload for the /calculate API
type CalculationRequest struct {
	Items int `json:"items"`
}

// CalculationResponse is the JSON response for the /calculate API
type CalculationResponse struct {
	PackDistribution map[int]int `json:"pack_distribution"`
	TotalItems       int         `json:"total_items"`
}

// validateCalculationRequest checks if the request is valid.
func validateCalculationRequest(req CalculationRequest) error {
	if req.Items <= 0 {
		return fmt.Errorf("items must be a positive integer")
	}
	// Limit items to 1,000,000 to prevent abuse
	if req.Items > 1_000_000 {
		return fmt.Errorf("items must be less than 1,000,000")
	}
	return nil
}

// NewCalculateHandler creates a handler for pack calculations, loading pack sizes from the config file
func NewCalculateHandler(configPath string) (http.HandlerFunc, error) {
	// Load pack sizes at startup
	packSizes, err := config.LoadPackSizes(configPath)
	if err != nil {
		return nil, fmt.Errorf("error initializing pack sizes: %w", err)
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers once
		h := w.Header()
		h.Set(headerAllowOrigin, "*")
		h.Set(headerAllowMethods, "POST, OPTIONS")
		h.Set(headerAllowHeaders, headerContentType)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		var req CalculationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Printf("Invalid JSON in request body: %v", err)
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if err := validateCalculationRequest(req); err != nil {
			log.Printf("Validation error: %v, items=%d", err, req.Items)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		distribution, err := packcalculator.CalculatePacks(req.Items, packSizes)
		if err != nil {
			log.Printf("Could not calculate packs: %v, items=%d", err, req.Items)
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		total := calculateTotal(distribution)

		resp := CalculationResponse{
			PackDistribution: distribution,
			TotalItems:       total,
		}

		log.Printf("Successful calculation, items=%d, distribution=%v", req.Items, distribution)
		w.Header().Set(headerContentType, contentTypeJSON)
		json.NewEncoder(w).Encode(resp)
	}, nil
}

// calculateTotal sums the (packSize*count) from a distribution
func calculateTotal(distribution map[int]int) int {
	var total int
	for size, count := range distribution {
		total += size * count
	}
	return total
}
