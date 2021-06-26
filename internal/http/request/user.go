package request

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

const (
	PasswordMinLength = 6
	PasswordMaxLength = 0
)

// Register represents a register request payload.
type Register struct {
	Name     string
	Email    string
	Password string
}

// Validate register request payload.
func (r Register) Validate() error {
	if err := validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required, validation.Length(PasswordMinLength, PasswordMaxLength)),
	); err != nil {
		return fmt.Errorf("register request validation failed: %w", err)
	}

	return nil
}