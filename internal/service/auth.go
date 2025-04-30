package service

import (
	"context"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
	"lost-items-service/internal/model"
	"lost-items-service/internal/service/pkg/hash"
	"lost-items-service/pkg/jwtutils"
)

type AuthService struct {
	userRepository UserRepository
	jwtSecret      string
}

func NewAuthService(
	repo UserRepository, jwt string,
) *AuthService {
	return &AuthService{
		userRepository: repo,
		jwtSecret:      jwt,
	}
}

func (s *AuthService) Registration(ctx context.Context, user model.User) (*model.User, error) {
	hashPass, err := hash.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash pass")
	}

	user.Password = hashPass

	if _, err = s.userRepository.UserByEmail(ctx, user.Email); err == nil {
		return nil, fmt.Errorf("user already exist")
	}

	userID, err := s.userRepository.CreateUser(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user")
	}
	user.ID = userID

	return &user, nil
}

func (s *AuthService) Authenticate(ctx context.Context, user model.User) (string, error) {
	current, err := s.userRepository.UserByEmail(ctx, user.Email)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(current.Password), []byte(user.Password)); err != nil {
		return "", fmt.Errorf("invalid email or password")
	}

	token, err := s.generateJWT(current.ID.String(), current.Role)

	return token, err
}

func (s *AuthService) generateJWT(userID string, role string) (string, error) {
	claims := map[string]interface{}{
		"userId": userID,
		"role":   role,
	}

	token, err := jwtutils.Generate(claims, 24*time.Hour, s.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to generate JWT token")
	}

	return token, nil
}
