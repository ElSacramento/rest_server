package middleware

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfig_ValidateConfig(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		cfg := Config{
			Database: DatabaseConfig{
				Host:     "t.bk.ru",
				Port:     "1234",
				User:     "test",
				Password: "test",
				Name:     "db"},
			Server: ServerConfig{
				Host: "e.mail.ru",
				Port: "1234",
			}}
		err := cfg.ValidateConfig()
		require.NoError(t, err)
	})
	t.Run("FAIL", func(t *testing.T) {
		t.Run("empty db password", func(t *testing.T) {
			cfg := Config{
				Database: DatabaseConfig{
					Host: "t.bk.ru",
					Port: "1234",
					User: "test",
					Name: "db"},
				Server: ServerConfig{
					Host: "e.mail.ru",
					Port: "1234",
				}}
			err := cfg.ValidateConfig()
			require.Error(t, err)
		})
		t.Run("empty server port", func(t *testing.T) {
			cfg := Config{
				Database: DatabaseConfig{
					Host:     "t.bk.ru",
					Port:     "1234",
					User:     "test",
					Password: "test",
					Name:     "db"},
				Server: ServerConfig{
					Host: "e.mail.ru",
				}}
			err := cfg.ValidateConfig()
			require.Error(t, err)
		})
		t.Run("empty server cfg", func(t *testing.T) {
			cfg := Config{
				Database: DatabaseConfig{
					Host:     "t.bk.ru",
					Port:     "1234",
					User:     "test",
					Password: "test",
					Name:     "db",
				}}
			err := cfg.ValidateConfig()
			require.Error(t, err)
		})
		t.Run("empty db cfg", func(t *testing.T) {
			cfg := Config{
				Server: ServerConfig{
					Port: "1234",
				}}
			err := cfg.ValidateConfig()
			require.Error(t, err)
		})
	})
}
