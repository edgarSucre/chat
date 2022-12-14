package usecase

import (
	"context"

	"github.com/edgarSucre/chat/internal/domain"
)

type (
	AdminRepository interface {
		CreateUser(context.Context, domain.UserParam) (domain.UserResponse, *domain.Err)
		GetUser(context.Context, string) (domain.UserResponse, *domain.Err)
		CreateRoom(context.Context, string) (domain.RoomResponse, *domain.Err)
	}

	Secure interface {
		IsPasswordValid(pass, hashed string) bool
		SecurePassword(pass string) (string, error)
	}
)
