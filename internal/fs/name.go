package fs

import (
	"errors"
	"fmt"
)

var ErrInvalidName = errors.New("invalid name")

const components = 2

// Generate filename for given photo of the given home.
func Generate(homeID string, name string) string {
	return fmt.Sprintf("%s_%s", homeID, name)
}

// Parse given filename to its home and photo name.
func Parse(name string) (string, string, error) {
	var id, photo string

	n, err := fmt.Sscanf(name, "%s_%s", &id, &photo)
	if err != nil {
		return "", "", ErrInvalidName
	}

	if n != components {
		return "", "", ErrInvalidName
	}

	return id, photo, nil
}
