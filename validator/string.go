package validator

import (
	"fmt"
	"regexp"

	"github.com/akornatskyy/goext/errorstate"
	"github.com/google/uuid"
)

type StringValidatorBuilder interface {
	Required() StringValidatorBuilder
	Min(min int) StringValidatorBuilder
	Max(max int) StringValidatorBuilder
	Exactly(expected int) StringValidatorBuilder
	Pattern(pattern string, message string) StringValidatorBuilder
	Email() StringValidatorBuilder
	UUID() StringValidatorBuilder

	Build() StringValidator
}

// StringValidator validates a string value and adds any errors into
// error state.
type StringValidator interface {
	Validate(e *errorstate.ErrorState, value string) bool
}

// String creates string validator builder to setup validation rules.
func String(location string) StringValidatorBuilder {
	return &stringValidator{
		location: location,
	}
}

type stringValidator struct {
	location   string
	validators []func(*errorstate.ErrorState, string) bool
}

func (v *stringValidator) Required() StringValidatorBuilder {
	v.validators = append(v.validators, func(e *errorstate.ErrorState, value string) bool {
		if value == "" {
			e.Add(&errorstate.Detail{
				Domain:   e.Domain,
				Type:     "field",
				Location: v.location,
				Reason:   "required",
				Message:  msgRequiredField,
			})
			return false
		}
		return true
	})
	return v
}

func (v *stringValidator) Min(min int) StringValidatorBuilder {
	msg := fmt.Sprintf(msgMinLength, min)
	v.validators = append(v.validators, func(e *errorstate.ErrorState, value string) bool {
		l := len(value)
		if l != 0 && l < min {
			e.Add(&errorstate.Detail{
				Domain:   e.Domain,
				Type:     "field",
				Location: v.location,
				Reason:   "min length",
				Message:  msg,
			})
			return false
		}
		return true
	})
	return v
}

func (v *stringValidator) Max(max int) StringValidatorBuilder {
	msg := fmt.Sprintf(msgMaxLength, max)
	v.validators = append(v.validators, func(e *errorstate.ErrorState, value string) bool {
		if len(value) > max {
			e.Add(&errorstate.Detail{
				Domain:   e.Domain,
				Type:     "field",
				Location: v.location,
				Reason:   "max length",
				Message:  msg,
			})
			return false
		}
		return true
	})
	return v
}

func (v *stringValidator) Exactly(expected int) StringValidatorBuilder {
	msg := fmt.Sprintf(msgExactLength, expected)
	v.validators = append(v.validators, func(e *errorstate.ErrorState, value string) bool {
		l := len(value)
		if l != 0 && l != expected {
			e.Add(&errorstate.Detail{
				Domain:   e.Domain,
				Type:     "field",
				Location: v.location,
				Reason:   "exactly",
				Message:  msg,
			})
			return false
		}
		return true
	})
	return v
}

func (v *stringValidator) Pattern(pattern string, message string) StringValidatorBuilder {
	r := regexp.MustCompile(pattern)
	v.validators = append(v.validators, func(e *errorstate.ErrorState, value string) bool {
		if value != "" && !r.MatchString(value) {
			e.Add(&errorstate.Detail{
				Domain:   e.Domain,
				Type:     "field",
				Location: v.location,
				Reason:   "pattern",
				Message:  message,
			})
			return false
		}
		return true
	})
	return v
}

func (v *stringValidator) Email() StringValidatorBuilder {
	return v.Pattern("^[a-z0-9._%+\\-]+@[a-z0-9.\\-]+\\.[a-z]{2,4}$", msgEmail)
}

func (v *stringValidator) UUID() StringValidatorBuilder {
	v.validators = append(v.validators, func(e *errorstate.ErrorState, value string) bool {
		if value == "" {
			return true
		}
		_, err := uuid.Parse(value)
		if err != nil {
			e.Add(&errorstate.Detail{
				Domain:   e.Domain,
				Type:     "field",
				Location: v.location,
				Reason:   "pattern",
				Message:  msgUUID,
			})
			return false
		}
		return true
	})
	return v
}

func (v *stringValidator) Build() StringValidator {
	return v
}

func (v *stringValidator) Validate(s *errorstate.ErrorState, value string) bool {
	for _, validator := range v.validators {
		if !validator(s, value) {
			return false
		}
	}
	return true
}
