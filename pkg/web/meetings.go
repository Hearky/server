package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/hearky/server/pkg/domain"
)

func (s *Server) HandleCreateMeeting(c *fiber.Ctx) error {
	uid, err := s.Authorize(c)
	if err != nil || uid == "" {
		return nil
	}
	var dto domain.CreateMeetingDto
	err = c.BodyParser(&dto)
	if err != nil {
		return s.BadRequest(c)
	}

	mId, err := s.meetingService.CreateMeeting(&dto, uid)
	if err != nil {
		return s.DomainError(c, err)
	}
	return c.JSON(&domain.IDMessage{ID: mId})
}

func (s *Server) HandleGetMeetingByID(c *fiber.Ctx) error {
	uid, err := s.Authorize(c)
	if err != nil || uid == "" {
		return nil
	}
	mId := c.Params("id")

	m, err := s.meetingService.GetMeetingByID(mId, uid)
	if err != nil {
		return s.DomainError(c, err)
	}
	return c.JSON(m)
}

func (s *Server) HandleDeleteMeeting(c *fiber.Ctx) error {
	uid, err := s.Authorize(c)
	if err != nil || uid == "" {
		return nil
	}
	mId := c.Params("id")

	err = s.meetingService.DeleteMeeting(mId, uid)
	if err != nil {
		return s.DomainError(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (s *Server) HandleGetMeetingInvites(c *fiber.Ctx) error {
	uid, err := s.Authorize(c)
	if err != nil || uid == "" {
		return nil
	}
	mId := c.Params("id")

	i, err := s.inviteService.GetInvitesByMeeting(mId, uid)
	if err != nil {
		return s.DomainError(c, err)
	}
	return c.JSON(i)
}

func (s *Server) HandleGetMeetingInvitesCount(c *fiber.Ctx) error {
	uid, err := s.Authorize(c)
	if err != nil || uid == "" {
		return nil
	}
	mId := c.Params("id")

	count, err := s.inviteService.GetInvitesByMeetingCount(mId, uid)
	if err != nil {
		return s.DomainError(c, err)
	}
	return c.JSON(&domain.CountMessage{Count: count})
}
