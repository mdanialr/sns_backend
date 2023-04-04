package requests

import (
	"strings"

	"github.com/go-playground/validator/v10"
	paginate "github.com/mdanialr/sns_backend/pkg/pagination"
)

// Shorten standard request object that may be used to parse request in
// /shorten endpoint.
type Shorten struct {
	Url         string `json:"url,omitempty" validate:"required"`
	Shorten     string `json:"shorten,omitempty" validate:"required,url"`
	IsPermanent bool   `json:"is_permanent,omitempty" validate:"required,boolean"`

	paginate.M
	// Order the field name to query Order. Default to id.
	Order string `json:"-" query:"order"`
	// Sort to query Order. Should be filled with either asc or desc. Default
	// to asc.
	Sort string `json:"-" query:"sort"`
}

// SetQuery do setup Order and Sort.
func (s *Shorten) SetQuery() {
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
func (s *Shorten) sanitizeQuerySort() string {
	switch strings.ToLower(s.Sort) {
	case "asc", "desc":
		return s.Sort
	}
	return ""
}

// Validate validation rules for Shorten that should be parsed from request
// body.
func (s *Shorten) Validate() validator.ValidationErrors {
	if err := validate.Struct(s); err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}
