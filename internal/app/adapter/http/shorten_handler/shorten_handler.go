package shorten_handler

import (
	"github.com/gofiber/fiber/v2"
	md "github.com/mdanialr/sns_backend/internal/app/adapter/http/middleware"
	"github.com/mdanialr/sns_backend/internal/core/service/shorten_service"
	"github.com/mdanialr/sns_backend/internal/requests"
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
