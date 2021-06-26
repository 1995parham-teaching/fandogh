package model

type Home struct {
	Title           string
	Location        string
	Description     string
	Peoples         int
	Room            int
	Bed             int
	Rooms           int
	Bathrooms       int
	Smoking         bool
	Guest           bool
	Pet             bool
	BillsIncluded   bool
	Contract        int
	SecurityDeposit int
	Photos          map[string][]byte
	Price           int
}
