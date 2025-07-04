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
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate register request payload.
func (r Register) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required, validation.Length(PasswordMinLength, PasswordMaxLength)),
	)
	if err != nil {
		return fmt.Errorf("register request validation failed: %w", err)
	}

	return nil
}

// Login represents a login request payload.
type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Validate login request payload.
func (r Login) Validate() error {
	err := validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required, validation.Length(PasswordMinLength, PasswordMaxLength)),
	)
	if err != nil {
		return fmt.Errorf("login request validation failed: %w", err)
	}

	return nil
}
