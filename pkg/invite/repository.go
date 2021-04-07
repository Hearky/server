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

package invite

import (
	"context"
	"github.com/getsentry/sentry-go"
	"github.com/hearky/server/pkg/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type repository struct {
	col *mongo.Collection
}

func NewRepository(db *mongo.Database) domain.InviteRepository {
	return &repository{
		col: db.Collection("invites"),
	}
}

func (r *repository) CreateInvite(ctx context.Context, i *domain.Invite) error {
	res, err := r.col.InsertOne(ctx, i)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to insert invite", zap.Any("invite", i), zap.Error(err))
		return domain.ErrInternal
	}
	zap.L().Info("inserted new invite", zap.Any("id", res.InsertedID))
	return nil
}

func (r *repository) GetInvitesByReceiver(ctx context.Context, id string) ([]*domain.Invite, error) {
	var i []*domain.Invite
	c, err := r.col.Find(ctx, bson.D{primitive.E{Key: "receiver_id", Value: id}})
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to find invites by receiver", zap.String("id", id), zap.Error(err))
		return nil, domain.ErrInternal
	}
	err = c.All(ctx, &i)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to parse invite elements from cursor in slice", zap.Error(err))
		return nil, domain.ErrInternal
	}
	if i == nil {
		i = make([]*domain.Invite, 0)
	}
	return i, nil
}

func (r *repository) GetInviteByReceiverAndMeeting(ctx context.Context, id string, mId string) (*domain.Invite, error) {
	var i domain.Invite
	err := r.col.FindOne(ctx, bson.D{primitive.E{Key: "receiver_id", Value: id}, primitive.E{Key: "meeting_id", Value: mId}}).Decode(&i)
	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrNotFound
	} else if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to find invite by receiver and meeting", zap.String("id", id), zap.Error(err))
		return nil, domain.ErrInternal
	}
	return &i, nil
}

func (r *repository) GetInviteByID(ctx context.Context, id string) (*domain.Invite, error) {
	var i domain.Invite
	err := r.col.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&i)
	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrNotFound
	} else if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to find invite by id", zap.String("id", id), zap.Error(err))
		return nil, domain.ErrInternal
	}
	return &i, nil
}

func (r *repository) GetInvitesByReceiverCount(ctx context.Context, uid string) (int64, error) {
	c, err := r.col.CountDocuments(ctx, bson.D{primitive.E{Key: "receiver_id", Value: uid}})
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to find invite count by user", zap.String("id", uid), zap.Error(err))
		return -1, domain.ErrInternal
	}
	return c, nil
}

func (r *repository) GetInvitesByMeeting(ctx context.Context, mid string) ([]*domain.Invite, error) {
	var i []*domain.Invite
	c, err := r.col.Find(ctx, bson.D{primitive.E{Key: "meeting_id", Value: mid}})
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to find invites by meeting", zap.String("id", mid), zap.Error(err))
		return nil, domain.ErrInternal
	}
	err = c.All(ctx, &i)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to parse invite elements from cursor in slice", zap.Error(err))
		return nil, domain.ErrInternal
	}
	if i == nil {
		i = make([]*domain.Invite, 0)
	}
	return i, nil
}

func (r *repository) GetInvitesByMeetingCount(ctx context.Context, mid string) (int64, error) {
	c, err := r.col.CountDocuments(ctx, bson.D{primitive.E{Key: "meeting_id", Value: mid}})
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to find invite count by meeting", zap.String("id", mid), zap.Error(err))
		return -1, domain.ErrInternal
	}
	return c, nil
}

func (r *repository) DeleteInvite(ctx context.Context, id string) error {
	_, err := r.col.DeleteOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}})
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to delete invite", zap.Any("id", id), zap.Error(err))
		return domain.ErrInternal
	}
	zap.L().Info("deleted invite", zap.String("id", id))
	return nil
}
