package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/thegenggo/equipment-loan/api/internal/model"
	"github.com/thegenggo/equipment-loan/api/internal/repository"
	"github.com/thegenggo/equipment-loan/api/pkg/token"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailTaken         = errors.New("email already registered")
	ErrInvalidCredentials = errors.New("invalid email or password")
)

type AuthService struct {
	users  *repository.UserRepository
	tokens *token.Manager
}

func NewAuthService(users *repository.UserRepository, tokens *token.Manager) *AuthService {
	return &AuthService{users: users, tokens: tokens}
}

func (s *AuthService) Register(ctx context.Context, email, password, name string) (*model.User, error) {
	_, err := s.users.FindByEmail(ctx, email)
	switch {
	case err == nil:
		return nil, ErrEmailTaken
	case !errors.Is(err, repository.ErrUserNotFound):
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &model.User{
		Email:        email,
		PasswordHash: string(hash),
		Name:         name,
		Role:         model.RoleStaff,
	}

	id, err := s.users.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	user.ID = id

	return user, nil
}

func (s *AuthService) Login(ctx context.Context, email, password string) (string, *model.User, error) {
	user, err := s.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	signed, err := s.tokens.Issue(user.ID, user.Role)
	if err != nil {
		return "", nil, err
	}

	return signed, user, nil
}
