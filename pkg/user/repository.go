package user

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

func NewRepository(db *mongo.Database) domain.UserRepository {
	return &repository{
		col: db.Collection("users"),
	}
}

func (r *repository) CreateUser(ctx context.Context, u *domain.User) error {
	res, err := r.col.InsertOne(ctx, u)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to insert user", zap.Any("user", u), zap.Error(err))
		return domain.ErrInternal
	}
	zap.L().Info("inserted new user", zap.Any("id", res.InsertedID))
	return nil
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var u domain.User
	err := r.col.FindOne(ctx, bson.D{primitive.E{Key: "username", Value: username}}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrNotFound
	} else if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to find user by username", zap.String("username", username), zap.Error(err))
		return nil, domain.ErrInternal
	}
	return &u, nil
}

func (r *repository) GetUserByID(ctx context.Context, uid string) (*domain.User, error) {
	var u domain.User
	err := r.col.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: uid}}).Decode(&u)
	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrNotFound
	} else if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to find user by id", zap.String("id", uid), zap.Error(err))
		return nil, domain.ErrInternal
	}
	return &u, nil
}

func (r *repository) SaveUser(ctx context.Context, u *domain.User) error {
	_, err := r.col.UpdateByID(ctx, u.ID, bson.D{{
		"$set", bson.D{
			{"username", u.Username},
			{"upgrade", u.Upgrade},
		},
	}})
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to update user", zap.Any("user", u), zap.Error(err))
		return domain.ErrInternal
	}
	return nil
}

func (r *repository) DeleteUser(ctx context.Context, uid string) error {
	_, err := r.col.DeleteOne(ctx, bson.D{primitive.E{Key: "_id", Value: uid}})
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to delete user", zap.Any("id", uid), zap.Error(err))
		return domain.ErrInternal
	}
	return nil
}
