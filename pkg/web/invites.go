/*
 * Hearky Server
 * Copyright (C) 2021 Hearky
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

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
