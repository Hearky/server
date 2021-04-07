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

package api

import (
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/websocket/v2"
	"go.uber.org/zap"
)

type WSMessage struct {
	OP   string      `json:"op"`
	Data interface{} `json:"d"`
}

type Error struct {
	Code string
}

var ErrInvalidPayload = Error{Code: "invalid-payload"}

func HandleMessage(c *websocket.Conn, m *WSMessage) {
	switch m.OP {
	case "login":
		login(c, m)
		break
	}
}

func sendError(c *websocket.Conn, err Error) {
	if err := c.WriteJSON(err); err != nil {
		sentry.CaptureException(err)
		zap.L().Error("could not send response", zap.Error(err))
	}
	return
}
