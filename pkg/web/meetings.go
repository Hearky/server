package web

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/hearky/server/pkg/domain"
	"time"
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

	ctx, ccl := context.WithTimeout(context.Background(), 10*time.Second)
	defer ccl()
	mId, err := s.meetingService.CreateMeeting(ctx, &dto, uid)
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

	ctx, ccl := context.WithTimeout(context.Background(), 10*time.Second)
	defer ccl()
	m, err := s.meetingService.GetMeetingByID(ctx, mId)
	if err != nil {
		return s.DomainError(c, err)
	}
	if !m.IsParticipant(uid) {
		return s.Forbidden(c)
	}
	return c.JSON(m)
}

func (s *Server) HandleDeleteMeetingByID(c *fiber.Ctx) error {
	uid, err := s.Authorize(c)
	if err != nil || uid == "" {
		return nil
	}
	mId := c.Params("id")

	ctx, ccl := context.WithTimeout(context.Background(), 10*time.Second)
	defer ccl()
	m, err := s.meetingService.GetMeetingByID(ctx, mId)
	if err != nil {
		return s.DomainError(c, err)
	}
	if m.Owner != uid {
		return s.Forbidden(c)
	}

	err = s.meetingService.DeleteMeeting(ctx, mId)
	if err != nil {
		return s.DomainError(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}
