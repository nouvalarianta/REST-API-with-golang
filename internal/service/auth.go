package service

import (
	"context"
	"errors"
	"golang-yt/domain"
	"golang-yt/dto"
	"golang-yt/internal/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	conf           *config.Config
	userRepository domain.UserRepository
}

func NewAuth(cnf *config.Config, userRepository domain.UserRepository) domain.AuthService {
	return authService{
		conf:           cnf,
		userRepository: userRepository,
	}
}

func (a authService) Login(ctx context.Context, req dto.AuthRequest) (dto.AuthResponse, error) {
	user, err := a.userRepository.GetByEmail(ctx, req.Email)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	if user.ID == "" {
		return dto.AuthResponse{}, errors.New("autentikasi gagal")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return dto.AuthResponse{}, errors.New("autentikasi gagal")
	}

	claims := jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Duration(a.conf.Jwt.Exp) * time.Minute).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(a.conf.Jwt.Key))
	if err != nil {
		return dto.AuthResponse{}, errors.New("autentikasi gagal")
	}

	return dto.AuthResponse{
		Token: tokenStr,
	}, nil
}
