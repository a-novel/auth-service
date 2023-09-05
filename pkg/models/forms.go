package models

import (
	"github.com/google/uuid"
	"time"
)

type LoginForm struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type RegisterForm struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`

	FirstName string    `json:"firstName" form:"firstName"`
	LastName  string    `json:"lastName" form:"lastName"`
	Sex       Sex       `json:"sex" form:"sex"`
	Birthday  time.Time `json:"birthday" form:"birthday"`

	Slug     string `json:"slug" form:"slug"`
	Username string `json:"username" form:"username"`
}

type UpdateEmailForm struct {
	NewEmail string `json:"newEmail" form:"newEmail"`
}

type UpdateIdentityForm struct {
	FirstName string    `json:"firstName" form:"firstName"`
	LastName  string    `json:"lastName" form:"lastName"`
	Sex       Sex       `json:"sex" form:"sex"`
	Birthday  time.Time `json:"birthday" form:"birthday"`
}

type UpdateProfileForm struct {
	Slug     string `json:"slug" form:"slug"`
	Username string `json:"username" form:"username"`
}

type UpdatePasswordForm struct {
	ID          uuid.UUID `json:"id" form:"id"`
	Code        string    `json:"code" form:"code"`
	OldPassword string    `json:"oldPassword" form:"oldPassword"`
	NewPassword string    `json:"newPassword" form:"newPassword"`
}
