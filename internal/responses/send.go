package responses

import (
	"time"

	"github.com/mdanialr/sns_backend/internal/domain"
	paginate "github.com/mdanialr/sns_backend/pkg/pagination"
)

// SendResponse adapted response for Send from domain.SNS.
type SendResponse struct {
	ID          uint       `json:"id,omitempty"`
	Url         string     `json:"url,omitempty"`
	Description string     `json:"description"`
	Send        *string    `json:"send,omitempty"`
	FileSize    *string    `json:"file_size,omitempty"`
	IsPermanent *bool      `json:"permanent,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// FromDomain adapt given domain.SNS to SendResponse.
func (s *SendResponse) FromDomain(sns *domain.SNS) {
	if sns != nil {
		s.ID = sns.ID
		s.Url = sns.Url
		s.Description = sns.Description
		s.Send = sns.Send
		s.FileSize = sns.FileSize
		s.IsPermanent = sns.IsPermanent
		s.CreatedAt = sns.CreatedAt
		s.UpdatedAt = sns.UpdatedAt
	}
}

// SendIndexResponse holds necessary data that should be used by handler to
// give response.
type SendIndexResponse struct {
	Data       []*SendResponse
	Pagination *paginate.M
}

// FromDomain setup Data from given sns data from domain/DB.
func (s *SendIndexResponse) FromDomain(sns []*domain.SNS) {
	for _, sn := range sns {
		d := &SendResponse{
			ID:          sn.ID,
			Url:         sn.Url,
			Description: sn.Description,
			Send:        sn.Send,
			FileSize:    sn.FileSize,
			IsPermanent: sn.IsPermanent,
			CreatedAt:   sn.CreatedAt,
			UpdatedAt:   sn.UpdatedAt,
		}
		s.Data = append(s.Data, d)
	}
}
