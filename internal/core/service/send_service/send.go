package send_service

import (
	"context"
	"errors"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/google/uuid"
	repo "github.com/mdanialr/sns_backend/internal/core/repository"
	"github.com/mdanialr/sns_backend/internal/core/repository/sns_repository"
	"github.com/mdanialr/sns_backend/internal/domain"
	req "github.com/mdanialr/sns_backend/internal/requests"
	res "github.com/mdanialr/sns_backend/internal/responses"
	h "github.com/mdanialr/sns_backend/pkg/helper"
	"github.com/mdanialr/sns_backend/pkg/logger"
	"github.com/mdanialr/sns_backend/pkg/storage"
	"github.com/spf13/viper"
)

type sendSvc struct {
	log  logger.Writer
	st   storage.IStorage
	v    *viper.Viper
	repo sns_repository.IRepository
}

// New return implementation of core business logic for Send service layer.
func New(l logger.Writer, s storage.IStorage, v *viper.Viper, r sns_repository.IRepository) IService {
	return &sendSvc{l, s, v, r}
}

func (s *sendSvc) Index(ctx context.Context, sn *req.Send) (*res.SendIndexResponse, error) {
	shortens, err := s.repo.FindSend(ctx, repo.Paginate(&sn.M), repo.Order(sn.Order+" "+sn.Sort))
	if err != nil {
		errMsg := "failed to retrieve all send data"
		s.log.Err(errMsg+":", err)
		return nil, errors.New(errMsg)
	}

	r := &res.SendIndexResponse{Pagination: &sn.M}
	r.Pagination.Paginate()
	r.FromDomain(shortens)

	return r, nil
}

func (s *sendSvc) Create(ctx context.Context, req *req.Send) (*res.SendResponse, error) {
	// make sure given url is not used yet in the db
	o, _ := s.repo.GetByUrl(ctx, req.Url, repo.Cols("id"))
	if o.ID != 0 {
		return nil, errors.New("url already been taken")
	}

	// save multipart to Storage
	fn, err := s.saveFile(req.Send)
	if err != nil {
		errMsg := "failed to save uploaded file"
		s.log.Err(errMsg+":", err)
		return nil, errors.New(errMsg)
	}

	// prepare new object to be saved to DB
	sn := &domain.SNS{
		Url:         req.Url,
		Send:        &fn,
		FileSize:    h.Ptr(h.BytesToHumanize(req.Send.Size)),
		IsPermanent: h.Ptr(req.PermanentToBool()),
	}
	if _, err = s.repo.Create(ctx, sn); err != nil {
		errMsg := "failed to create new Send"
		s.log.Err(errMsg+":", err)
		return nil, errors.New(errMsg)
	}

	// adapt data from domain.SNS to required SendResponse
	var r res.SendResponse
	r.FromDomain(sn)

	return &r, nil
}

func (s *sendSvc) Update(ctx context.Context, req *req.SendUpdate) (*res.SendResponse, error) {
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
	sn := &domain.SNS{
		ID:          req.ID,
		Url:         req.Url,
		IsPermanent: h.Ptr(req.PermanentToBool()),
	}
	newSn, err := s.repo.Update(ctx, sn)
	if err != nil {
		errMsg := "failed to update Send with id " + strconv.Itoa(int(req.ID))
		s.log.Err(errMsg+":", err)
		return nil, errors.New(errMsg)
	}

	var r res.SendResponse
	r.FromDomain(newSn)

	return &r, nil
}

func (s *sendSvc) Delete(ctx context.Context, req *req.SendDelete) error {
	// check first if given id is exists in DB
	sn, err := s.repo.GetByID(ctx, req.ID, repo.Cols("id", "send"))
	if err != nil {
		errMsg := "data with id " + strconv.Itoa(int(req.ID)) + " was not found"
		s.log.Err(errMsg+":", err)
		return errors.New(errMsg)
	}

	// then delete it using the id from query
	if err = s.repo.DeleteByID(ctx, sn.ID); err != nil {
		errMsg := "failed to delete SNS data with id " + strconv.Itoa(int(req.ID))
		s.log.Err(errMsg+":", err)
		return errors.New(errMsg)
	}

	// then delete the file
	go s.removeFile(*sn.Send)

	return nil
}

// saveFile save given multipart to Storage after prepend it with target file
// path from config and random string as the filename.
func (s *sendSvc) saveFile(f *multipart.FileHeader) (string, error) {
	fl, err := f.Open()
	if err != nil {
		return "", err
	}
	defer fl.Close()

	// set up the target path
	pt := strings.TrimSuffix(s.v.GetString("storage.path"), "/") + "/" // make sure to manually append slice
	// generate random name then append it with the file extension
	fn := uuid.NewString() + filepath.Ext(f.Filename)

	// save using separate goroutine
	go s.st.Save(fl, pt+fn)

	return fn, nil
}

// removeFile delete given filename after append it with Storage path from
// config.
func (s *sendSvc) removeFile(fn string) {
	pt := strings.TrimSuffix(s.v.GetString("storage.path"), "/") + "/" // make sure to manually append slice
	s.st.Remove(pt + fn)
}
