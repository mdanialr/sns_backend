package sns_repository

import (
	"context"

	r "github.com/mdanialr/sns_backend/internal/core/repository"
	"github.com/mdanialr/sns_backend/internal/domain"
	"gorm.io/gorm"
)

type snsRepo struct {
	db *gorm.DB
}

// NewSNSRepository return implementation that can be used to interact with
// object domain.SNS.
func NewSNSRepository(db *gorm.DB) IRepository {
	return &snsRepo{db}
}

func (s *snsRepo) FindShorten(ctx context.Context, opts ...r.IOptions) ([]*domain.SNS, error) {
	// prepend condition to first element
	opts = append([]r.IOptions{r.Cons("send IS NULL")}, opts...)
	return s.findSNS(ctx, opts...)
}

// findSNS general method that may be used to retrieve all domain.SNS data.
func (s *snsRepo) findSNS(ctx context.Context, opts ...r.IOptions) ([]*domain.SNS, error) {
	q := s.db.WithContext(ctx).Model(&domain.SNS{})

	for _, opt := range opts {
		q = opt.Set(q)
	}

	var sns []*domain.SNS
	return sns, q.Find(&sns).Error
}

func (s *snsRepo) GetByID(ctx context.Context, id uint, opts ...r.IOptions) (*domain.SNS, error) {
	q := s.db.WithContext(ctx)

	for _, opt := range opts {
		q = opt.Set(q)
	}

	sns := domain.SNS{ID: id}
	return &sns, q.First(&sns).Error
}

func (s *snsRepo) GetByUrl(ctx context.Context, url string, opts ...r.IOptions) (*domain.SNS, error) {
	q := s.db.WithContext(ctx)

	for _, opt := range opts {
		q = opt.Set(q)
	}

	sns := domain.SNS{Url: url}
	return &sns, q.Where(&sns, "Url").First(&sns).Error
}

func (s *snsRepo) Create(ctx context.Context, sns *domain.SNS) (*domain.SNS, error) {
	return sns, s.db.WithContext(ctx).Create(&sns).Error
}

func (s *snsRepo) Update(ctx context.Context, sns *domain.SNS, opts ...r.IOptions) (*domain.SNS, error) {
	q := s.db.WithContext(ctx)

	for _, opt := range opts {
		q = opt.Set(q)
	}

	return sns, q.Updates(&sns).Error
}

func (s *snsRepo) DeleteByID(ctx context.Context, id uint) error {
	return s.db.WithContext(ctx).Delete(&domain.SNS{ID: id}).Error
}
