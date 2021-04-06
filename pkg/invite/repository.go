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

type Repository struct {
	col *mongo.Collection
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{
		col: db.Collection("invites"),
	}
}

func (r *Repository) CreateInvite(ctx context.Context, i *domain.Invite) error {
	res, err := r.col.InsertOne(ctx, i)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to insert invite", zap.Any("invite", i), zap.Error(err))
		return domain.ErrInternal
	}
	zap.L().Info("inserted new invite", zap.Any("id", res.InsertedID))
	return nil
}

func (r *Repository) GetInvitesByReceiver(ctx context.Context, id string) ([]*domain.Invite, error) {
	var i []*domain.Invite
	c, err := r.col.Find(ctx, bson.D{primitive.E{Key: "receiver", Value: id}})
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

func (r *Repository) GetInviteByReceiverAndMeeting(ctx context.Context, id string, mId string) (*domain.Invite, error) {
	var i domain.Invite
	err := r.col.FindOne(ctx, bson.D{primitive.E{Key: "receiver", Value: id}, primitive.E{Key: "meeting", Value: mId}}).Decode(&i)
	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrNotFound
	} else if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to find invite by receiver and meeting", zap.String("id", id), zap.Error(err))
		return nil, domain.ErrInternal
	}
	return &i, nil
}
