package jwt

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	timeUtils "github.com/toji-dev/go-shop/internal/pkg/time"
	"github.com/toji-dev/go-shop/internal/services/user-service/internal/config"
	errorConstants "github.com/toji-dev/go-shop/internal/services/user-service/internal/pkg/errors"
)

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
	cfg *config.Config
}

func NewJwtService(cfg *config.Config) JwtService {
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
			ExpiresAt: now.Add(s.cfg.JWT.AccessTokenTTL).Unix(),
			NotBefore: now.Unix(),
			Issuer:    s.cfg.JWT.Issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWT.SecretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *jwtService) GenerateRefreshToken(input *GenerateTokenInput) (string, error) {
	claims := &CustomJwtClaims{
		*input,
		jwt.StandardClaims{
			ExpiresAt: time.Now().UTC().Add(s.cfg.JWT.RefreshTokenTTL).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.cfg.JWT.SecretKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *jwtService) ValidateAccessToken(ctx context.Context, tokenString string) (*CustomJwtClaims, error) {
	// Remove "Bearer " prefix if present
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)

	if len(tokenString) == 0 {
		log.Println("Token string is empty after processing")
		return nil, errorConstants.ErrTokenInvalid
	}

	claims := &CustomJwtClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra thuật toán
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Println("Invalid signing method")
			return nil, errorConstants.ErrTokenInvalid
		}
		return []byte(s.cfg.JWT.SecretKey), nil
	})

	if err != nil {
		log.Printf("Error parsing token: %v (Type: %T)", err, err)

		// Kiểm tra lỗi cụ thể
		if ve, ok := err.(*jwt.ValidationError); ok {
			log.Printf("ValidationError details: Errors=%d, Inner=%v", ve.Errors, ve.Inner)

			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				log.Println("Token expired")
				return nil, errorConstants.ErrTokenExpired
			}
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				log.Println("Token malformed")
				return nil, errorConstants.ErrTokenInvalid
			}
			if ve.Errors&jwt.ValidationErrorSignatureInvalid != 0 {
				log.Println("Token signature invalid")
				return nil, errorConstants.ErrTokenInvalid
			}
			if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				log.Println("Token not valid yet")
				return nil, errorConstants.ErrTokenInvalid
			}
			if ve.Errors&jwt.ValidationErrorUnverifiable != 0 {
				log.Println("Token unverifiable")
				return nil, errorConstants.ErrTokenInvalid
			}
		}

		return nil, errorConstants.ErrTokenInvalid
	}

	if !token.Valid {
		log.Println("Token is not valid")
		return nil, errorConstants.ErrTokenInvalid
	}

	return claims, nil
}

func (s *jwtService) ValidateRefreshToken(tokenString string) (*CustomJwtClaims, error) {
	claims := &CustomJwtClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.JWT.SecretKey), nil
	})

	if err != nil {
		v, _ := err.(*jwt.ValidationError)

		if v.Errors == jwt.ValidationErrorExpired {
			return nil, errorConstants.ErrTokenExpired
		}

		return nil, errorConstants.ErrTokenInvalid
	}

	if !token.Valid {
		return nil, errorConstants.ErrTokenInvalid
	}

	return claims, nil
}
