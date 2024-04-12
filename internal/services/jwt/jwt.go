package jwtservice

import (
	"github.com/SakuraBurst/vigilant-octo-meme/internal/config"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/services"
	"github.com/go-faster/errors"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

type Service struct {
	secret   string
	tokenTTL time.Duration
}

func New(cfg *config.Config) *Service {
	return &Service{
		secret:   cfg.App.AppSecret,
		tokenTTL: cfg.App.TokenTTL,
	}
}
func (s *Service) NewToken(isAdmin bool) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["isAdmin"] = isAdmin
	claims["exp"] = time.Now().Add(s.tokenTTL).Unix()

	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", errors.Wrap(err, "sign token failed")
	}
	return tokenString, nil
}

// err == nil - значит токен валидный, булевое значение - админ ли пользователь
func (s *Service) ParseToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(_ *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		return false, errors.Wrap(err, "parse token failed")
	}
	claims := token.Claims.(jwt.MapClaims)
	t, err := claims.GetExpirationTime()
	if err != nil {
		return false, services.ErrUserNotAuthorized
	}
	if t.Time.Before(time.Now()) {
		return false, errors.Wrap(services.ErrUserNotAuthorized, "token is expired")
	}
	isAdmin, ok := claims["isAdmin"].(bool)
	if !ok {
		return false, nil
	}
	return isAdmin, nil
}
