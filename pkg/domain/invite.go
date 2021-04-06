package domain

import (
	"context"
	"time"
)

type Invite struct {
	ID         string    `json:"id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	MeetingID  string    `json:"meeting_id"`
	Timestamp  time.Time `json:"timestamp"`
}

type CreateInviteDto struct {
	ReceiverID string `json:"receiver_id"`
	MeetingID  string `json:"meeting_id"`
}

type InviteRepository interface {
	CreateInvite(ctx context.Context, i *Invite) error
	GetInviteByID(ctx context.Context, id string) (*Invite, error)
	GetInvitesByReceiver(ctx context.Context, uid string) ([]*Invite, error)
	GetInviteByReceiverAndMeeting(ctx context.Context, uid string, mid string) (*Invite, error)
	DeleteInvite(ctx context.Context, id string) error
}

type InviteService interface {
	SendInvite(dto *CreateInviteDto, uid string) error
	GetInvitesByReceiver(uid string) error
	AcceptInvite(id string) error
	RejectInvite(id string) error
}
