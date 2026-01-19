package model

type User struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
	Name     string `bson:"name"`
	Admin    bool   `bson:"admin"`
}
