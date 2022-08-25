package usecase

import (
	"context"

	"github.com/edgarSucre/chat/internal/domain"
)

type AdminUsecase struct {
	repo   AdminRepository
	hasher Secure
}

func NewAdminUsecase(repository AdminRepository, opt ...AdminUsecaseOption) *AdminUsecase {
	uc := &AdminUsecase{
		repo:   repository,
		hasher: Hasher{},
	}

	for _, v := range opt {
		v(uc)
	}

	return uc
}

func (uc *AdminUsecase) CreateUser(ctx context.Context, params domain.UserParam) (domain.UserResponse, *domain.Err) {
	hashed, err := uc.hasher.SecurePassword(params.Password)
	if err != nil {
		return domain.UserResponse{}, domain.WrapErrorf(
			err,
			domain.ErrorCodeInvalidParams,
			"can't hash password",
		)
	}

	params.Password = hashed

	return uc.repo.CreateUser(ctx, params)
}

func (uc *AdminUsecase) Login(ctx context.Context, params domain.UserParam) *domain.Err {
	user, err := uc.repo.GetUser(ctx, params.UserName)
	if err != nil {
		return err
	}

	if !uc.hasher.IsPasswordValid(params.Password, user.Password) {
		return domain.WrapErrorf(
			nil,
			domain.ErrCodeWrongPassword,
			"wrong password",
		)
	}

	return nil
}

func (uc *AdminUsecase) CreateRoom(ctx context.Context, params domain.RoomParam) (domain.RoomResponse, *domain.Err) {
	room, err := uc.repo.CreateRoom(ctx, params.Name)
	if err != nil {
		return domain.RoomResponse{}, err
	}

	return room, nil
}
