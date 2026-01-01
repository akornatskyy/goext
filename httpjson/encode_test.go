package httpjson

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEncode(t *testing.T) {
	type testStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	w := httptest.NewRecorder()
	data := testStruct{Name: "test", Value: 42}

	Encode(w, data, http.StatusOK)

	if w.Code != http.StatusOK {
		t.Errorf("expected status code %d, got %d", http.StatusOK, w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json; charset=UTF-8" {
		t.Errorf("expected Content-Type 'application/json; charset=UTF-8', got %q", contentType)
	}

	var result testStruct
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if result.Name != "test" {
		t.Errorf("expected name 'test', got %q", result.Name)
	}
	if result.Value != 42 {
		t.Errorf("expected value 42, got %d", result.Value)
	}
}
