package models

import (
	"github.com/google/uuid"
	"time"
)

type UserTokenHeader struct {
	IAT time.Time `json:"iat"`
	EXP time.Time `json:"exp"`
	ID  uuid.UUID `json:"id"`
}

type UserTokenPayload struct {
	ID uuid.UUID `json:"id"`
}

// UserToken represents the token issued to a user, for authentication.
type UserToken struct {
	Header  UserTokenHeader  `json:"header"`
	Payload UserTokenPayload `json:"payload"`
}

// UserTokenStatus gives information about a provided token.
type UserTokenStatus struct {
	OK        bool       `json:"ok"`
	Expired   bool       `json:"expired"`
	NotIssued bool       `json:"notIssued"`
	Malformed bool       `json:"malformed"`
	Token     *UserToken `json:"token,omitempty"`
	TokenRaw  string     `json:"tokenRaw,omitempty"`
}
