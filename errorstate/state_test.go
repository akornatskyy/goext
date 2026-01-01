package errorstate

import (
	"encoding/json"
	"testing"
)

func TestSingle(t *testing.T) {
	d := &Detail{
		Domain:   "validation",
		Type:     "field",
		Location: "email",
		Reason:   "invalid_format",
		Message:  "Invalid email format",
	}

	state := Single(d)

	if state == nil {
		t.Fatal("Single() returned nil")
	}
	if len(state.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(state.Errors))
	}
	if state.Errors[0] != d {
		t.Error("Error detail does not match")
	}
}

func TestErrorState_Add(t *testing.T) {
	state := &ErrorState{}

	d1 := &Detail{
		Domain:   "validation",
		Location: "email",
		Message:  "Invalid email",
	}
	d2 := &Detail{
		Domain:   "validation",
		Location: "password",
		Message:  "Password too short",
	}

	state.Add(d1)
	if len(state.Errors) != 1 {
		t.Errorf("Expected 1 error after first Add, got %d", len(state.Errors))
	}

	state.Add(d2)
	if len(state.Errors) != 2 {
		t.Errorf("Expected 2 errors after second Add, got %d", len(state.Errors))
	}

	if state.Errors[0] != d1 {
		t.Error("First error does not match")
	}
	if state.Errors[1] != d2 {
		t.Error("Second error does not match")
	}
}

func TestErrorState_Add_Multiple(t *testing.T) {
	state := &ErrorState{}

	// Add multiple errors
	for i := 0; i < 5; i++ {
		state.Add(&Detail{
			Location: "field",
			Message:  "error",
		})
	}

	if len(state.Errors) != 5 {
		t.Errorf("Expected 5 errors, got %d", len(state.Errors))
	}
}

func TestErrorState_OrNil_WithErrors(t *testing.T) {
	state := &ErrorState{}
	state.Add(&Detail{Message: "error"})

	err := state.OrNil()
	if err == nil {
		t.Error("OrNil() should return error when errors exist")
	}
	if err != state {
		t.Error("OrNil() should return the ErrorState itself")
	}
}

func TestErrorState_OrNil_NoErrors(t *testing.T) {
	state := &ErrorState{}

	err := state.OrNil()
	if err != nil {
		t.Errorf("OrNil() should return nil when no errors, got %v", err)
	}
}

func TestErrorState_OrNil_NilErrors(t *testing.T) {
	state := &ErrorState{
		Errors: nil,
	}

	err := state.OrNil()
	if err != nil {
		t.Errorf("OrNil() should return nil when Errors is nil, got %v", err)
	}
}

func TestErrorState_Error(t *testing.T) {
	state := &ErrorState{}
	state.Add(&Detail{Message: "some error"})

	errMsg := state.Error()
	expected := "errorstate: multiple errors"
	if errMsg != expected {
		t.Errorf("Error() = %q, want %q", errMsg, expected)
	}
}

func TestErrorState_Error_Empty(t *testing.T) {
	state := &ErrorState{}

	errMsg := state.Error()
	expected := "errorstate: multiple errors"
	if errMsg != expected {
		t.Errorf("Error() = %q, want %q", errMsg, expected)
	}
}

func TestDetail_Fields(t *testing.T) {
	d := &Detail{
		Domain:   "auth",
		Type:     "header",
		Location: "Authorization",
		Reason:   "missing_token",
		Message:  "Authorization token is required",
	}

	if d.Domain != "auth" {
		t.Errorf("Domain = %q, want %q", d.Domain, "auth")
	}
	if d.Type != "header" {
		t.Errorf("Type = %q, want %q", d.Type, "header")
	}
	if d.Location != "Authorization" {
		t.Errorf("Location = %q, want %q", d.Location, "Authorization")
	}
	if d.Reason != "missing_token" {
		t.Errorf("Reason = %q, want %q", d.Reason, "missing_token")
	}
	if d.Message != "Authorization token is required" {
		t.Errorf("Message = %q, want %q", d.Message, "Authorization token is required")
	}
}

func TestErrorState_JSON_Marshaling(t *testing.T) {
	state := &ErrorState{
		Domain: "validation",
	}
	state.Add(&Detail{
		Domain:   "user",
		Type:     "field",
		Location: "email",
		Reason:   "invalid",
		Message:  "Invalid email format",
	})
	state.Add(&Detail{
		Domain:   "user",
		Type:     "field",
		Location: "age",
		Reason:   "out_of_range",
		Message:  "Age must be between 18 and 120",
	})

	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("Failed to marshal ErrorState: %v", err)
	}

	// Verify Domain is omitted (json:"-" tag)
	jsonStr := string(data)
	if contains(jsonStr, "validation") {
		t.Error("Domain field should be omitted from JSON")
	}

	// Unmarshal and verify
	var unmarshaled ErrorState
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("Failed to unmarshal ErrorState: %v", err)
	}

	if len(unmarshaled.Errors) != 2 {
		t.Errorf("Expected 2 errors after unmarshal, got %d", len(unmarshaled.Errors))
	}
	if unmarshaled.Errors[0].Location != "email" {
		t.Errorf("First error location = %q, want %q", unmarshaled.Errors[0].Location, "email")
	}
	if unmarshaled.Errors[1].Location != "age" {
		t.Errorf("Second error location = %q, want %q", unmarshaled.Errors[1].Location, "age")
	}
}

func TestErrorState_JSON_EmptyErrors(t *testing.T) {
	state := &ErrorState{
		Domain: "test",
		Errors: nil,
	}

	data, err := json.Marshal(state)
	if err != nil {
		t.Fatalf("Failed to marshal ErrorState: %v", err)
	}

	// With omitempty, nil Errors should result in empty object
	expected := "{}"
	if string(data) != expected {
		t.Errorf("JSON = %s, want %s", string(data), expected)
	}
}

func TestDetail_JSON_OmitEmpty(t *testing.T) {
	// Detail with only Message field set
	d := &Detail{
		Message: "Error message",
	}

	data, err := json.Marshal(d)
	if err != nil {
		t.Fatalf("Failed to marshal Detail: %v", err)
	}

	jsonStr := string(data)
	// All fields except Message should be omitted
	if contains(jsonStr, "domain") {
		t.Error("Empty domain should be omitted")
	}
	if contains(jsonStr, "type") {
		t.Error("Empty type should be omitted")
	}
	if contains(jsonStr, "location") {
		t.Error("Empty location should be omitted")
	}
	if contains(jsonStr, "reason") {
		t.Error("Empty reason should be omitted")
	}
	if !contains(jsonStr, "message") {
		t.Error("Message should be present")
	}
}

func TestErrorState_AsError_Interface(t *testing.T) {
	var err error
	state := &ErrorState{}
	state.Add(&Detail{Message: "test error"})

	// Verify ErrorState implements error interface
	err = state

	errMsg := err.Error()
	if errMsg != "errorstate: multiple errors" {
		t.Errorf("Error message = %q, want %q", errMsg, "errorstate: multiple errors")
	}
}

func TestSingle_WithNilDetail(t *testing.T) {
	state := Single(nil)

	if state == nil {
		t.Fatal("Single(nil) returned nil")
	}
	if len(state.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(state.Errors))
	}
	if state.Errors[0] != nil {
		t.Error("Expected nil detail in errors slice")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(substr) == 0 ||
			indexOfSubstring(s, substr) >= 0)
}

func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
