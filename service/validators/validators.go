package validators

import "github.com/go-playground/validator/v10"

type Validator struct {
	*validator.Validate
}

// Return a new validator instance
func New() (v *Validator) {
	v = &Validator{validator.New()}

	return
}
