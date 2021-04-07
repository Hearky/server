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

package domain

import "context"

// CreateMeetingDto represents the needed data to create a new meeting
type CreateMeetingDto struct {
	Name         string   `json:"name"`
	Participants []string `json:"participants"`
}

// Meeting represents a Hearky meeting
type Meeting struct {
	ID           string         `json:"id" bson:"_id"`
	Name         string         `json:"name"`
	OwnerID      string         `json:"owner_id" bson:"owner_id"`
	Organizers   []string       `json:"organizers"`
	Participants []string       `json:"participants"`
	Upgrade      MeetingUpgrade `json:"upgrade"`
}

// MeetingUpgrade contains the upgrade data
type MeetingUpgrade struct {
	Participants      int `json:"participants"`
	ConcurrentInvites int `json:"concurrent_invites" bson:"concurrent_invites"`
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
	GetMeetingsByUser(ctx context.Context, id string) ([]*Meeting, error)
	GetMeetingsByUserCount(ctx context.Context, id string) (int64, error)
	DeleteMeeting(ctx context.Context, id string) error
}

type MeetingService interface {
	CreateMeeting(dto *CreateMeetingDto, uid string) (string, error)
	GetMeetingByID(mid string, uid string) (*Meeting, error)
	GetMeetingsByUser(uid string) ([]*Meeting, error)
	GetMeetingsByUserCount(uid string) (int64, error)
	DeleteMeeting(mid string, uid string) error
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
	return m.OwnerID == uid
}

// IsParticipant returns true if the passed user is an organizer
func (m *Meeting) IsOrganizer(uid string) bool {
	if m.OwnerID == uid {
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
