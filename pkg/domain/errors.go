package domain

import "errors"

var (
	ErrForbidden      = errors.New("forbidden")
	ErrInternal       = errors.New("internal")
	ErrNotFound       = errors.New("not-found")
	ErrUsernameExists = errors.New("username-already-exists")
	ErrInviteExists   = errors.New("invite-already-exists")
)
