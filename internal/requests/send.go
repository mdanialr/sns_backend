package requests

import (
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	paginate "github.com/mdanialr/sns_backend/pkg/pagination"
)

// Send standard request object that may be used to parse request in
// /send.
type Send struct {
	Url       string                `form:"url" validate:"required"`
	Send      *multipart.FileHeader `form:"send" validate:"required"`
	Permanent string                `form:"permanent" validate:"required,boolean"`

	paginate.M
	// Order the field name to query Order. Default to id.
	Order string `json:"-" query:"order"`
	// Sort to query Order. Should be filled with either asc or desc. Default
	// to asc.
	Sort string `json:"-" query:"sort"`
}

// PermanentToBool convert Permanent field to bool.
func (s *Send) PermanentToBool() bool {
	b, _ := strconv.ParseBool(s.Permanent)
	return b
}

// SetQuery do setup Order and Sort.
func (s *Send) SetQuery() {
	if s.Order == "" {
		s.Order = "id" // set default to id
	}
	// sanitize Sort
	s.Sort = s.sanitizeQuerySort()
	if s.Sort == "" {
		s.Sort = "asc" // set default to asc
	}
	// make sure the Sort is upper case
	s.Sort = strings.ToUpper(s.Sort)
}

// sanitizeQuerySort make sure Sort has the expected value.
func (s *Send) sanitizeQuerySort() string {
	switch strings.ToLower(s.Sort) {
	case "asc", "desc":
		return s.Sort
	}
	return ""
}

// Validate validation rules for Send that should be parsed from request
// body.
func (s *Send) Validate() validator.ValidationErrors {
	if err := validate.Struct(s); err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}

// SendUpdate standard request object that may be used to parse request in
// /send/update endpoint.
type SendUpdate struct {
	ID        uint   `form:"id" validate:"required,numeric"`
	Url       string `form:"url" validate:"required"`
	Send      *multipart.FileHeader
	Permanent string `form:"permanent" validate:"required,boolean"`
}

// Validate validation rules for SendUpdate.
func (s *SendUpdate) Validate() validator.ValidationErrors {
	if err := validate.Struct(s); err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}

// PermanentToBool convert Permanent field to bool.
func (s *SendUpdate) PermanentToBool() bool {
	b, _ := strconv.ParseBool(s.Permanent)
	return b
}

// SendDelete standard request object that may be used to parse request in
// /send/delete endpoint.
type SendDelete struct {
	ID uint `json:"id" validate:"required,numeric"`
}

// Validate validation rules for SendDelete.
func (s *SendDelete) Validate() validator.ValidationErrors {
	if err := validate.Struct(s); err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}
