package models

import "time"

type Credentials struct {
	Email     string `json:"email"`
	NewEmail  string `json:"newEmail"`
	Validated bool   `json:"validated"`
}

type Identity struct {
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Sex       Sex       `json:"sex"`
	Birthday  time.Time `json:"birthday"`
}

type Profile struct {
	Username string `json:"username"`
	Slug     string `json:"slug"`
}
