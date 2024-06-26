package jwtservice

import (
	"github.com/SakuraBurst/vigilant-octo-meme/internal/config"
	"github.com/SakuraBurst/vigilant-octo-meme/internal/services"
	"github.com/go-faster/errors"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"time"
)

type Service struct {
	secret   string
	tokenTTL time.Duration
	log      *slog.Logger
}

func New(cfg *config.Config, log *slog.Logger) *Service {
	return &Service{
		secret:   cfg.App.AppSecret,
		tokenTTL: cfg.App.TokenTTL,
		log:      log,
	}
}
func (s *Service) NewToken(isAdmin bool) (string, error) {
	log := s.log.With(slog.String("method", "NewToken"))
	log.Info("NewToken", slog.Bool("isAdmin", isAdmin))
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["isAdmin"] = isAdmin
	claims["exp"] = time.Now().Add(s.tokenTTL).Unix()

	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		log.Error(err.Error())
		return "", errors.Wrap(err, "sign token failed")
	}
	log.Info("NewToken", slog.String("token", tokenString))
	return tokenString, nil
}

// ParseToken возвращает флаг и ошибку. Флаг - админ ли пользователь, err - является ли токен валидным
func (s *Service) ParseToken(tokenString string) (bool, error) {
	log := s.log.With(slog.String("method", "ParseToken"))
	token, err := jwt.Parse(tokenString, func(_ *jwt.Token) (interface{}, error) {
		return []byte(s.secret), nil
	})
	if err != nil {
		log.Error(err.Error())
		return false, services.ErrTokenInvalid
	}
	claims := token.Claims.(jwt.MapClaims)
	t, err := claims.GetExpirationTime()
	if err != nil {
		log.Error(err.Error())
		return false, services.ErrUserNotAuthorized
	}
	if t.Time.Before(time.Now()) {
		log.Error("token is expired", slog.Time("exp", t.Time))
		return false, errors.Wrap(services.ErrUserNotAuthorized, "token is expired")
	}
	isAdmin, ok := claims["isAdmin"].(bool)
	if !ok {
		return false, nil
	}
	return isAdmin, nil
}
