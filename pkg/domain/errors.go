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

import "errors"

var (
	ErrForbidden      = errors.New("forbidden")
	ErrInternal       = errors.New("internal")
	ErrNotFound       = errors.New("not-found")
	ErrUsernameExists = errors.New("username-already-exists")
	ErrUserExists     = errors.New("user-already-exists")
	ErrInviteExists   = errors.New("invite-already-exists")
	ErrOwnerOfMeeting = errors.New("owner-of-meeting")
	ErrTooManyMeetings = errors.New("too-many-meetings")
)
