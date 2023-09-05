package models

import (
	"github.com/a-novel/go-framework/types"
)

type ListQuery struct {
	IDs types.StringUUIDs `json:"ids" form:"ids"`
}

type SearchQuery struct {
	Query  string `json:"query" form:"query"`
	Limit  int    `json:"limit" form:"limit"`
	Offset int    `json:"offset" form:"offset"`
}

type ValidateEmailQuery struct {
	ID   types.StringUUID `json:"id" form:"id"`
	Code string           `json:"code" form:"code"`
}
