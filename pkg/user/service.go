package user

import (
	"context"
	"github.com/hearky/server/pkg/domain"
)

type Service struct {
	userRepo *Repository
}

func NewService(r *Repository) *Service {
	return &Service{
		userRepo: r,
	}
}

func (s *Service) CreateUser(ctx context.Context, dto *domain.CreateUserDto, uid string) error {
	_, err := s.userRepo.GetUserByUsername(ctx, dto.Username)
	if err != nil && err != domain.ErrNotFound {
		return err
	} else if err != domain.ErrNotFound {
		return domain.ErrUsernameExists
	}
	u := &domain.User{
		ID:       uid,
		Username: dto.Username,
	}
	return s.userRepo.CreateUser(ctx, u)
}
