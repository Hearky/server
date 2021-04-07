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
	userRepository := user.NewRepository(db)
	inviteRepository := invite.NewRepository(db)

	meetingService := meeting.NewService(meetingRepository, inviteRepository, userRepository)
	userService := user.NewService(userRepository, meetingRepository, inviteRepository)
	inviteService := invite.NewService(inviteRepository, meetingRepository, userRepository)

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
