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
