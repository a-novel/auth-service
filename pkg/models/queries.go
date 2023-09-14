package models

import (
	"github.com/a-novel/go-apis"
)

type ListQuery struct {
	IDs apis.StringUUIDs `json:"ids" form:"ids"`
}

type SearchQuery struct {
	Query  string `json:"query" form:"query"`
	Limit  int    `json:"limit" form:"limit"`
	Offset int    `json:"offset" form:"offset"`
}

type ValidateEmailQuery struct {
	ID   apis.StringUUID `json:"id" form:"id"`
	Code string          `json:"code" form:"code"`
}
