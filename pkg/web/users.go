package web

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/hearky/server/pkg/domain"
	"time"
)

func (s *Server) HandleCreateUser(c *fiber.Ctx) error {
	uid, err := s.Authorize(c)
	if err != nil || uid != "" {
		return nil
	}
	var dto domain.CreateUserDto
	err = c.BodyParser(&dto)
	if err != nil {
		return s.BadRequest(c)
	}

	ctx, ccl := context.WithTimeout(context.Background(), 10 * time.Second)
	defer ccl()
	err = s.userService.CreateUser(ctx, &dto, uid)
	if err != nil {
		return s.DomainError(c, err)
	}
	return c.SendStatus(fiber.StatusOK)
}

func (s *Server) HandleGetUserMeetings(c *fiber.Ctx) error {
	uid, err := s.Authorize(c)
	if err != nil || uid == "" {
		return nil
	}

	ctx, ccl := context.WithTimeout(context.Background(), 10 * time.Second)
	defer ccl()
	m, err := s.meetingService.GetMeetingsByUserID(ctx, uid)
	if err != nil {
		return s.DomainError(c, err)
	}

	return c.JSON(m)
}

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

	ctx, ccl := context.WithTimeout(context.Background(), 10 * time.Second)
	defer ccl()
	m, err := s.meetingService.GetMeetingByID(ctx, dto.MeetingID)
	if err != nil {
		return s.DomainError(c, err)
	}

	// Check if the user is organizer
	if !m.IsOrganizer(uid) {
		return s.Forbidden(c)
	}

	return c.JSON(m)
}