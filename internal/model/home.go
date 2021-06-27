package model

type Bed int

const (
	Single Bed = 1
	Double Bed = 2
)

// Photo contains the information use for minio.
type Photo struct {
	Name        string
	ContentType string
}

// Home represents a home to rent. contract types and room types are string to handle them more easier.
type Home struct {
	ID              string            `bson:"_id"`
	Title           string            `bson:"title"`
	Location        string            `bson:"location"`
	Description     string            `bson:"description"`
	Peoples         int               `bson:"peoples"`
	Room            string            `bson:"room"`
	Bed             Bed               `bson:"bed"`
	Rooms           int               `bson:"rooms"`
	Bathrooms       int               `bson:"bathrooms"`
	Smoking         bool              `bson:"smoking"`
	Guest           bool              `bson:"guest"`
	Pet             bool              `bson:"pet"`
	BillsIncluded   bool              `bson:"bills_included"`
	Contract        string            `bson:"contract"`
	SecurityDeposit int               `bson:"security_deposit"`
	Photos          map[string]string `bson:"photos"`
	Price           int               `bson:"price"`
}
