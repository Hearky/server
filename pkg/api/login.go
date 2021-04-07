package api

import (
	"github.com/gofiber/websocket/v2"
	"go.uber.org/zap"
)

type LoginRequest struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

func login(c *websocket.Conn, m *WSMessage) {
	zap.L().Info("mes", zap.Any("message", m.Data))
	req, ok := m.Data.(*LoginRequest)
	if !ok {
		sendError(c, ErrInvalidPayload)
		return
	}
	zap.L().Info("Login Request", zap.Any("req", req))
}
