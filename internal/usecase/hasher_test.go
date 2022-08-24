package usecase_test

import (
	"testing"

	"github.com/edgarSucre/chat/internal/usecase"
	"github.com/stretchr/testify/require"
)

func TestHasher(t *testing.T) {
	hasher := usecase.Hasher{}
	pass := "secret"

	hashed, err := hasher.SecurePassword(pass)
	require.NoError(t, err)
	require.NotEmpty(t, hashed)

	require.True(t, hasher.IsPasswordValid(pass, hashed))
	require.False(t, hasher.IsPasswordValid("secrets", hashed))
}
