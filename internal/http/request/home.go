package request

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// PhotoInput represents a photo in base64 format for JSON requests.
type PhotoInput struct {
	Name    string `json:"name"`
	Content string `json:"content"` // base64-encoded image data
}

// NewHome contains the home creation request payload.
type NewHome struct {
	Title           string       `json:"title"`
	Location        string       `json:"location"`
	Description     string       `json:"description"`
	Peoples         int          `json:"peoples"`
	Room            string       `json:"room"`
	Bed             string       `json:"bed"`
	Rooms           int          `json:"rooms"`
	Bathrooms       int          `json:"bathrooms"`
	Smoking         bool         `json:"smoking"`
	Guest           bool         `json:"guest"`
	Pet             bool         `json:"pet"`
	BillsIncluded   bool         `json:"bills_included"`
	Contract        string       `json:"contract"`
	SecurityDeposit int          `json:"security_deposit"`
	Price           int          `json:"price"`
	Photos          []PhotoInput `json:"photos"`
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

// UpdateHome contains the home update request payload.
type UpdateHome struct {
	Title           string `json:"title"`
	Location        string `json:"location"`
	Description     string `json:"description"`
	Peoples         int    `json:"peoples"`
	Room            string `json:"room"`
	Bed             string `json:"bed"`
	Rooms           int    `json:"rooms"`
	Bathrooms       int    `json:"bathrooms"`
	Smoking         bool   `json:"smoking"`
	Guest           bool   `json:"guest"`
	Pet             bool   `json:"pet"`
	BillsIncluded   bool   `json:"bills_included"`
	Contract        string `json:"contract"`
	SecurityDeposit int    `json:"security_deposit"`
	Price           int    `json:"price"`
}

// Validate home update request payload.
func (r UpdateHome) Validate() error {
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
		return fmt.Errorf("home update request validation failed: %w", err)
	}

	return nil
}
