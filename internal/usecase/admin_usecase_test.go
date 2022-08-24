package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/edgarSucre/chat/internal/domain"
	mockhash "github.com/edgarSucre/chat/internal/mock/hasher"
	mockrepo "github.com/edgarSucre/chat/internal/mock/repo"
	"github.com/edgarSucre/chat/internal/usecase"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestAdminUsecase(t *testing.T) {

	params := domain.UserParam{
		UserName: gofakeit.Username(),
		Password: gofakeit.Password(false, false, false, false, false, 10),
	}

	testCases := []struct {
		name      string
		buildStub func(repo *mockrepo.MockAdminRepository, hasher *mockhash.MockSecure)
		checkAns  func(t *testing.T, response domain.UserResponse, err error)
	}{
		{
			name: "OK",
			buildStub: func(repo *mockrepo.MockAdminRepository, hasher *mockhash.MockSecure) {
				hasher.EXPECT().
					SecurePassword(params.Password).
					Times(1).
					Return("hashed", nil)

				repoParams := domain.UserParam{
					UserName: params.UserName,
					Password: "hashed",
				}

				repo.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(repoParams)).
					Times(1).
					Return(domain.UserResponse{UserName: params.UserName}, nil)
			},
			checkAns: func(t *testing.T, response domain.UserResponse, err error) {
				require.NoError(t, err)
				require.NotEmpty(t, response)
				require.Equal(t, params.UserName, response.UserName)
			},
		},

		{
			name: "Failed to hash password",
			buildStub: func(repo *mockrepo.MockAdminRepository, hasher *mockhash.MockSecure) {
				hasher.EXPECT().
					SecurePassword(params.Password).
					Times(1).
					Return("", fmt.Errorf("can't parse"))

				repo.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkAns: func(t *testing.T, response domain.UserResponse, err error) {
				require.Error(t, err)
				require.ErrorIs(t, err, domain.ErrBadParamInput)
				require.Empty(t, response)
			},
		},

		{
			name: "Failed to create user",
			buildStub: func(repo *mockrepo.MockAdminRepository, hasher *mockhash.MockSecure) {
				hasher.EXPECT().
					SecurePassword(params.Password).
					Times(1).
					Return("hashed", nil)

				repoParams := domain.UserParam{
					UserName: params.UserName,
					Password: "hashed",
				}

				repo.EXPECT().
					CreateUser(gomock.Any(), gomock.Eq(repoParams)).
					Times(1).
					Return(domain.UserResponse{}, domain.ErrInternalServerError)
			},
			checkAns: func(t *testing.T, response domain.UserResponse, err error) {
				require.Error(t, err)
				require.ErrorIs(t, err, domain.ErrInternalServerError)
				require.Empty(t, response)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			repo := mockrepo.NewMockAdminRepository(ctrl)
			hasher := mockhash.NewMockSecure(ctrl)

			tc.buildStub(repo, hasher)

			uc := usecase.NewAdminUsecase(repo, usecase.WithHasher(hasher))
			uc.CreateUser(context.Background(), params)
		})
	}
}
