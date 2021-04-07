package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hearky/server/pkg/domain"
)

func (s *Server) HandleCreateUser(c *fiber.Ctx) error {
	uid, err := s.TokenAuth(c)
	if err != nil || uid == "" {
		return nil
	}
	var dto domain.CreateUserDto
	err = c.BodyParser(&dto)
	if err != nil {
		return s.BadRequest(c)
	}

	err = s.userService.CreateUser(&dto, uid)
	if err != nil {
		return s.DomainError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (s *Server) HandleGetMe(c *fiber.Ctx) error {
	uid, err := s.Authorize(c)
	if err != nil || uid == "" {
		return nil
	}

	u, err := s.userService.GetUser(uid, uid)
	if err != nil {
		return s.DomainError(c, err)
	}

	return c.JSON(u)
}

func (s *Server) HandleDeleteMe(c *fiber.Ctx) error {
	uid, err := s.Authorize(c)
	if err != nil || uid == "" {
		return nil
	}

	err = s.userService.DeleteUser(uid, uid)
	if err != nil {
		return s.DomainError(c, err)
	}

	return c.SendStatus(fiber.StatusOK)
}

func (s *Server) HandleGetMyMeetings(c *fiber.Ctx) error {
	uid, err := s.Authorize(c)
	if err != nil || uid == "" {
		return nil
	}

	m, err := s.meetingService.GetMeetingsByUser(uid)
	if err != nil {
		return s.DomainError(c, err)
	}

	return c.JSON(m)
}

func (s *Server) HandleGetMyMeetingsCount(c *fiber.Ctx) error {
	uid, err := s.Authorize(c)
	if err != nil || uid == "" {
		return nil
	}

	count, err := s.meetingService.GetMeetingsByUserCount(uid)
	if err != nil {
		return s.DomainError(c, err)
	}

	return c.JSON(&domain.CountMessage{Count: count})
}

func (s *Server) HandleGetMyInvites(c *fiber.Ctx) error {
	uid, err := s.Authorize(c)
	if err != nil || uid == "" {
		return nil
	}

	m, err := s.inviteService.GetInvitesByReceiver(uid)
	if err != nil {
		return s.DomainError(c, err)
	}

	return c.JSON(m)
}

func (s *Server) HandleGetMyInvitesCount(c *fiber.Ctx) error {
	uid, err := s.Authorize(c)
	if err != nil || uid == "" {
		return nil
	}

	count, err := s.inviteService.GetInvitesByReceiverCount(uid)
	if err != nil {
		return s.DomainError(c, err)
	}

	return c.JSON(&domain.CountMessage{Count: count})
}
