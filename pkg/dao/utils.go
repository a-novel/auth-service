package dao

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	goerrors "errors"
	"fmt"
	"io"
	"strings"
)

var (
	ErrInvalidEmailFormat             = goerrors.New("invalid email format")
	ErrMarshalSignatureKey            = goerrors.New("failed to marshal signature key")
	ErrEncodeSignatureKey             = goerrors.New("failed to encode signature key")
	ErrInvalidSignatureKeyFileContent = goerrors.New("file does not contain a valid ed25519 private key: no block found")
)

// Email represents an email address as a structure, rather than a single string. This facilitates indexing:
// for example, when looking for a user, only the User is relevant, so we may only index this field for searching.
// The String method converts the email object back to the standard string representation, in the format [user]@[domain].
type Email struct {
	// Validation is the hashed key used to validate a user email. The raw key is sent to the email address only.
	// When user manages to successfully prove its authenticity, the email is validated and this code is removed.
	// Like a password, the raw key should never be stored or cached.
	Validation string `bun:"validation_code"`
	// User of the email. This is the unique name that comes before the provider.
	User string `bun:"user"`
	// Domain is the host of the mailing service provider, for example 'gmail.com'.
	Domain string `bun:"domain"`
}

// String converts the email object back to the standard string representation, in the format [user]@[domain].
// This method should never return an invalid format, so if the email is incomplete, it should return an empty string.
func (email Email) String() string {
	// Email is not valid, so it cannot be represented properly.
	if email.User == "" || email.Domain == "" {
		return ""
	}

	return fmt.Sprintf("%s@%s", email.User, email.Domain)
}

func ParseEmail(source string) (Email, error) {
	var model Email

	parts := strings.Split(source, "@")
	if len(parts) != 2 {
		return model, ErrInvalidEmailFormat
	}

	model.User = parts[0]
	model.Domain = parts[1]

	return model, nil
}

// Password represents a hashed password, that can be safely stored in the database.
type Password struct {
	// Validation is used to reset a password, for example when the original one has been forgotten. This field
	// contains the hashed key only. The raw key is sent to the user through a secure channel (an email address),
	// and once the user has managed to prove its identity, it can then create a new password.
	Validation string `bun:"validation_code"`
	// Hashed is the hashed password, used to validate user claims when trying to authenticate.
	Hashed string `bun:"hashed"`
}

// WhereEmail returns arguments for a bun Where clause, to search for a precise email value.
//
//	db.NewSelect().Model(model).Where(WhereEmail("email", email))
func WhereEmail(source string, value Email) (string, string, string) {
	return fmt.Sprintf("%[1]s_user = ? AND %[1]s_domain = ?", source), value.User, value.Domain
}

func writeKeyToOutput(out io.Writer, key ed25519.PrivateKey) error {
	marshalledKey, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		return goerrors.Join(ErrMarshalSignatureKey, err)
	}

	if err := pem.Encode(out, &pem.Block{Type: "PRIVATE KEY", Bytes: marshalledKey}); err != nil {
		return goerrors.Join(ErrEncodeSignatureKey, err)
	}

	return nil
}

func unmarshalPrivateKey(data []byte) (ed25519.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, ErrInvalidSignatureKeyFileContent
	}
	keyData, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	key, ok := keyData.(ed25519.PrivateKey)
	if !ok {
		return nil, ErrInvalidSignatureKeyFileContent
	}

	return key, nil
}
