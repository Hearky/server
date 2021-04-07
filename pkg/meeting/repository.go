package meeting

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

func NewRepository(db *mongo.Database) domain.MeetingRepository {
	return &repository{
		col: db.Collection("meetings"),
	}
}

func (r *repository) CreateMeeting(ctx context.Context, m *domain.Meeting) error {
	res, err := r.col.InsertOne(ctx, m)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to insert meeting", zap.Any("meeting", m), zap.Error(err))
		return domain.ErrInternal
	}
	zap.L().Info("inserted new meeting", zap.Any("id", res.InsertedID))
	return nil
}

func (r *repository) GetMeetingByID(ctx context.Context, id string) (*domain.Meeting, error) {
	var m domain.Meeting
	err := r.col.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}}).Decode(&m)
	if err == mongo.ErrNoDocuments {
		return nil, domain.ErrNotFound
	} else if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to find meeting by id", zap.String("id", id), zap.Error(err))
		return nil, domain.ErrInternal
	}
	return &m, nil
}

func (r *repository) GetMeetingsByUser(ctx context.Context, id string) ([]*domain.Meeting, error) {
	var m []*domain.Meeting
	c, err := r.col.Find(ctx, bson.M{"$or": []bson.M{
		{"owner_id": id},
		{"organizers": id},
		{"participants": id},
	}})
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to find meetings by user", zap.String("id", id), zap.Error(err))
		return nil, domain.ErrInternal
	}
	err = c.All(ctx, &m)
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to parse meeting elements from cursor in slice", zap.Error(err))
		return nil, domain.ErrInternal
	}
	if m == nil {
		m = make([]*domain.Meeting, 0)
	}
	return m, nil
}

func (r *repository) GetMeetingsByUserCount(ctx context.Context, id string) (int64, error) {
	c, err := r.col.CountDocuments(ctx, bson.M{"$or": []bson.M{
		{"owner_id": id},
		{"organizers": id},
		{"participants": id},
	}})
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to find meetings count by user", zap.String("id", id), zap.Error(err))
		return -1, domain.ErrInternal
	}
	return c, nil
}

func (r *repository) DeleteMeeting(ctx context.Context, id string) error {
	_, err := r.col.DeleteOne(ctx, bson.D{primitive.E{Key: "_id", Value: id}})
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to delete meeting", zap.Any("id", id), zap.Error(err))
		return domain.ErrInternal
	}
	zap.L().Info("deleted meeting", zap.String("id", id))
	return nil
}

func (r *repository) SaveMeeting(ctx context.Context, m *domain.Meeting) error {
	_, err := r.col.UpdateByID(ctx, m.ID, bson.D{
		{"$set", bson.D{
			{"name", m.Name},
			{"owner_id", m.OwnerID},
			{"participants", m.Participants},
			{"organizers", m.Organizers},
			{"upgrade", m.Upgrade},
		}},
	})
	if err != nil {
		sentry.CaptureException(err)
		zap.L().Error("failed to update meeting", zap.Any("meeting", m), zap.Error(err))
		return domain.ErrInternal
	}
	return nil
}
