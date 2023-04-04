package shorten_handler

import (
	"github.com/gofiber/fiber/v2"
	md "github.com/mdanialr/sns_backend/internal/app/adapter/http/middleware"
	"github.com/mdanialr/sns_backend/internal/core/service/shorten_service"
	"github.com/mdanialr/sns_backend/internal/requests"
	cons "github.com/mdanialr/sns_backend/pkg/constant"
	resp "github.com/mdanialr/sns_backend/pkg/response"
	"github.com/spf13/viper"
)

type shortenHandler struct {
	v     *viper.Viper
	route fiber.Router
	shSvc shorten_service.IService
}

// NewShortenHandler init all endpoints within `/shorten`.
func NewShortenHandler(route fiber.Router, v *viper.Viper, svc shorten_service.IService) {
	sh := &shortenHandler{v, route, svc}

	api := sh.route.Group("/shorten", md.JWT(sh.v))
	api.Get("/", sh.Index)
	api.Post("/create", sh.Create)
	api.Post("/update", sh.Update)
	api.Post("/delete", sh.Delete)
}

// Index retrieve all data in shorten category.
func (s *shortenHandler) Index(c *fiber.Ctx) error {
	req := new(requests.Shorten)
	c.QueryParser(req)
	// set up the query order and sort
	req.SetQuery()

	res, err := s.shSvc.Index(c.Context(), req)
	if err != nil {
		return resp.Error(c, resp.WithErr(err))
	}

	return resp.Success(c, resp.WithData(res.Data), resp.WithMeta(res.Pagination))
}

// Create save new Shorten instance to DB.
func (s *shortenHandler) Create(c *fiber.Ctx) error {
	req := new(requests.Shorten)
	c.BodyParser(req)

	// validate the request
	if err := req.Validate(); err != nil {
		return resp.Error(c, resp.WithErrMsg(cons.InvalidPayload), resp.WithErrValidation(err))
	}

	res, err := s.shSvc.Create(c.Context(), req)
	if err != nil {
		return resp.Error(c, resp.WithErr(err))
	}

	return resp.Success(c, resp.WithData(res))
}

// Update do update an existing Shorten instance in DB.
func (s *shortenHandler) Update(c *fiber.Ctx) error {
	req := new(requests.ShortenUpdate)
	c.BodyParser(req)

	// validate the request
	if err := req.Validate(); err != nil {
		return resp.Error(c, resp.WithErrMsg(cons.InvalidPayload), resp.WithErrValidation(err))
	}

	res, err := s.shSvc.Update(c.Context(), req)
	if err != nil {
		return resp.Error(c, resp.WithErr(err))
	}

	return resp.Success(c, resp.WithData(res))
}

// Delete remove a Shorten instance from DB.
func (s *shortenHandler) Delete(c *fiber.Ctx) error {
	req := new(requests.ShortenDelete)
	c.BodyParser(req)

	// validate the request
	if err := req.Validate(); err != nil {
		return resp.Error(c, resp.WithErrMsg(cons.InvalidPayload), resp.WithErrValidation(err))
	}

	if err := s.shSvc.Delete(c.Context(), req); err != nil {
		return resp.Error(c, resp.WithErr(err))
	}

	return resp.Success(c)
}
