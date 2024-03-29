package response

import "github.com/1995parham-teaching/fandogh/internal/model"

// Login contains the information from the login endpoint in case of successful login.
type Login struct {
	AccessToken string `json:"accessToken"`
	model.User
}
