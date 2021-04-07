package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hearky/server/pkg/domain"
	"go.uber.org/zap"
)

type errorMsg struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (s *Server) Error(c *fiber.Ctx, code int, m string) error {
	return c.Status(code).JSON(&errorMsg{Code: code, Message: m})
}

func (s *Server) DomainError(c *fiber.Ctx, err error) error {
	switch err {
	case domain.ErrTooManyMeetings:
		return s.Error(c, fiber.StatusBadRequest, err.Error())
	case domain.ErrOwnerOfMeeting:
		return s.Error(c, fiber.StatusBadRequest, err.Error())
	case domain.ErrUsernameExists:
		return s.Error(c, fiber.StatusBadRequest, err.Error())
	case domain.ErrUserExists:
		return s.Error(c, fiber.StatusBadRequest, err.Error())
	case domain.ErrInviteExists:
		return s.Error(c, fiber.StatusBadRequest, err.Error())
	case domain.ErrInternal:
		return s.InternalError(c)
	case domain.ErrNotFound:
		return s.Error(c, fiber.StatusNotFound, err.Error())
	case domain.ErrForbidden:
		return s.Forbidden(c)

	}
	zap.L().Error("caught an unknown error", zap.Error(err))
	return s.InternalError(c)
}

func (s *Server) Forbidden(c *fiber.Ctx) error {
	return s.Error(c, fiber.StatusForbidden, "forbidden")
}

func (s *Server) InternalError(c *fiber.Ctx) error {
	return s.Error(c, fiber.StatusInternalServerError, "internal")
}

func (s *Server) BadRequest(c *fiber.Ctx) error {
	return s.Error(c, fiber.StatusBadRequest, "bad-request")
}

func (s *Server) Unauthorized(c *fiber.Ctx, m string) error {
	return s.Error(c, fiber.StatusUnauthorized, m)
}
