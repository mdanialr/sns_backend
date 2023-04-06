package send_handler

import (
	"github.com/gofiber/fiber/v2"
	md "github.com/mdanialr/sns_backend/internal/app/adapter/http/middleware"
	"github.com/mdanialr/sns_backend/internal/core/service/send_service"
	"github.com/mdanialr/sns_backend/internal/requests"
	cons "github.com/mdanialr/sns_backend/pkg/constant"
	resp "github.com/mdanialr/sns_backend/pkg/response"
	"github.com/spf13/viper"
)

type sendHandler struct {
	v     *viper.Viper
	route fiber.Router
	svc   send_service.IService
}

// New init all endpoints within `/send`.
func New(r fiber.Router, v *viper.Viper, svc send_service.IService) {
	sn := &sendHandler{v, r, svc}

	api := sn.route.Group("/send", md.JWT(sn.v))
	api.Get("/", sn.Index)
	api.Post("/create", sn.Create)
	api.Post("/update", sn.Update)
	api.Post("/delete", sn.Delete)
}

func (s *sendHandler) Index(c *fiber.Ctx) error {
	req := new(requests.Send)
	c.QueryParser(req)
	// set up the query order and sort
	req.SetQuery()

	res, err := s.svc.Index(c.Context(), req)
	if err != nil {
		return resp.Error(c, resp.WithErr(err))
	}

	return resp.Success(c, resp.WithData(res.Data), resp.WithMeta(res.Pagination))
}

func (s *sendHandler) Create(c *fiber.Ctx) error {
	req := new(requests.Send)
	c.BodyParser(req)
	// manually retrieve binary file for 'send' param
	req.Send, _ = c.FormFile("send")

	// validate the request
	if err := req.Validate(); err != nil {
		return resp.Error(c, resp.WithErrMsg(cons.InvalidPayload), resp.WithErrValidation(err))
	}

	res, err := s.svc.Create(c.Context(), req)
	if err != nil {
		return resp.Error(c, resp.WithErr(err))
	}

	return resp.Success(c, resp.WithData(res))
}

func (s *sendHandler) Update(c *fiber.Ctx) error {
	req := new(requests.SendUpdate)
	c.BodyParser(req)

	// validate the request
	if err := req.Validate(); err != nil {
		return resp.Error(c, resp.WithErrMsg(cons.InvalidPayload), resp.WithErrValidation(err))
	}

	res, err := s.svc.Update(c.Context(), req)
	if err != nil {
		return resp.Error(c, resp.WithErr(err))
	}

	return resp.Success(c, resp.WithData(res))
}

func (s *sendHandler) Delete(c *fiber.Ctx) error {
	req := new(requests.SendDelete)
	c.BodyParser(req)

	// validate the request
	if err := req.Validate(); err != nil {
		return resp.Error(c, resp.WithErrMsg(cons.InvalidPayload), resp.WithErrValidation(err))
	}

	if err := s.svc.Delete(c.Context(), req); err != nil {
		return resp.Error(c, resp.WithErr(err))
	}

	return resp.Success(c)
}
