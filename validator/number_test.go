package validator

import (
	"testing"

	"github.com/akornatskyy/goext/errorstate"
)

func TestNumber_Min_Valid(t *testing.T) {
	validator := Number("age").Min(18).Build()
	state := &errorstate.ErrorState{Domain: "test"}

	// Test with value greater than minimum
	valid := validator.Validate(state, 25)
	if !valid {
		t.Error("Expected validation to pass for value greater than minimum")
	}
	if len(state.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(state.Errors))
	}

	// Test with value equal to minimum
	state2 := &errorstate.ErrorState{Domain: "test"}
	valid = validator.Validate(state2, 18)
	if !valid {
		t.Error("Expected validation to pass for value equal to minimum")
	}
	if len(state2.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(state2.Errors))
	}
}

func TestNumber_Min_Invalid(t *testing.T) {
	validator := Number("age").Min(18).Build()
	state := &errorstate.ErrorState{Domain: "test"}

	valid := validator.Validate(state, 15)

	if valid {
		t.Error("Expected validation to fail for value less than minimum")
	}
	if len(state.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(state.Errors))
	}

	err := state.Errors[0]
	if err.Location != "age" {
		t.Errorf("Expected location 'age', got '%s'", err.Location)
	}
	if err.Reason != "min range" {
		t.Errorf("Expected reason 'min range', got '%s'", err.Reason)
	}
	if err.Type != "field" {
		t.Errorf("Expected type 'field', got '%s'", err.Type)
	}
}

func TestNumber_Max_Valid(t *testing.T) {
	validator := Number("score").Max(100).Build()
	state := &errorstate.ErrorState{Domain: "test"}

	// Test with value less than maximum
	valid := validator.Validate(state, 75)
	if !valid {
		t.Error("Expected validation to pass for value less than maximum")
	}
	if len(state.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(state.Errors))
	}

	// Test with value equal to maximum
	state2 := &errorstate.ErrorState{Domain: "test"}
	valid = validator.Validate(state2, 100)
	if !valid {
		t.Error("Expected validation to pass for value equal to maximum")
	}
	if len(state2.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(state2.Errors))
	}
}

func TestNumber_Max_Invalid(t *testing.T) {
	validator := Number("score").Max(100).Build()
	state := &errorstate.ErrorState{Domain: "test"}

	valid := validator.Validate(state, 150)

	if valid {
		t.Error("Expected validation to fail for value greater than maximum")
	}
	if len(state.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(state.Errors))
	}

	err := state.Errors[0]
	if err.Location != "score" {
		t.Errorf("Expected location 'score', got '%s'", err.Location)
	}
	if err.Reason != "max range" {
		t.Errorf("Expected reason 'max range', got '%s'", err.Reason)
	}
}

func TestNumber_Exactly_Valid(t *testing.T) {
	validator := Number("quantity").Exactly(10).Build()
	state := &errorstate.ErrorState{Domain: "test"}

	valid := validator.Validate(state, 10)

	if !valid {
		t.Error("Expected validation to pass for exact value")
	}
	if len(state.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(state.Errors))
	}
}

func TestNumber_Exactly_Invalid(t *testing.T) {
	validator := Number("quantity").Exactly(10).Build()

	// Test with less than expected
	state := &errorstate.ErrorState{Domain: "test"}
	valid := validator.Validate(state, 5)
	if valid {
		t.Error("Expected validation to fail for value less than expected")
	}
	if len(state.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(state.Errors))
	}
	if state.Errors[0].Reason != "exactly" {
		t.Errorf("Expected reason 'exactly', got '%s'", state.Errors[0].Reason)
	}

	// Test with greater than expected
	state2 := &errorstate.ErrorState{Domain: "test"}
	valid = validator.Validate(state2, 15)
	if valid {
		t.Error("Expected validation to fail for value greater than expected")
	}
	if len(state2.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(state2.Errors))
	}
	if state2.Errors[0].Reason != "exactly" {
		t.Errorf("Expected reason 'exactly', got '%s'", state2.Errors[0].Reason)
	}
}

func TestNumber_Range_Valid(t *testing.T) {
	validator := Number("temperature").Range(-10, 50).Build()

	validValues := []int{-10, -5, 0, 25, 50}

	for _, value := range validValues {
		state := &errorstate.ErrorState{Domain: "test"}
		valid := validator.Validate(state, value)
		if !valid {
			t.Errorf("Expected validation to pass for value %d in range [-10, 50]", value)
		}
		if len(state.Errors) != 0 {
			t.Errorf("Expected no errors for value %d, got %d", value, len(state.Errors))
		}
	}
}

func TestNumber_Range_Invalid(t *testing.T) {
	validator := Number("temperature").Range(-10, 50).Build()

	invalidValues := []int{-15, -11, 51, 100}

	for _, value := range invalidValues {
		state := &errorstate.ErrorState{Domain: "test"}
		valid := validator.Validate(state, value)
		if valid {
			t.Errorf("Expected validation to fail for value %d outside range [-10, 50]", value)
		}
		if len(state.Errors) != 1 {
			t.Errorf("Expected 1 error for value %d, got %d", value, len(state.Errors))
		}
		if state.Errors[0].Reason != "range" {
			t.Errorf("Expected reason 'range', got '%s'", state.Errors[0].Reason)
		}
	}
}

func TestNumber_ChainedValidators(t *testing.T) {
	validator := Number("age").Min(18).Max(120).Build()

	// Test valid age
	state := &errorstate.ErrorState{Domain: "test"}
	valid := validator.Validate(state, 25)
	if !valid {
		t.Error("Expected validation to pass for valid age")
	}
	if len(state.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(state.Errors))
	}

	// Test too young (should fail on min)
	state2 := &errorstate.ErrorState{Domain: "test"}
	valid = validator.Validate(state2, 15)
	if valid {
		t.Error("Expected validation to fail for age below minimum")
	}
	if len(state2.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(state2.Errors))
	}
	if state2.Errors[0].Reason != "min range" {
		t.Errorf("Expected failure on 'min range', got '%s'", state2.Errors[0].Reason)
	}

	// Test too old (should fail on max)
	state3 := &errorstate.ErrorState{Domain: "test"}
	valid = validator.Validate(state3, 150)
	if valid {
		t.Error("Expected validation to fail for age above maximum")
	}
	if len(state3.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(state3.Errors))
	}
	if state3.Errors[0].Reason != "max range" {
		t.Errorf("Expected failure on 'max range', got '%s'", state3.Errors[0].Reason)
	}
}

func TestNumber_MultipleValidators_StopOnFirstError(t *testing.T) {
	validator := Number("score").Min(0).Max(100).Exactly(50).Build()
	state := &errorstate.ErrorState{Domain: "test"}

	// Test with value that fails first validator (min)
	valid := validator.Validate(state, -10)
	if valid {
		t.Error("Expected validation to fail")
	}
	if len(state.Errors) != 1 {
		t.Errorf("Expected 1 error (stopped at first failure), got %d", len(state.Errors))
	}
	if state.Errors[0].Reason != "min range" {
		t.Errorf("Expected to stop at 'min range', got '%s'", state.Errors[0].Reason)
	}
}

func TestNumber_Builder_Chaining(t *testing.T) {
	// Test that builder methods return the builder for chaining
	builder := Number("field")

	b1 := builder.Min(0)
	if b1 != builder {
		t.Error("Min() should return the same builder instance")
	}

	b2 := builder.Max(100)
	if b2 != builder {
		t.Error("Max() should return the same builder instance")
	}

	b3 := builder.Exactly(50)
	if b3 != builder {
		t.Error("Exactly() should return the same builder instance")
	}

	b4 := builder.Range(10, 90)
	if b4 != builder {
		t.Error("Range() should return the same builder instance")
	}
}

func TestNumber_NegativeValues(t *testing.T) {
	validator := Number("balance").Min(-1000).Max(1000).Build()
	state := &errorstate.ErrorState{Domain: "test"}

	// Test negative value within range
	valid := validator.Validate(state, -500)
	if !valid {
		t.Error("Expected validation to pass for negative value in range")
	}
	if len(state.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(state.Errors))
	}

	// Test negative value outside range
	state2 := &errorstate.ErrorState{Domain: "test"}
	valid = validator.Validate(state2, -1500)
	if valid {
		t.Error("Expected validation to fail for negative value outside range")
	}
	if len(state2.Errors) != 1 {
		t.Errorf("Expected 1 error, got %d", len(state2.Errors))
	}
}

func TestNumber_ZeroValue(t *testing.T) {
	validator := Number("count").Min(0).Build()
	state := &errorstate.ErrorState{Domain: "test"}

	// Test zero value
	valid := validator.Validate(state, 0)
	if !valid {
		t.Error("Expected validation to pass for zero value at minimum")
	}
	if len(state.Errors) != 0 {
		t.Errorf("Expected no errors, got %d", len(state.Errors))
	}
}

func TestNumber_ErrorStateDomainPropagation(t *testing.T) {
	validator := Number("age").Min(18).Build()
	state := &errorstate.ErrorState{Domain: "registration"}

	validator.Validate(state, 15)

	if len(state.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(state.Errors))
	}

	err := state.Errors[0]
	if err.Domain != "registration" {
		t.Errorf("Expected domain 'registration', got '%s'", err.Domain)
	}
}

func TestNumber_MultipleFields(t *testing.T) {
	ageValidator := Number("age").Min(18).Max(120).Build()
	scoreValidator := Number("score").Min(0).Max(100).Build()

	state := &errorstate.ErrorState{Domain: "test"}

	// Validate multiple fields
	ageValidator.Validate(state, 15)
	scoreValidator.Validate(state, 150)

	if len(state.Errors) != 2 {
		t.Fatalf("Expected 2 errors (one per field), got %d", len(state.Errors))
	}

	// Check that errors are for different fields
	if state.Errors[0].Location == state.Errors[1].Location {
		t.Error("Expected errors for different fields")
	}

	locations := map[string]bool{}
	for _, err := range state.Errors {
		locations[err.Location] = true
	}
	if !locations["age"] || !locations["score"] {
		t.Error("Expected errors for both 'age' and 'score' fields")
	}
}

func TestNumber_MessageFormatting(t *testing.T) {
	// Test Min message
	validator := Number("age").Min(18).Build()
	state := &errorstate.ErrorState{Domain: "test"}
	validator.Validate(state, 10)
	if len(state.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(state.Errors))
	}
	expectedMsg := "Required to be greater or equal to 18."
	if state.Errors[0].Message != expectedMsg {
		t.Errorf("Expected message '%s', got '%s'", expectedMsg, state.Errors[0].Message)
	}

	// Test Max message
	validator2 := Number("score").Max(100).Build()
	state2 := &errorstate.ErrorState{Domain: "test"}
	validator2.Validate(state2, 150)
	if len(state2.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(state2.Errors))
	}
	expectedMsg2 := "Exceeds maximum allowed value of 100."
	if state2.Errors[0].Message != expectedMsg2 {
		t.Errorf("Expected message '%s', got '%s'", expectedMsg2, state2.Errors[0].Message)
	}

	// Test Exactly message
	validator3 := Number("quantity").Exactly(10).Build()
	state3 := &errorstate.ErrorState{Domain: "test"}
	validator3.Validate(state3, 5)
	if len(state3.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(state3.Errors))
	}
	expectedMsg3 := "The value must be exactly 10."
	if state3.Errors[0].Message != expectedMsg3 {
		t.Errorf("Expected message '%s', got '%s'", expectedMsg3, state3.Errors[0].Message)
	}

	// Test Range message
	validator4 := Number("temperature").Range(0, 100).Build()
	state4 := &errorstate.ErrorState{Domain: "test"}
	validator4.Validate(state4, 150)
	if len(state4.Errors) != 1 {
		t.Fatalf("Expected 1 error, got %d", len(state4.Errors))
	}
	expectedMsg4 := "The value must fall within the range 0 - 100."
	if state4.Errors[0].Message != expectedMsg4 {
		t.Errorf("Expected message '%s', got '%s'", expectedMsg4, state4.Errors[0].Message)
	}
}

func TestNumber_RangeVsMinMax(t *testing.T) {
	// Test that Range validator works correctly compared to separate Min/Max
	rangeValidator := Number("value").Range(10, 20).Build()
	minMaxValidator := Number("value").Min(10).Max(20).Build()

	// Test value in range
	state1 := &errorstate.ErrorState{Domain: "test"}
	state2 := &errorstate.ErrorState{Domain: "test"}

	rangeValidator.Validate(state1, 15)
	minMaxValidator.Validate(state2, 15)

	if len(state1.Errors) != len(state2.Errors) {
		t.Error("Range and Min/Max should behave the same for valid values")
	}

	// Test value below range
	state3 := &errorstate.ErrorState{Domain: "test"}
	state4 := &errorstate.ErrorState{Domain: "test"}

	rangeValidator.Validate(state3, 5)
	minMaxValidator.Validate(state4, 5)

	if len(state3.Errors) == 0 || len(state4.Errors) == 0 {
		t.Error("Both should fail for value below range")
	}

	// Test value above range
	state5 := &errorstate.ErrorState{Domain: "test"}
	state6 := &errorstate.ErrorState{Domain: "test"}

	rangeValidator.Validate(state5, 25)
	minMaxValidator.Validate(state6, 25)

	if len(state5.Errors) == 0 || len(state6.Errors) == 0 {
		t.Error("Both should fail for value above range")
	}
}

func TestNumber_EdgeCases(t *testing.T) {
	// Test with very large numbers
	validator := Number("big").Min(1000000).Max(9999999).Build()
	state := &errorstate.ErrorState{Domain: "test"}

	valid := validator.Validate(state, 5000000)
	if !valid {
		t.Error("Expected validation to pass for large numbers in range")
	}

	// Test with very negative numbers
	validator2 := Number("small").Min(-9999999).Max(-1000000).Build()
	state2 := &errorstate.ErrorState{Domain: "test"}

	valid = validator2.Validate(state2, -5000000)
	if !valid {
		t.Error("Expected validation to pass for very negative numbers in range")
	}
}
