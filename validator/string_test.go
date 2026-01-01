package validator

import (
	"testing"

	"github.com/akornatskyy/goext/errorstate"
	"github.com/google/uuid"
)

func TestString_Required_Valid(t *testing.T) {
	validator := String("email").Required().Build()
	state := &errorstate.ErrorState{Domain: "test"}

	valid := validator.Validate(state, "test@example.com")

	if !valid {
		t.Error("Expected validation to pass for non-empty string")
	}
	if len(state.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(state.Errors))
	}
}

func TestString_Required_Empty(t *testing.T) {
	validator := String("email").Required().Build()
	state := &errorstate.ErrorState{Domain: "test"}

	valid := validator.Validate(state, "")

	if valid {
		t.Error("Expected validation to fail for empty string")
	}
	if len(state.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(state.Errors))
	}

	err := state.Errors[0]
	if err.Location != "email" {
		t.Errorf("Expected location 'email', got '%s'", err.Location)
	}
	if err.Reason != "required" {
		t.Errorf("Expected reason 'required', got '%s'", err.Reason)
	}
	if err.Type != "field" {
		t.Errorf("Expected type 'field', got '%s'", err.Type)
	}
	if err.Message != msgRequiredField {
		t.Errorf("Expected message '%s', got '%s'", msgRequiredField, err.Message)
	}
}

func TestString_Min_Valid(t *testing.T) {
	validator := String("username").Min(3).Build()
	state := &errorstate.ErrorState{Domain: "test"}

	// Test with string longer than minimum
	valid := validator.Validate(state, "john")
	if !valid {
		t.Error("Expected validation to pass for string longer than minimum")
	}
	if len(state.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(state.Errors))
	}

	// Test with string exactly at minimum
	state2 := &errorstate.ErrorState{Domain: "test"}
	valid = validator.Validate(state2, "joe")
	if !valid {
		t.Error("Expected validation to pass for string at minimum length")
	}
	if len(state2.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(state2.Errors))
	}

	// Test with empty string (should pass - min only validates non-empty)
	state3 := &errorstate.ErrorState{Domain: "test"}
	valid = validator.Validate(state3, "")
	if !valid {
		t.Error("Expected validation to pass for empty string")
	}
	if len(state3.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(state3.Errors))
	}
}

func TestString_Min_Invalid(t *testing.T) {
	validator := String("username").Min(5).Build()
	state := &errorstate.ErrorState{Domain: "test"}

	valid := validator.Validate(state, "joe")

	if valid {
		t.Error("Expected validation to fail for string shorter than minimum")
	}
	if len(state.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(state.Errors))
	}

	err := state.Errors[0]
	if err.Location != "username" {
		t.Errorf("Expected location 'username', got '%s'", err.Location)
	}
	if err.Reason != "min length" {
		t.Errorf("Expected reason 'min length', got '%s'", err.Reason)
	}
}

func TestString_Max_Valid(t *testing.T) {
	validator := String("bio").Max(100).Build()
	state := &errorstate.ErrorState{Domain: "test"}

	valid := validator.Validate(state, "Short bio")

	if !valid {
		t.Error("Expected validation to pass for string shorter than maximum")
	}
	if len(state.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(state.Errors))
	}
}

func TestString_Max_Invalid(t *testing.T) {
	validator := String("bio").Max(10).Build()
	state := &errorstate.ErrorState{Domain: "test"}

	valid := validator.Validate(state, "This is a very long biography that exceeds the limit")

	if valid {
		t.Error("Expected validation to fail for string longer than maximum")
	}
	if len(state.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(state.Errors))
	}

	err := state.Errors[0]
	if err.Location != "bio" {
		t.Errorf("Expected location 'bio', got '%s'", err.Location)
	}
	if err.Reason != "max length" {
		t.Errorf("Expected reason 'max length', got '%s'", err.Reason)
	}
}

func TestString_Exactly_Valid(t *testing.T) {
	validator := String("code").Exactly(5).Build()
	state := &errorstate.ErrorState{Domain: "test"}

	// Test with exact length
	valid := validator.Validate(state, "12345")
	if !valid {
		t.Error("Expected validation to pass for string with exact length")
	}
	if len(state.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(state.Errors))
	}

	// Test with empty string (should pass - exactly only validates non-empty)
	state2 := &errorstate.ErrorState{Domain: "test"}
	valid = validator.Validate(state2, "")
	if !valid {
		t.Error("Expected validation to pass for empty string")
	}
	if len(state2.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(state2.Errors))
	}
}

func TestString_Exactly_Invalid(t *testing.T) {
	validator := String("code").Exactly(5).Build()

	// Test shorter string
	state := &errorstate.ErrorState{Domain: "test"}
	valid := validator.Validate(state, "123")
	if valid {
		t.Error("Expected validation to fail for string shorter than exact length")
	}
	if len(state.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(state.Errors))
	}

	// Test longer string
	state2 := &errorstate.ErrorState{Domain: "test"}
	valid = validator.Validate(state2, "1234567")
	if valid {
		t.Error("Expected validation to fail for string longer than exact length")
	}
	if len(state2.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(state2.Errors))
	}

	err := state2.Errors[0]
	if err.Reason != "exactly" {
		t.Errorf("Expected reason 'exactly', got '%s'", err.Reason)
	}
}

func TestString_Pattern_Valid(t *testing.T) {
	validator := String("username").Pattern("^[a-z]+$", "Only lowercase letters allowed").Build()
	state := &errorstate.ErrorState{Domain: "test"}

	// Test matching pattern
	valid := validator.Validate(state, "johnsmith")
	if !valid {
		t.Error("Expected validation to pass for string matching pattern")
	}
	if len(state.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(state.Errors))
	}

	// Test empty string (should pass - pattern only validates non-empty)
	state2 := &errorstate.ErrorState{Domain: "test"}
	valid = validator.Validate(state2, "")
	if !valid {
		t.Error("Expected validation to pass for empty string")
	}
	if len(state2.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(state2.Errors))
	}
}

func TestString_Pattern_Invalid(t *testing.T) {
	validator := String("username").Pattern("^[a-z]+$", "Only lowercase letters allowed").Build()
	state := &errorstate.ErrorState{Domain: "test"}

	valid := validator.Validate(state, "John123")

	if valid {
		t.Error("Expected validation to fail for string not matching pattern")
	}
	if len(state.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(state.Errors))
	}

	err := state.Errors[0]
	if err.Location != "username" {
		t.Errorf("Expected location 'username', got '%s'", err.Location)
	}
	if err.Reason != "pattern" {
		t.Errorf("Expected reason 'pattern', got '%s'", err.Reason)
	}
	if err.Message != "Only lowercase letters allowed" {
		t.Errorf("Expected custom message, got '%s'", err.Message)
	}
}

func TestString_Email_Valid(t *testing.T) {
	validator := String("email").Email().Build()

	validEmails := []string{
		"test@example.com",
		"user.name@example.com",
		"user+tag@example.co.uk",
		"user_name@example-domain.com",
		"123@example.com",
		"a@b.co",
	}

	for _, email := range validEmails {
		state := &errorstate.ErrorState{Domain: "test"}
		valid := validator.Validate(state, email)
		if !valid {
			t.Errorf("Expected validation to pass for valid email: %s", email)
		}
		if len(state.Errors) != 0 {
			t.Errorf("Expected no errors for %s, got %d", email, len(state.Errors))
		}
	}
}

func TestString_Email_Invalid(t *testing.T) {
	validator := String("email").Email().Build()

	invalidEmails := []string{
		"notanemail",
		"@example.com",
		"user@",
		"user@domain",
		"user name@example.com",
		"user@domain .com",
		"User@Example.Com", // uppercase not allowed by regex
	}

	for _, email := range invalidEmails {
		state := &errorstate.ErrorState{Domain: "test"}
		valid := validator.Validate(state, email)
		if valid {
			t.Errorf("Expected validation to fail for invalid email: %s", email)
		}
		if len(state.Errors) != 1 {
			t.Errorf("Expected 1 error for %s, got %d", email, len(state.Errors))
		}
	}
}

func TestString_Email_Empty(t *testing.T) {
	validator := String("email").Email().Build()
	state := &errorstate.ErrorState{Domain: "test"}

	valid := validator.Validate(state, "")

	if !valid {
		t.Error("Expected validation to pass for empty string")
	}
	if len(state.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(state.Errors))
	}
}

func TestString_UUID_Valid(t *testing.T) {
	validator := String("id").UUID().Build()

	validUUIDs := []string{
		"550e8400-e29b-41d4-a716-446655440000",
		"6ba7b810-9dad-11d1-80b4-00c04fd430c8",
		uuid.New().String(),
	}

	for _, id := range validUUIDs {
		state := &errorstate.ErrorState{Domain: "test"}
		valid := validator.Validate(state, id)
		if !valid {
			t.Errorf("Expected validation to pass for valid UUID: %s", id)
		}
		if len(state.Errors) != 0 {
			t.Errorf("Expected no errors for %s, got %d", id, len(state.Errors))
		}
	}
}

func TestString_UUID_Invalid(t *testing.T) {
	validator := String("id").UUID().Build()

	invalidUUIDs := []string{
		"not-a-uuid",
		"12345",
		"550e8400-e29b-41d4-a716",
		"550e8400-e29b-41d4-a716-446655440000-extra",
	}

	for _, id := range invalidUUIDs {
		state := &errorstate.ErrorState{Domain: "test"}
		valid := validator.Validate(state, id)
		if valid {
			t.Errorf("Expected validation to fail for invalid UUID: %s", id)
		}
		if len(state.Errors) != 1 {
			t.Errorf("Expected 1 error for %s, got %d", id, len(state.Errors))
		}
		if state.Errors[0].Message != msgUUID {
			t.Errorf("Expected UUID message, got: %s", state.Errors[0].Message)
		}
	}
}

func TestString_UUID_Empty(t *testing.T) {
	validator := String("id").UUID().Build()
	state := &errorstate.ErrorState{Domain: "test"}

	valid := validator.Validate(state, "")

	if !valid {
		t.Error("Expected validation to pass for empty string")
	}
	if len(state.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(state.Errors))
	}
}

func TestString_ChainedValidators(t *testing.T) {
	validator := String("email").Required().Min(5).Max(100).Email().Build()

	// Test valid email
	state := &errorstate.ErrorState{Domain: "test"}
	valid := validator.Validate(state, "user@example.com")
	if !valid {
		t.Error("Expected validation to pass for valid email")
	}
	if len(state.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(state.Errors))
	}

	// Test empty (should fail on required)
	state2 := &errorstate.ErrorState{Domain: "test"}
	valid = validator.Validate(state2, "")
	if valid {
		t.Error("Expected validation to fail for empty string")
	}
	if len(state2.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(state2.Errors))
	}
	if state2.Errors[0].Reason != "required" {
		t.Errorf("Expected first failure on 'required', got '%s'", state2.Errors[0].Reason)
	}

	// Test too short (should fail on min)
	state3 := &errorstate.ErrorState{Domain: "test"}
	valid = validator.Validate(state3, "a@b")
	if valid {
		t.Error("Expected validation to fail for too short string")
	}
	if len(state3.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(state3.Errors))
	}
	if state3.Errors[0].Reason != "min length" {
		t.Errorf("Expected failure on 'min length', got '%s'", state3.Errors[0].Reason)
	}
}

func TestString_MultipleValidators_StopOnFirstError(t *testing.T) {
	validator := String("username").Required().Min(3).Max(10).Pattern("^[a-z]+$", "Lowercase only").Build()
	state := &errorstate.ErrorState{Domain: "test"}

	// Test with empty string - should stop at Required
	valid := validator.Validate(state, "")
	if valid {
		t.Error("Expected validation to fail")
	}
	if len(state.Errors) != 1 {
		t.Errorf("Expected 1 error (stopped at first failure), got %d", len(state.Errors))
	}
	if state.Errors[0].Reason != "required" {
		t.Errorf("Expected to stop at 'required', got '%s'", state.Errors[0].Reason)
	}
}

func TestString_Builder_Chaining(t *testing.T) {
	// Test that builder methods return the builder for chaining
	builder := String("field")

	b1 := builder.Required()
	if b1 != builder {
		t.Error("Required() should return the same builder instance")
	}

	b2 := builder.Min(5)
	if b2 != builder {
		t.Error("Min() should return the same builder instance")
	}

	b3 := builder.Max(10)
	if b3 != builder {
		t.Error("Max() should return the same builder instance")
	}

	b4 := builder.Exactly(7)
	if b4 != builder {
		t.Error("Exactly() should return the same builder instance")
	}

	b5 := builder.Pattern(".*", "test")
	if b5 != builder {
		t.Error("Pattern() should return the same builder instance")
	}

	b6 := builder.Email()
	if b6 != builder {
		t.Error("Email() should return the same builder instance")
	}

	b7 := builder.UUID()
	if b7 != builder {
		t.Error("UUID() should return the same builder instance")
	}
}

func TestString_ErrorStateDomainPropagation(t *testing.T) {
	validator := String("email").Required().Build()
	state := &errorstate.ErrorState{Domain: "registration"}

	validator.Validate(state, "")

	if len(state.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(state.Errors))
	}

	err := state.Errors[0]
	if err.Domain != "registration" {
		t.Errorf("Expected domain 'registration', got '%s'", err.Domain)
	}
}

func TestString_MultipleFields(t *testing.T) {
	emailValidator := String("email").Required().Email().Build()
	usernameValidator := String("username").Required().Min(3).Max(20).Build()

	state := &errorstate.ErrorState{Domain: "test"}

	// Validate multiple fields
	emailValidator.Validate(state, "")
	usernameValidator.Validate(state, "ab")

	if len(state.Errors) != 2 {
		t.Fatalf("Expected 2 errors (one per field), got %d", len(state.Errors))
	}

	// Check that errors are for different fields
	if state.Errors[0].Location == state.Errors[1].Location {
		t.Error("Expected errors for different fields")
	}
}
