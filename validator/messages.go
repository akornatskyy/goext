package validator

const (
	msgRequiredField = "Required field cannot be left blank."
	msgMinLength     = "Required to be a minimum of %d characters in length."
	msgMaxLength     = "Exceeds maximum length of %d."
	msgExactLength   = "Required to be exactly %d characters in length."
	msgEmail         = "Required to be a valid email address."

	msgMinRange = "Required to be greater or equal to %d."
	msgMaxRange = "Exceeds maximum allowed value of %d."
	msgExactly  = "The value must be exactly %d."
	msgRange    = "The value must fall within the range %d - %d."
)
