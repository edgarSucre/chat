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
		input     domain.UserParam
		buildStub func(repo *mockrepo.MockAdminRepository, hasher *mockhash.MockSecure)
		checkAns  func(t *testing.T, response domain.UserResponse, err *domain.Err)
	}{
		{
			name:  "OK",
			input: params,
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
			checkAns: func(t *testing.T, response domain.UserResponse, err *domain.Err) {
				require.Nil(t, err)
				require.NotEmpty(t, response)
				require.Equal(t, params.UserName, response.UserName)
			},
		},

		{
			name:  "Failed to hash password",
			input: params,
			buildStub: func(repo *mockrepo.MockAdminRepository, hasher *mockhash.MockSecure) {
				hasher.EXPECT().
					SecurePassword(params.Password).
					Times(1).
					Return("", domain.WrapErrorf(fmt.Errorf(""), domain.ErrorCodeInvalidParams, ""))

				repo.EXPECT().
					CreateUser(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkAns: func(t *testing.T, response domain.UserResponse, err *domain.Err) {
				require.Error(t, err)
				require.Equal(t, err.Code(), domain.ErrorCodeInvalidParams)
				require.Empty(t, response)
			},
		},

		{
			name:  "Failed to create user",
			input: params,
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
					Return(domain.UserResponse{}, domain.WrapErrorf(fmt.Errorf(""), domain.ErrorCodeInternalRepository, ""))
			},
			checkAns: func(t *testing.T, response domain.UserResponse, err *domain.Err) {
				require.Error(t, err)
				require.Equal(t, err.Code(), domain.ErrorCodeInternalRepository)
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
			rsp, err := uc.CreateUser(context.Background(), tc.input)
			tc.checkAns(t, rsp, err)
		})
	}
}

func TestLogin(t *testing.T) {
	params := domain.UserParam{
		UserName: gofakeit.Username(),
		Password: gofakeit.Password(false, false, false, false, false, 10),
	}

	testCases := []struct {
		name      string
		input     domain.UserParam
		buildStub func(repo *mockrepo.MockAdminRepository, hasher *mockhash.MockSecure)
		checkAns  func(t *testing.T, err *domain.Err)
	}{
		{
			name:  "OK",
			input: params,
			buildStub: func(repo *mockrepo.MockAdminRepository, hasher *mockhash.MockSecure) {
				repo.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(params.UserName)).
					Times(1).
					Return(
						domain.UserResponse{UserName: params.UserName, Password: "hashed"},
						nil,
					)

				hasher.EXPECT().
					IsPasswordValid(gomock.Eq(params.Password), gomock.Eq("hashed")).
					Times(1).
					Return(true)
			},
			checkAns: func(t *testing.T, err *domain.Err) {
				require.Nil(t, err)
			},
		},

		{
			name:  "User not found",
			input: params,
			buildStub: func(repo *mockrepo.MockAdminRepository, hasher *mockhash.MockSecure) {
				repo.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(params.UserName)).
					Times(1).
					Return(
						domain.UserResponse{},
						domain.WrapErrorf(fmt.Errorf("nothing"), domain.ErrorCodeUserNotFound, ""),
					)

				hasher.EXPECT().
					IsPasswordValid(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkAns: func(t *testing.T, err *domain.Err) {
				require.Error(t, err)
				require.Equal(t, err.Code(), domain.ErrorCodeUserNotFound)
			},
		},

		{
			name:  "Wrong Password",
			input: params,
			buildStub: func(repo *mockrepo.MockAdminRepository, hasher *mockhash.MockSecure) {
				repo.EXPECT().
					GetUser(gomock.Any(), gomock.Eq(params.UserName)).
					Times(1).
					Return(
						domain.UserResponse{UserName: params.UserName, Password: "hashed"},
						nil,
					)

				hasher.EXPECT().
					IsPasswordValid(gomock.Eq(params.Password), gomock.Eq("hashed")).
					Times(1).
					Return(false)
			},
			checkAns: func(t *testing.T, err *domain.Err) {
				require.Error(t, err)
				require.Equal(t, err.Code(), domain.ErrCodeWrongPassword)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			repo := mockrepo.NewMockAdminRepository(ctrl)
			hasher := mockhash.NewMockSecure(ctrl)

			tc.buildStub(repo, hasher)

			uc := usecase.NewAdminUsecase(repo, usecase.WithHasher(hasher))
			tc.checkAns(t, uc.Login(context.Background(), tc.input))
		})
	}
}
