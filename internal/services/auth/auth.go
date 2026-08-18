package auth

import (
	"context"
	"time"

	"github.com/Dmitro44/relper/internal/config"
	"github.com/Dmitro44/relper/internal/entities/user"
	"github.com/Dmitro44/relper/internal/services/users"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	userservice *users.UserService
	config      *config.Config
}

func NewAuthService(userservice *users.UserService, config *config.Config) *AuthService {
	return &AuthService{
		userservice: userservice,
		config:      config,
	}
}

type CreateInput struct {
	Email    string
	Password string
	Nickname string
	Habits   []string
	Contacts user.Contacts
}

func (s *AuthService) Register(ctx context.Context, input CreateInput) (*user.User, string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	_, err = s.userservice.GetByEmail(ctx, input.Email)
	if err != nil {
		return nil, "", err
	}

	user, err := s.userservice.Create(ctx, users.CreateInput{
		Email:        input.Email,
		PasswordHash: hash,
		Nickname:     input.Nickname,
		Habits:       input.Habits,
		Contacts:     input.Contacts,
	})
	if err != nil {
		return nil, "", err
	}

	accessToken, err := s.issueToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, accessToken, nil
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (*user.User, string, error) {
	user, err := s.userservice.GetByEmail(ctx, email)
	if err != nil {
		return nil, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", err
	}

	accessToken, err := s.issueToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, accessToken, nil
}

func (s *AuthService) ValidateAccessToken(ctx context.Context, token string) (*jwt.RegisteredClaims, error) {
	var claims jwt.RegisteredClaims
	keyFunc := func(t *jwt.Token) (any, error) {
		// if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
		// 	return apperr...
		// }
		return s.config.JWTSecret, nil
	}
	_, err := jwt.ParseWithClaims(
		token, &claims, keyFunc,
		jwt.WithValidMethods([]string{"HS256"}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, err
	}

	return &claims, nil
}

func (s *AuthService) issueToken(userID uuid.UUID) (string, error) {
	now := time.Now()
	claims := jwt.RegisteredClaims{
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	access, err := token.SignedString(s.config.JWTSecret)
	if err != nil {
		return "", err
	}

	return access, nil
}
