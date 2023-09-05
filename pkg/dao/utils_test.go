package dao_test

import (
	"github.com/a-novel/auth-service/pkg/dao"
	"time"
)

var (
	baseTime   = time.Date(2020, time.May, 4, 8, 0, 0, 0, time.UTC)
	updateTime = time.Date(2020, time.May, 4, 9, 0, 0, 0, time.UTC)
)

func MustParseEmail(email string) dao.Email {
	e, err := dao.ParseEmail(email)
	if err != nil {
		panic(err)
	}

	return e
}

func MustParseEmailWithValidation(email, validationCode string) dao.Email {
	daoEmail := MustParseEmail(email)
	daoEmail.Validation = validationCode

	return daoEmail
}
