package request

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// NewHome contains the home creation request payload.
type NewHome struct {
	Title           string `form:"title"`
	Location        string `form:"location"`
	Description     string `form:"description"`
	Peoples         int    `form:"peoples"`
	Room            string `form:"room"`
	Bed             string `form:"bed"`
	Rooms           int    `form:"rooms"`
	Bathrooms       int    `form:"bathrooms"`
	Smoking         bool   `form:"smoking"`
	Guest           bool   `form:"guest"`
	Pet             bool   `form:"pet"`
	BillsIncluded   bool   `form:"bills_included"`
	Contract        string `form:"contract"`
	SecurityDeposit int    `form:"security_deposit"`
	Price           int    `form:"price"`
}

// Validate home creation request payload.
func (r NewHome) Validate() error {
	if err := validation.ValidateStruct(&r,
		validation.Field(&r.Title, validation.Required),
		validation.Field(&r.Location, validation.Required),
		validation.Field(&r.Description, validation.Required),
		validation.Field(&r.Peoples, validation.Required),
		validation.Field(&r.Room, validation.Required),
		validation.Field(&r.Bed, validation.In("single", "double"), validation.Required),
		validation.Field(&r.Rooms, validation.Required),
		validation.Field(&r.Bathrooms, validation.Required),
		validation.Field(&r.Contract, validation.Required),
		validation.Field(&r.SecurityDeposit, validation.Required),
		validation.Field(&r.Price, validation.Required),
	); err != nil {
		return fmt.Errorf("home creation request validation failed: %w", err)
	}

	return nil
}
