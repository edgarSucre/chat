package config_test

import (
	"os"
	"testing"

	"github.com/edgarSucre/chat/config"
	"github.com/stretchr/testify/require"
)

func TestGetEnv(t *testing.T) {
	content := `SOME_KEY=secret`
	err := os.WriteFile("temp.env", []byte(content), 0666)
	require.NoError(t, err)

	defer func() {
		err := os.Remove("temp.env")
		require.NoError(t, err)
	}()

	t.Run("Reading from file", func(t *testing.T) {
		temp := &(struct{ SOME_KEY string }{})
		err := config.GetEnv("temp", temp)
		require.NoError(t, err)
		require.NotEmpty(t, temp)
		require.Equal(t, temp.SOME_KEY, "secret")
	})

	t.Run("Overriding values", func(t *testing.T) {
		os.Setenv("SOME_KEY", "another key")
		temp := &(struct{ SOME_KEY string }{})
		err := config.GetEnv("temp", temp)
		require.NoError(t, err)
		require.NotEmpty(t, temp)
		require.Equal(t, temp.SOME_KEY, "another key")
	})

	t.Run("Non existing file", func(t *testing.T) {
		temp := &(struct{ SOME_KEY string }{})
		err := config.GetEnv("asdasdasd", temp)
		require.Error(t, err)
		require.Empty(t, temp)
	})

	t.Run("Parse error", func(t *testing.T) {
		_ = &(struct{ SOME_KEY string }{})
		temp := 0
		err := config.GetEnv("temp", &temp)
		require.Error(t, err)
		require.Empty(t, temp)
	})
}
