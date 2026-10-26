package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadUsesDefaultDatabasePath(t *testing.T) {
	settings, err := Load(map[string]string{})

	require.NoError(t, err)
	assert.Equal(t, "data/tiebreaker.db", settings.DatabasePath)
}

func TestLoadReadsDatabasePathFromEnvironment(t *testing.T) {
	settings, err := Load(map[string]string{
		"TIEBREAKER_DATABASE_PATH": "/var/lib/tiebreaker/tiebreaker.db",
	})

	require.NoError(t, err)
	assert.Equal(t, "/var/lib/tiebreaker/tiebreaker.db", settings.DatabasePath)
}
