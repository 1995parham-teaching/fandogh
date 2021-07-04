package request_test

import (
	"testing"

	"github.com/1995parham/fandogh/internal/http/request"
)

// nolint: funlen
func TestNewHomeValidation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		rq      request.NewHome
		isValid bool
	}{
		{
			rq: request.NewHome{
				Title:           "",
				Location:        "",
				Description:     "",
				Peoples:         0,
				Room:            "",
				Bed:             "",
				Rooms:           0,
				Bathrooms:       0,
				Smoking:         false,
				Guest:           false,
				Pet:             false,
				BillsIncluded:   false,
				Contract:        "",
				SecurityDeposit: 0,
				Price:           0,
			},
			isValid: false,
		},
		{
			rq: request.NewHome{
				Title:           "sweet",
				Location:        "127.0.0.1",
				Description:     "very good home",
				Peoples:         4,
				Room:            "good",
				Bed:             "single",
				Rooms:           3,
				Bathrooms:       1,
				Smoking:         false,
				Guest:           false,
				Pet:             false,
				BillsIncluded:   false,
				Contract:        "good",
				SecurityDeposit: 100,
				Price:           100,
			},
			isValid: true,
		},
		{
			rq: request.NewHome{
				Title:           "sweet",
				Location:        "127.0.0.1",
				Description:     "very good home",
				Peoples:         4,
				Room:            "good",
				Bed:             "s",
				Rooms:           3,
				Bathrooms:       1,
				Smoking:         false,
				Guest:           false,
				Pet:             false,
				BillsIncluded:   false,
				Contract:        "good",
				SecurityDeposit: 100,
				Price:           100,
			},
			isValid: false,
		},
	}

	for _, c := range cases {
		err := c.rq.Validate()

		if c.isValid && err != nil {
			t.Fatalf("valid request %+v has error %s", c.rq, err)
		}

		if !c.isValid && err == nil {
			t.Fatalf("invalid request %+v has no error", c.rq)
		}
	}
}
