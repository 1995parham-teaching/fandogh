package request

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// NewHome contains the home creation request payload.
type NewHome struct {
	Title           string `form:"title"            json:"title"`
	Location        string `form:"location"         json:"location"`
	Description     string `form:"description"      json:"description"`
	Peoples         int    `form:"peoples"          json:"peoples"`
	Room            string `form:"room"             json:"room"`
	Bed             string `form:"bed"              json:"bed"`
	Rooms           int    `form:"rooms"            json:"rooms"`
	Bathrooms       int    `form:"bathrooms"        json:"bathrooms"`
	Smoking         bool   `form:"smoking"          json:"smoking"`
	Guest           bool   `form:"guest"            json:"guest"`
	Pet             bool   `form:"pet"              json:"pet"`
	BillsIncluded   bool   `form:"bills_included"   json:"bills_included"`
	Contract        string `form:"contract"         json:"contract"`
	SecurityDeposit int    `form:"security_deposit" json:"security_deposit"`
	Price           int    `form:"price"            json:"price"`
}

// Validate home creation request payload.
func (r NewHome) Validate() error {
	err := validation.ValidateStruct(&r,
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
	)
	if err != nil {
		return fmt.Errorf("home creation request validation failed: %w", err)
	}

	return nil
}
