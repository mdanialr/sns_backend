package responses

import (
	"time"

	"github.com/mdanialr/sns_backend/internal/domain"
	paginate "github.com/mdanialr/sns_backend/pkg/pagination"
)

// ShortenResponse adapted response for Shorten from domain.SNS.
type ShortenResponse struct {
	ID          uint       `json:"id,omitempty"`
	Url         string     `json:"url,omitempty"`
	Shorten     *string    `json:"shorten,omitempty"`
	IsPermanent *bool      `json:"permanent,omitempty"`
	CreatedAt   *time.Time `json:"created_at,omitempty"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

// FromDomain adapt given domain.SNS to ShortenResponse.
func (s *ShortenResponse) FromDomain(sns *domain.SNS) {
	if sns != nil {
		s.ID = sns.ID
		s.Url = sns.Url
		s.Shorten = sns.Shorten
		s.IsPermanent = sns.IsPermanent
		s.CreatedAt = sns.CreatedAt
		s.UpdatedAt = sns.UpdatedAt
	}
}

// ShortenIndexResponse holds necessary data that should be used by handler to
// give response.
type ShortenIndexResponse struct {
	Data       []*ShortenResponse
	Pagination *paginate.M
}

// FromDomain setup Data from given sns data from domain/DB.
func (s *ShortenIndexResponse) FromDomain(sns []*domain.SNS) {
	for _, sn := range sns {
		d := &ShortenResponse{
			ID:          sn.ID,
			Url:         sn.Url,
			Shorten:     sn.Shorten,
			IsPermanent: sn.IsPermanent,
			CreatedAt:   sn.CreatedAt,
			UpdatedAt:   sn.UpdatedAt,
		}
		s.Data = append(s.Data, d)
	}
}
