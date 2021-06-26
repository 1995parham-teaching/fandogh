package model

type Bed int

const (
	Single Bed = 1
	Double Bed = 2
)

// Home represents a home to rent. contract types and room types are string to handle them more easier.
type Home struct {
	Title           string
	Location        string
	Description     string
	Peoples         int
	Room            string
	Bed             int
	Rooms           int
	Bathrooms       int
	Smoking         bool
	Guest           bool
	Pet             bool
	BillsIncluded   bool
	Contract        string
	SecurityDeposit int
	Photos          map[string][]byte
	Price           int
}
