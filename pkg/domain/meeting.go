package domain

import "context"

// CreateMeetingDto represents the needed data to create a new meeting
type CreateMeetingDto struct {
	Name         string   `json:"name"`
	Organizers   []string `json:"organizers"`
	Participants []string `json:"participants"`
}

// Meeting represents a Hearky meeting
type Meeting struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	Owner        string         `json:"owner"`
	Organizers   []string       `json:"organizers"`
	Participants []string       `json:"participants"`
	Upgrade      MeetingUpgrade `json:"upgrade"`
}

// MeetingUpgrade contains the upgrade data
type MeetingUpgrade struct {
	Participants      int `json:"participants"`
	ConcurrentInvites int `json:"concurrent_invites"`
}

// PartialMeeting is a subset with necessary data of a meeting
type PartialMeeting struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// MeetingRepository defines an interface for managing meetings in the database
type MeetingRepository interface {
	CreateMeeting(ctx context.Context, m *Meeting) error
	SaveMeeting(ctx context.Context, m *Meeting) error
	GetMeetingByID(ctx context.Context, id string) (*Meeting, error)
	GetMeetingsByUserID(ctx context.Context, id string) ([]*Meeting, error)
	DeleteMeetingByID(ctx context.Context, id string) error
}

// AsPartial returns a subset of a Meeting with only the necessary data
func (m *Meeting) AsPartial() *PartialMeeting {
	return &PartialMeeting{
		ID:   m.ID,
		Name: m.Name,
	}
}

// IsOwner returns true if the passed user is the owner
func (m *Meeting) IsOwner(uid string) bool {
	return m.Owner == uid
}

// IsParticipant returns true if the passed user is an organizer
func (m *Meeting) IsOrganizer(uid string) bool {
	if m.Owner == uid {
		return true
	}
	for _, o := range m.Organizers {
		if o == uid {
			return true
		}
	}
	return false
}

// IsParticipant returns true if the passed user is a participant
func (m *Meeting) IsParticipant(uid string) bool {
	for _, p := range m.Participants {
		if p == uid {
			return true
		}
	}
	return m.IsOrganizer(uid)
}

// AddOrganizer adds a new organizer to the meeting
func (m *Meeting) AddOrganizer(uid string) {
	m.Organizers = append(m.Organizers, uid)
}

// AddParticipant adds a new participant to the meeting
func (m *Meeting) AddParticipant(uid string) {
	m.Participants = append(m.Participants, uid)
}
