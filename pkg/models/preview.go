package models

import (
	"github.com/google/uuid"
	"time"
)

// UserPreview is used as a generic subset of data for a user preview.
//
// FirstName and LastName are given separately, and their content is empty if the Username is set.
// This allows the frontend locales to control the display order of the name.
type UserPreview struct {
	// FirstName is used for display.
	FirstName string `json:"firstName"`
	// LastName is used for display.
	LastName string `json:"lastName"`
	// Username is used for display.
	Username string `json:"username"`
	// Slug is used to access the public URL of the current user.
	Slug string `json:"slug"`
	// CreatedAt gives information about the creation date of the user.
	CreatedAt time.Time `json:"createdAt"`
}

// UserPreviewPrivate extends the UserPreview object, with some private data for the current user.
type UserPreviewPrivate struct {
	ID uuid.UUID `json:"id"`
	// Email is used for display in the menu bar.
	Email string `json:"email"`
	// NewEmail indicates whether a new email is pending update for the current user.
	NewEmail string `json:"newEmail"`
	// Validated indicates whether the current user has validated their email address.
	Validated bool `json:"validated"`

	UserPreview
}
