package web

import (
	"firebase.google.com/go/v4/auth"
	"github.com/ansrivas/fiberprometheus/v2"
	"github.com/getsentry/sentry-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/hearky/server/pkg/domain"
	"go.uber.org/zap"
)

type Server struct {
	app            *fiber.App
	fbAuth         *auth.Client
	userService    domain.UserService
	meetingService domain.MeetingService
	inviteService  domain.InviteService
}

func New(dev bool, fbAuth *auth.Client, userService domain.UserService, meetingService domain.MeetingService, inviteService domain.InviteService) *Server {
	app := fiber.New()

	s := &Server{
		app:            app,
		fbAuth:         fbAuth,
		userService:    userService,
		meetingService: meetingService,
		inviteService:  inviteService,
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
	api.Delete("/meetings/:id", s.HandleDeleteMeeting)
	api.Get("/meetings/:id/invites", s.HandleGetMeetingInvites)
	api.Get("/meetings/:id/invites/count", s.HandleGetMeetingInvitesCount)

	api.Post("/users", s.HandleCreateUser)
	api.Get("/users/@me", s.HandleGetMe)
	api.Delete("/users/@me", s.HandleDeleteMe)
	api.Get("/users/@me/meetings", s.HandleGetMyMeetings)
	api.Get("/users/@me/meetings/count", s.HandleGetMyMeetingsCount)
	api.Get("/users/@me/invites", s.HandleGetMyInvites)
	api.Get("/users/@me/invites/count", s.HandleGetMyInvitesCount)

	api.Post("/invites", s.HandleSendInvite)
	api.Post("/invites/:id/accept", s.HandleAcceptInvite)
	api.Delete("/invites/:id", s.HandleDeleteInvite)
	return s
}

func (s *Server) Start(addr string) {
	err := s.app.Listen(addr)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Fatal("failed to serve webserver", zap.String("address", addr), zap.Error(err))
	}
}
