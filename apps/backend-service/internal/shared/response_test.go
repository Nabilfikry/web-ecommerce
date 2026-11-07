package shared

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestJSONResponseFormatting(t *testing.T) {
	rec := httptest.NewRecorder()

	type DummyData struct {
		Name string `json:"name"`
	}

	Success(rec, http.StatusOK, DummyData{Name: "Item 1"}, "Success message")

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var resp ApiResponse[DummyData]
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse json response: %v", err)
	}

	if !resp.Success {
		t.Errorf("expected success true, got false")
	}
	if resp.Message != "Success message" {
		t.Errorf("expected message 'Success message', got %s", resp.Message)
	}
	if resp.Data.Name != "Item 1" {
		t.Errorf("expected data name 'Item 1', got %s", resp.Data.Name)
	}
}

func TestErrorResponseFormatting(t *testing.T) {
	rec := httptest.NewRecorder()

	Error(rec, http.StatusBadRequest, "Invalid input", "VALIDATION_ERROR", []ErrorDetail{
		{Field: "email", Message: "email is invalid"},
	})

	if rec.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}

	var raw map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &raw); err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}

	if raw["success"] != false {
		t.Errorf("expected success false, got %v", raw["success"])
	}
}
