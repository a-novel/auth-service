package models

// Sex represents the biological gender of a user, either SexMale or SexFemale.
type Sex string

const (
	SexMale   Sex = "male"
	SexFemale Sex = "female"
)
