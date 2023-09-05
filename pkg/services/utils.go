package services

import (
	goerrors "errors"
	"regexp"
	"time"
)

var (
	ErrTaken            = goerrors.New("this value is already used by another user")
	ErrNoSignatureMatch = goerrors.New("no secret key match the current token signature")
	ErrWrongPassword    = goerrors.New("wrong password")

	ErrMissingSignatureKeys      = goerrors.New("no signature key provided")
	ErrMissingPasswordValidation = goerrors.New("you must provide either a code or an old password")
	ErrMissingPendingValidation  = goerrors.New("no pending validation found on the user")

	ErrInvalidToken          = goerrors.New("(data) invalid token")
	ErrInvalidEmail          = goerrors.New("(data) invalid email")
	ErrInvalidPassword       = goerrors.New("(data) invalid password")
	ErrInvalidFirstName      = goerrors.New("(data) invalid first name")
	ErrInvalidLastName       = goerrors.New("(data) invalid last name")
	ErrInvalidSlug           = goerrors.New("(data) invalid slug")
	ErrInvalidUsername       = goerrors.New("(data) invalid username")
	ErrInvalidSex            = goerrors.New("(data) invalid sex")
	ErrInvalidAge            = goerrors.New("(data) invalid age")
	ErrInvalidSearchLimit    = goerrors.New("(data) invalid search limit")
	ErrInvalidTokenHeader    = goerrors.New("(data) invalid token header")
	ErrInvalidTokenPayload   = goerrors.New("(data) invalid token payload")
	ErrInvalidTokenSignature = goerrors.New("(data) invalid token signature")
	ErrInvalidValidationCode = goerrors.New("(data) invalid validation code")

	ErrIntrospectToken      = goerrors.New("(dep) failed to introspect token")
	ErrCheckPassword        = goerrors.New("(dep) failed to check password")
	ErrValidateToken        = goerrors.New("(dep) failed to validate token")
	ErrVerifyValidationCode = goerrors.New("(dep) failed to verify validation code")

	ErrCancelNewEmail           = goerrors.New("(dao) failed to cancel new email")
	ErrEmailExists              = goerrors.New("(dao) failed to check if email exists")
	ErrGetProfileBySlug         = goerrors.New("(dao) failed to retrieve profile by slug")
	ErrGetTokenStatus           = goerrors.New("(dao) failed to get token status")
	ErrListUsers                = goerrors.New("(dao) failed to list users")
	ErrGetCredentialsByEmail    = goerrors.New("(dao) failed to retrieve credentials by email")
	ErrGenerateToken            = goerrors.New("(dao) failed to generate token")
	ErrGetIdentity              = goerrors.New("(dao) failed to get identity")
	ErrGetCredentials           = goerrors.New("(dao) failed to get credentials")
	ErrGetProfile               = goerrors.New("(dao) failed to get profile")
	ErrSlugExists               = goerrors.New("(dao) failed to check if slug exists")
	ErrGenerateValidationCode   = goerrors.New("(dao) failed to generate validation code")
	ErrHashPassword             = goerrors.New("(dao) failed to hash password")
	ErrCreateUser               = goerrors.New("(dao) failed to create user")
	ErrSendValidationEmail      = goerrors.New("(dao) failed to send validation email")
	ErrUpdateEmailValidation    = goerrors.New("(dao) failed to update email validation")
	ErrUpdateNewEmailValidation = goerrors.New("(dao) failed to update new email validation")
	ErrResetPassword            = goerrors.New("(dao) failed to reset password")
	ErrGenerateSignatureKey     = goerrors.New("(dao) failed to generate signature key")
	ErrWriteSignatureKey        = goerrors.New("(dao) failed to write signature key")
	ErrListSignatureKeys        = goerrors.New("(dao) failed to list signature keys")
	ErrDeleteSignatureKey       = goerrors.New("(dao) failed to delete signature key")
	ErrSearchUsers              = goerrors.New("(dao) failed to search users")
	ErrUpdateEmail              = goerrors.New("(dao) failed to update email")
	ErrUpdateIdentity           = goerrors.New("(dao) failed to update identity")
	ErrUpdatePassword           = goerrors.New("(dao) failed to update password")
	ErrUpdateProfile            = goerrors.New("(dao) failed to update profile")
	ErrValidateEmail            = goerrors.New("(dao) failed to validate email")

	usernameRegexp = regexp.MustCompile(`^[\p{L}\p{N}\p{P}]+( ([\p{L}\p{N}\p{P}]+))*$`)
	slugRegexp     = regexp.MustCompile(`^[a-z\d]+(-[a-z\d]+)*$`)
	nameRegexp     = regexp.MustCompile(`^\p{L}+([- ']\p{L}+)*$`)
)

const (
	MinEmailLength    = 3
	MaxEmailLength    = 128
	MinPasswordLength = 2
	MaxPasswordLength = 256
	MaxSlugLength     = 64
	MaxNameLength     = 32
	MaxUsernameLength = 64
	MinAge            = 16
	MaxAge            = 150
)

func getUserAge(birthday, now time.Time) int {
	return now.In(birthday.Location()).AddDate(
		-birthday.Year(),
		// Because month and day start at 1 rather than 0, we have to account for this difference.
		// -x+1 = -(x-1), to translate them back to 0 based values.
		//
		// I did math in high school btw.
		-int(birthday.Month())+1,
		-birthday.Day()+1,
	).Year()
}
