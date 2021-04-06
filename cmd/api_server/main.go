package main

import (
	"context"
	firebase "firebase.google.com/go/v4"
	"github.com/getsentry/sentry-go"
	"github.com/hearky/server/pkg/config"
	"github.com/hearky/server/pkg/invite"
	"github.com/hearky/server/pkg/logger"
	"github.com/hearky/server/pkg/meeting"
	"github.com/hearky/server/pkg/user"
	"github.com/hearky/server/pkg/web"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"time"
)

func main() {
	cfg := config.Load()
	logger.Initialize(cfg.Dev)

	// Connect to to MongoDB
	ctx, ccl := context.WithTimeout(context.Background(), 10*time.Second)
	defer ccl()
	mgClientOpts := options.Client().ApplyURI(cfg.MongoURI)
	mgClient, err := mongo.Connect(ctx, mgClientOpts)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Fatal("failed to connect to MongoDB", zap.Error(err))
	}
	db := mgClient.Database(cfg.MongoDBName)

	meetingRepository := meeting.NewRepository(db)
	meetingService := meeting.NewService(meetingRepository)

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)

	inviteRepository := invite.NewRepository(db)
	inviteService := invite.NewService(inviteRepository, meetingRepository)

	// Initialize firebase
	ctx, ccl = context.WithTimeout(context.Background(), 10*time.Second)
	defer ccl()
	fbOpts := option.WithCredentialsFile("serviceAccountKey.json")
	fbApp, err := firebase.NewApp(ctx, nil, fbOpts)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Fatal("failed to initialize Firebase", zap.Error(err))
	}
	fbAuth, err := fbApp.Auth(context.Background())
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Fatal("failed to create FirebaseAuth client", zap.Error(err))
	}

	// Initialize and start server
	s := web.New(cfg.Dev, fbAuth, userService, meetingService, inviteService)
	s.Start(cfg.WebAddress)
}
