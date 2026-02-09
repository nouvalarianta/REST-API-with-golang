package domain

import (
	"context"
	"golang-yt/dto"
)

type AuthService interface {
	Login(ctx context.Context, req dto.AuthRequest) (dto.AuthResponse, error)
}
