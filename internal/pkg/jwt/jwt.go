package jwt

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/toji-dev/go-shop/internal/pkg/apperror"
	timeUtils "github.com/toji-dev/go-shop/internal/pkg/time"
)

type JWTConfig struct {
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
}

type CustomJwtClaims struct {
	GenerateTokenInput
	jwt.StandardClaims
}

type GenerateTokenInput struct {
	UserId string
	Email  string
	Role   string
}

type JwtService interface {
	GenerateAccessToken(input *GenerateTokenInput) (string, error)
	GenerateRefreshToken(input *GenerateTokenInput) (string, error)
	ValidateAccessToken(ctx context.Context, tokenString string) (*CustomJwtClaims, error)
	ValidateRefreshToken(tokenString string) (*CustomJwtClaims, error)
}

type jwtService struct {
	cfg JWTConfig
}

func NewJwtService(cfg JWTConfig) JwtService {
	return &jwtService{
		cfg: cfg,
	}
}

func (s *jwtService) GenerateAccessToken(input *GenerateTokenInput) (string, error) {
	now := timeUtils.GetUtcTime()
	claims := &CustomJwtClaims{
		*input,
		jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(s.cfg.AccessTokenTTL).Unix(),
			NotBefore: now.Unix(),
			Issuer:    s.cfg.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.SecretKey))
}

func (s *jwtService) GenerateRefreshToken(input *GenerateTokenInput) (string, error) {
	now := timeUtils.GetUtcTime()
	claims := &CustomJwtClaims{
		*input,
		jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(s.cfg.RefreshTokenTTL).Unix(),
			Issuer:    s.cfg.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.cfg.SecretKey))
}

func (s *jwtService) ValidateAccessToken(ctx context.Context, tokenString string) (*CustomJwtClaims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)

	if len(tokenString) == 0 {
		return nil, apperror.NewTokenInvalid("Token is missing")
	}

	claims := &CustomJwtClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperror.NewTokenInvalid("Invalid signing method")
		}
		return []byte(s.cfg.SecretKey), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return nil, apperror.NewTokenExpired("Token has expired")
			}
		}
		log.Printf("Error parsing token: %v", err)
		return nil, apperror.NewTokenInvalid("Token is invalid")
	}

	if !token.Valid {
		return nil, apperror.NewTokenInvalid("Token is not valid")
	}

	return claims, nil
}

func (s *jwtService) ValidateRefreshToken(tokenString string) (*CustomJwtClaims, error) {
	claims := &CustomJwtClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.SecretKey), nil
	})

	if err != nil {
		v, _ := err.(*jwt.ValidationError)

		if v.Errors == jwt.ValidationErrorExpired {
			return nil, apperror.NewTokenExpired("Refresh token has expired")
		}

		return nil, apperror.NewTokenInvalid("Refresh token is invalid")
	}

	if !token.Valid {
		return nil, apperror.NewTokenInvalid("Refresh token is invalid")
	}

	return claims, nil
}
