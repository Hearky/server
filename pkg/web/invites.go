package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hearky/server/pkg/domain"
)

func (s *Server) HandleSendInvite(c *fiber.Ctx) error {
	uid, err := s.Authorize(c)
	if err != nil || uid == "" {
		return nil
	}

	var dto domain.CreateInviteDto
	err = c.BodyParser(&dto)
	if err != nil {
		return s.BadRequest(c)
	}

	err = s.inviteService.SendInvite(&dto, uid)
	if err != nil {
		return s.DomainError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (s *Server) HandleAcceptInvite(c *fiber.Ctx) error {
	uid, err := s.Authorize(c)
	if err != nil || uid == "" {
		return nil
	}
	id := c.Params("id")

	err = s.inviteService.AcceptInvite(id, uid)
	if err != nil {
		return s.DomainError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (s *Server) HandleDeleteInvite(c *fiber.Ctx) error {
	uid, err := s.Authorize(c)
	if err != nil || uid == "" {
		return nil
	}
	id := c.Params("id")

	err = s.inviteService.DeleteInvite(id, uid)
	if err != nil {
		return s.DomainError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}
