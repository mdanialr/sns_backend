package shorten_service

import (
	"context"
	"errors"
	"strconv"

	repo "github.com/mdanialr/sns_backend/internal/core/repository"
	"github.com/mdanialr/sns_backend/internal/core/repository/sns_repository"
	"github.com/mdanialr/sns_backend/internal/domain"
	req "github.com/mdanialr/sns_backend/internal/requests"
	res "github.com/mdanialr/sns_backend/internal/responses"
	h "github.com/mdanialr/sns_backend/pkg/helper"
	"github.com/mdanialr/sns_backend/pkg/logger"
)

type shService struct {
	log  logger.Writer
	repo sns_repository.IRepository
}

// New return implementation of core business logic for Shorten service layer.
func New(l logger.Writer, repo sns_repository.IRepository) IService {
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

func (s *shService) Create(ctx context.Context, req *req.Shorten) (*res.ShortenResponse, error) {
	// make sure given url is not used yet in the db
	o, _ := s.repo.GetByUrl(ctx, req.Url, repo.Cols("id"))
	if o.ID != 0 {
		return nil, errors.New("url already been taken")
	}

	// prepare new object to be saved to DB
	sh := &domain.SNS{
		Url:         req.Url,
		Shorten:     &req.Shorten,
		IsPermanent: h.Ptr(req.PermanentToBool()),
	}
	if _, err := s.repo.Create(ctx, sh); err != nil {
		errMsg := "failed to create new Shorten"
		s.log.Err(errMsg+":", err)
		return nil, errors.New(errMsg)
	}

	// adapt data from domain.SNS to required ShortenResponse
	var r res.ShortenResponse
	r.FromDomain(sh)

	return &r, nil
}

func (s *shService) Update(ctx context.Context, req *req.ShortenUpdate) (*res.ShortenResponse, error) {
	// make sure given url is not used yet in the db -
	sns, _ := s.repo.GetByID(ctx, req.ID, repo.Cols("url"))
	if sns.Url != req.Url {
		// - if given id has different url from url in request
		o, _ := s.repo.GetByUrl(ctx, req.Url, repo.Cols("id"))
		if o.ID != 0 {
			return nil, errors.New("url already been taken")
		}
	}

	// prepare new object to be updated to DB
	sh := &domain.SNS{
		ID:          req.ID,
		Url:         req.Url,
		Shorten:     req.Shorten,
		IsPermanent: h.Ptr(req.PermanentToBool()),
	}
	newSh, err := s.repo.Update(ctx, sh)
	if err != nil {
		errMsg := "failed to update Shorten with id " + strconv.Itoa(int(req.ID))
		s.log.Err(errMsg+":", err)
		return nil, errors.New(errMsg)
	}

	var r res.ShortenResponse
	r.FromDomain(newSh)

	return &r, nil
}

func (s *shService) Delete(ctx context.Context, req *req.ShortenDelete) error {
	// check first if given id is exists in DB
	sh, err := s.repo.GetByID(ctx, req.ID, repo.Cols("id"))
	if err != nil {
		errMsg := "data with id " + strconv.Itoa(int(req.ID)) + " was not found"
		s.log.Err(errMsg+":", err)
		return errors.New(errMsg)
	}

	// then delete it using the id from query
	if err = s.repo.DeleteByID(ctx, sh.ID); err != nil {
		errMsg := "failed to delete SNS data with id " + strconv.Itoa(int(req.ID))
		s.log.Err(errMsg+":", err)
		return errors.New(errMsg)
	}
	return nil
}
