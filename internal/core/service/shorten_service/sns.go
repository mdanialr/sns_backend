package shorten_service

import (
	"context"
	"errors"

	repo "github.com/mdanialr/sns_backend/internal/core/repository"
	"github.com/mdanialr/sns_backend/internal/core/repository/sns_repository"
	req "github.com/mdanialr/sns_backend/internal/requests"
	res "github.com/mdanialr/sns_backend/internal/responses"
	"github.com/mdanialr/sns_backend/pkg/logger"
)

type shService struct {
	log  logger.Writer
	repo sns_repository.IRepository
}

// NewSNSService return implementation of core business logic for Shorten
// service layer.
func NewSNSService(l logger.Writer, repo sns_repository.IRepository) IService {
	return &shService{l, repo}
}

func (s *shService) Index(ctx context.Context, sh *req.Shorten) (*res.ShortenIndexResponse, error) {
	shortens, err := s.repo.FindShorten(ctx, repo.Paginate(&sh.M), repo.Order(sh.Order+" "+sh.Sort))
	if err != nil {
		errMsg := "failed to retrieve all shorten data"
		s.log.Err(errMsg+":", err)
		return nil, errors.New(errMsg)
	}

	r := &res.ShortenIndexResponse{Pagination: &sh.M}
	r.Pagination.Paginate()
	r.FromDomain(shortens)

	return r, nil
}
