package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"github.com/Dmitro44/relper/internal/config"
	"github.com/Dmitro44/relper/internal/entities/user"
	apperr "github.com/Dmitro44/relper/internal/errors"
	"github.com/Dmitro44/relper/internal/repositories/refreshtokens"
	"github.com/Dmitro44/relper/internal/services/users"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	jwt.RegisteredClaims
	Roles []string `json:"roles,omitempty"`
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type AuthService struct {
	userservice *users.UserService
	reftokrepo  refreshtokens.Repository
	config      *config.Config
}

func NewAuthService(userservice *users.UserService, reftokrepo refreshtokens.Repository, config *config.Config) *AuthService {
	return &AuthService{
		userservice: userservice,
		reftokrepo:  reftokrepo,
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

func (s *AuthService) Register(ctx context.Context, input CreateInput) (*user.User, *TokenPair, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, nil, err
	}

	_, err = s.userservice.GetByEmail(ctx, input.Email)
	if err == nil {
		return nil, nil, apperr.NewConflictError("user with this email already exists")
	}

	if appErr, ok := errors.AsType[*apperr.AppError](err); !ok || appErr.Code != apperr.ErrorCodeNotFound {
		return nil, nil, err
	}

	user, err := s.userservice.Create(ctx, users.CreateInput{
		Email:        input.Email,
		PasswordHash: hash,
		Nickname:     input.Nickname,
		Habits:       input.Habits,
		Contacts:     input.Contacts,
	})
	if err != nil {
		return nil, nil, err
	}

	pair, err := s.generateTokenPair(ctx, user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, pair, nil
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (*user.User, *TokenPair, error) {
	user, err := s.userservice.GetByEmail(ctx, email)
	if err != nil {
		return nil, nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, nil, apperr.NewUnauthorizedError("invalid credentials")
	}

	pair, err := s.generateTokenPair(ctx, user.ID)
	if err != nil {
		return nil, nil, err
	}

	return user, pair, nil
}

func (s *AuthService) Logout(ctx context.Context, token string) error {
	sum := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(sum[:])

	err := s.reftokrepo.RevokeByHash(ctx, tokenHash)
	if err != nil {
		return apperr.NewInternalError("failed to revoke refresh token", err)
	}

	return nil
}

func (s *AuthService) Refresh(ctx context.Context, token string) (*TokenPair, error) {
	sum := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(sum[:])

	storedToken, err := s.reftokrepo.GetByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, refreshtokens.ErrNotFound) {
			return nil, apperr.NewUnauthorizedError("invalid refresh token")
		}
		return nil, apperr.NewInternalError("failed to look up refresh token", err)
	}

	if storedToken.RevokedAt != nil {
		if err := s.reftokrepo.RevokeAllByUser(ctx, storedToken.UserID); err != nil {
			return nil, apperr.NewInternalError("failed to revoke user tokens", err)
		}
		return nil, apperr.NewUnauthorizedError("refresh token reuse detected")
	}

	if time.Now().After(storedToken.ExpiresAt) {
		return nil, apperr.NewUnauthorizedError("refresh token expired")
	}

	pair, err := s.generateTokenPair(ctx, storedToken.UserID)
	if err != nil {
		return nil, err
	}

	if err := s.reftokrepo.RevokeByHash(ctx, tokenHash); err != nil {
		return nil, apperr.NewInternalError("failed to revoke refresh token", err)
	}

	return pair, nil
}

func (s *AuthService) generateTokenPair(ctx context.Context, userID uuid.UUID) (*TokenPair, error) {
	accessToken, err := s.generateAccessToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(ctx, userID)
	if err != nil {
		return nil, err
	}

	pair := &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return pair, nil
}

func (s *AuthService) generateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
	bytes := make([]byte, 32)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", apperr.NewInternalError("failed to generate refresh token", err)
	}

	token := base64.RawURLEncoding.EncodeToString(bytes)

	hash := sha256.Sum256([]byte(token))
	tokenHash := hex.EncodeToString(hash[:])

	err = s.reftokrepo.Create(ctx, refreshtokens.CreateInput{
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(15 * 24 * time.Hour),
	})
	if err != nil {
		return "", apperr.NewInternalError("failed to save refresh token", err)
	}

	return token, nil
}

func (s *AuthService) ValidateAccessToken(ctx context.Context, token string) (*Claims, error) {
	var claims Claims
	keyFunc := func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperr.NewUnauthorizedError("jwt token is invalid")
		}
		return []byte(s.config.JWTSecret), nil
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

func (s *AuthService) generateAccessToken(ctx context.Context, userID uuid.UUID) (string, error) {
	roles, err := s.userservice.GetRoles(ctx, userID)
	if err != nil {
		return "", err
	}

	now := time.Now()
	claims := Claims{
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
		Roles:     roles,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	access, err := token.SignedString([]byte(s.config.JWTSecret))
	if err != nil {
		return "", err
	}

	return access, nil
}
