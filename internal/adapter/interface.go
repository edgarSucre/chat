package adapter

import (
	"context"

	"github.com/edgarSucre/chat/internal/domain"
)

type AdminUseCase interface {
	CreateUser(context.Context, domain.UserParam) (domain.UserResponse, error)
	Login(context.Context, domain.UserParam) error
	CreateRoom(context.Context, domain.RoomParam) (domain.RoomResponse, error)
}
