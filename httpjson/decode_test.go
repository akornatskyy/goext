package httpjson

import (
	"bytes"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/akornatskyy/goext/errorstate"
)

func TestDecode_Success(t *testing.T) {
	type testStruct struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	jsonBody := `{"name":"test","value":42}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	var result testStruct
	err := Decode(req, &result, 1024)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Name != "test" {
		t.Errorf("expected name 'test', got %q", result.Name)
	}
	if result.Value != 42 {
		t.Errorf("expected value 42, got %d", result.Value)
	}
}

func TestDecode_WrongContentType(t *testing.T) {
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"test"}`))
	req.Header.Set("Content-Type", "text/plain")

	var result map[string]interface{}
	err := Decode(req, &result, 1024)

	if err == nil {
		t.Fatal("expected error for wrong content type")
	}

	errState, ok := err.(*errorstate.ErrorState)
	if !ok {
		t.Fatalf("expected errorstate.ErrorState, got %T", err)
	}

	if len(errState.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errState.Errors))
	}

	detail := errState.Errors[0]
	if detail.Domain != "HTTP" {
		t.Errorf("expected domain 'HTTP', got %q", detail.Domain)
	}
	if detail.Type != "header" {
		t.Errorf("expected type 'header', got %q", detail.Type)
	}
	if detail.Location != "Content-Type" {
		t.Errorf("expected location 'Content-Type', got %q", detail.Location)
	}
	if detail.Reason != "unexpected content type" {
		t.Errorf("expected reason 'unexpected content type', got %q", detail.Reason)
	}
}

func TestDecode_MissingContentType(t *testing.T) {
	req := httptest.NewRequest("POST", "/", strings.NewReader(`{"name":"test"}`))

	var result map[string]interface{}
	err := Decode(req, &result, 1024)

	if err == nil {
		t.Fatal("expected error for missing content type")
	}

	errState, ok := err.(*errorstate.ErrorState)
	if !ok {
		t.Fatalf("expected errorstate.ErrorState, got %T", err)
	}

	if len(errState.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errState.Errors))
	}
}

func TestDecode_InvalidJSON(t *testing.T) {
	invalidJSON := `{"name":"test", "value":}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	var result map[string]interface{}
	err := Decode(req, &result, 1024)

	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}

	errState, ok := err.(*errorstate.ErrorState)
	if !ok {
		t.Fatalf("expected errorstate.ErrorState, got %T", err)
	}

	if len(errState.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errState.Errors))
	}

	detail := errState.Errors[0]
	if detail.Domain != "JSON" {
		t.Errorf("expected domain 'JSON', got %q", detail.Domain)
	}
	if detail.Type != "decode" {
		t.Errorf("expected type 'decode', got %q", detail.Type)
	}
	if detail.Location != "HTTP request body" {
		t.Errorf("expected location 'HTTP request body', got %q", detail.Location)
	}
	if detail.Message != "Unable to parse JSON." {
		t.Errorf("expected message 'Unable to parse JSON.', got %q", detail.Message)
	}
}

func TestDecode_RequestTooLarge(t *testing.T) {
	largeJSON := `{"data":"` + strings.Repeat("a", 1000) + `"}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(largeJSON))
	req.Header.Set("Content-Type", "application/json")

	var result map[string]interface{}
	err := Decode(req, &result, 10) // Limit to 10 bytes

	if err == nil {
		t.Fatal("expected error for request too large")
	}

	errState, ok := err.(*errorstate.ErrorState)
	if !ok {
		t.Fatalf("expected errorstate.ErrorState, got %T", err)
	}

	if len(errState.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errState.Errors))
	}

	detail := errState.Errors[0]
	if detail.Domain != "HTTP" {
		t.Errorf("expected domain 'HTTP', got %q", detail.Domain)
	}
	if detail.Type != "reader" {
		t.Errorf("expected type 'reader', got %q", detail.Type)
	}
	if detail.Location != "HTTP request body" {
		t.Errorf("expected location 'HTTP request body', got %q", detail.Location)
	}
	if detail.Reason != "request entity too large" {
		t.Errorf("expected reason 'request entity too large', got %q", detail.Reason)
	}
	if !strings.Contains(detail.Message, "10 bytes") {
		t.Errorf("expected message to contain '10 bytes', got %q", detail.Message)
	}
}

func TestDecode_EmptyBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/json")

	var result map[string]interface{}
	err := Decode(req, &result, 1024)

	if err == nil {
		t.Fatal("expected error for empty body")
	}

	errState, ok := err.(*errorstate.ErrorState)
	if !ok {
		t.Fatalf("expected errorstate.ErrorState, got %T", err)
	}

	if len(errState.Errors) != 1 {
		t.Fatalf("expected 1 error, got %d", len(errState.Errors))
	}

	detail := errState.Errors[0]
	if detail.Domain != "JSON" {
		t.Errorf("expected domain 'JSON', got %q", detail.Domain)
	}
}

func TestDecode_ExactLimit(t *testing.T) {
	jsonBody := `{"test":"value"}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	var result map[string]interface{}
	err := Decode(req, &result, int64(len(jsonBody)))

	if err != nil {
		t.Fatalf("expected no error for exact limit, got %v", err)
	}
	if result["test"] != "value" {
		t.Errorf("expected test='value', got %v", result["test"])
	}
}

func TestDecode_NilBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/", nil)
	req.Header.Set("Content-Type", "application/json")

	var result map[string]interface{}
	err := Decode(req, &result, 1024)

	if err == nil {
		t.Fatal("expected error for nil body")
	}
}

func TestDecode_NestedJSON(t *testing.T) {
	type Address struct {
		Street string `json:"street"`
		City   string `json:"city"`
	}
	type Person struct {
		Name    string  `json:"name"`
		Age     int     `json:"age"`
		Address Address `json:"address"`
	}

	jsonBody := `{"name":"John","age":30,"address":{"street":"Main St","city":"NYC"}}`
	req := httptest.NewRequest("POST", "/", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	var result Person
	err := Decode(req, &result, 1024)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Name != "John" {
		t.Errorf("expected name 'John', got %q", result.Name)
	}
	if result.Age != 30 {
		t.Errorf("expected age 30, got %d", result.Age)
	}
	if result.Address.Street != "Main St" {
		t.Errorf("expected street 'Main St', got %q", result.Address.Street)
	}
	if result.Address.City != "NYC" {
		t.Errorf("expected city 'NYC', got %q", result.Address.City)
	}
}

func TestDecode_Array(t *testing.T) {
	jsonBody := `[{"name":"item1"},{"name":"item2"}]`
	req := httptest.NewRequest("POST", "/", strings.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	var result []map[string]string
	err := Decode(req, &result, 1024)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("expected 2 items, got %d", len(result))
	}
	if result[0]["name"] != "item1" {
		t.Errorf("expected first item name 'item1', got %q", result[0]["name"])
	}
	if result[1]["name"] != "item2" {
		t.Errorf("expected second item name 'item2', got %q", result[1]["name"])
	}
}

func TestDecode_WithBodyReader(t *testing.T) {
	jsonBody := `{"key":"value"}`
	req := httptest.NewRequest("POST", "/", io.NopCloser(bytes.NewReader([]byte(jsonBody))))
	req.Header.Set("Content-Type", "application/json")

	var result map[string]string
	err := Decode(req, &result, 1024)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result["key"] != "value" {
		t.Errorf("expected key='value', got %v", result["key"])
	}
}
