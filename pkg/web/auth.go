package web

import (
	"context"
	"github.com/gofiber/fiber/v2"
	"strings"
	"time"
)

func (s *Server) Authorize(c *fiber.Ctx) (string, error) {
	uid, err := s.TokenAuth(c)
	if err != nil {
		return "", err
	}

	_, err = s.userService.GetUser(uid, uid)
	if err != nil {
		return "", s.DomainError(c, err)
	}
	return uid, nil
}

func (s *Server) TokenAuth(c *fiber.Ctx) (string, error) {
	rawToken := c.Get(fiber.HeaderAuthorization)

	// Check if token even exists
	if rawToken == "" {
		return "", s.Unauthorized(c, "missing-header")
	}

	// Split token in parts and check if length eq 2
	tokenParts := strings.Split(rawToken, " ")
	if len(tokenParts) != 2 {
		return "", s.Unauthorized(c, "invalid-header")
	}

	// Check if first part is Bearer
	if tokenParts[0] != "Bearer" {
		return "", s.Unauthorized(c, "invalid-header")
	}

	// Check with firebase
	ctx, ccl := context.WithTimeout(context.Background(), 5*time.Second)
	defer ccl()
	t, err := s.fbAuth.VerifyIDTokenAndCheckRevoked(ctx, tokenParts[1])
	if err != nil {
		return "", s.Unauthorized(c, "invalid-token")
	}

	return t.UID, nil
}
