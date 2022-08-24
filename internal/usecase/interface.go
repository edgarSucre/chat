package usecase

import (
	"context"

	"github.com/edgarSucre/chat/internal/domain"
)

type (
	AdminRepository interface {
		CreateUser(context.Context, domain.UserParam) (domain.UserResponse, error)
	}

	Secure interface {
		IsPasswordValid(pass, hashed string) bool
		SecurePassword(pass string) (string, error)
	}
)
