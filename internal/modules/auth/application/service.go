package application

import (
	"context"

	"bobshop/internal/modules/auth/domain"
	"bobshop/internal/platform/security"
)

type AuthService struct {
	userRepo  domain.AuthRepository
	tokenizer security.Tokenizer
}

func NewAuthService(userRepo domain.AuthRepository, tokenizer security.Tokenizer) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		tokenizer: tokenizer,
	}
}

func (s *AuthService) SignUp(ctx context.Context, email, password string) error {
	existingUser, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return domain.ErrUserAlreadyExists
	}

	newUser, err := domain.NewUser(email, password)
	if err != nil {
		return err
	}

	return s.userRepo.Create(ctx, newUser)
}

func (s *AuthService) SignIn(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if !user.CheckPassword(password) {
		return "", domain.ErrInvalidPassword
	}

	return s.tokenizer.GenerateToken(user.ID.String(), user.Role.String())
}
