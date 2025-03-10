package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"items-packs-calculator/internal/api"
)

// testConfigPath points to a small JSON file with test packSizes, e.g. configs/packs_test.json
const testConfigPath = "../configs/packs_test.json"

func TestCalculateEndpoint(t *testing.T) {
	type testCase struct {
		name       string
		items      int
		wantStatus int
		wantsBody  bool // Whether we expect a valid JSON response body with pack distribution
	}

	testCases := []testCase{
		{"NegativeItems", -5, http.StatusBadRequest, false},
		{"ZeroItems", 0, http.StatusBadRequest, false},
		{"NormalItems", 501, http.StatusOK, true},
		{"BigItems", 123456, http.StatusOK, true},
	}

	handler, err := api.NewCalculateHandler(testConfigPath)
	if err != nil {
		t.Fatalf("Failed to create handler: %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/calculate", handler)

	ts := httptest.NewServer(mux)
	defer ts.Close()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bodyBytes, _ := json.Marshal(map[string]int{"items": tc.items})

			resp, err := http.Post(ts.URL+"/calculate", "application/json", bytes.NewBuffer(bodyBytes))
			if err != nil {
				t.Fatalf("POST /calculate error: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tc.wantStatus {
				t.Fatalf("Expected %d, got %d", tc.wantStatus, resp.StatusCode)
			}

			// If we expect a valid response body, parse and check basic fields
			if tc.wantsBody && tc.wantStatus == http.StatusOK {
				var result struct {
					PackDistribution map[int]int `json:"pack_distribution"`
					TotalItems       int         `json:"total_items"`
				}
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					t.Fatalf("JSON decode error: %v", err)
				}

				if len(result.PackDistribution) == 0 {
					t.Errorf("Expected non-empty pack distribution")
				}
				if result.TotalItems < 1 {
					t.Errorf("Expected total_items > 0, got %d", result.TotalItems)
				}
			}
		})
	}
}
