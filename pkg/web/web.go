package web

import (
	"firebase.google.com/go/v4/auth"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/hearky/server/pkg/invite"
	"github.com/hearky/server/pkg/meeting"
	"github.com/hearky/server/pkg/user"
	"go.uber.org/zap"
)

type Server struct {
	app            *fiber.App
	fbAuth         *auth.Client
	userService    *user.Service
	meetingService *meeting.Service
	inviteService  *invite.service
}

func New(dev bool, fbAuth *auth.Client, userService *user.Service, meetingService *meeting.Service, inviteService *invite.service) *Server {
	app := fiber.New()

	s := &Server{
		app:            app,
		fbAuth:         fbAuth,
		userService:    userService,
		meetingService: meetingService,
	}

	// Register metrics endpoint for prometheus scraping
	prometheus := fiberprometheus.New("hearky_server")
	prometheus.RegisterAt(app, "/metrics")
	app.Use(prometheus.Middleware)

	// Monitoring Dashboard
	app.Get("/dashboard", monitor.New())

	// Register API routes
	api := app.Group("/api")
	api.Post("/meetings", s.HandleCreateMeeting)
	api.Get("/meetings/:id", s.HandleGetMeetingByID)
	api.Delete("/meetings/:id", s.HandleDeleteMeetingByID)

	api.Post("/users", s.HandleCreateUser)
	api.Get("/users/@me/meetings", s.HandleGetUserMeetings)
	return s
}

func (s *Server) Start(addr string) {
	err := s.app.Listen(addr)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Fatal("failed to serve webserver", zap.String("address", addr), zap.Error(err))
	}
}
