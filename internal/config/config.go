// Package config reads the app's settings from environment variables.
package config

import "github.com/caarlos0/env/v11"

// Settings hold everything the app reads from its environment.
type Settings struct {
	// DatabasePath is the location of the SQLite database file.
	DatabasePath string `env:"DATABASE_PATH" envDefault:"data/tiebreaker.db"`
}

// Load reads Settings from environ, a map of variable names to values.
func Load(environ map[string]string) (Settings, error) {
	return env.ParseAsWithOptions[Settings](env.Options{
		Environment: environ,
		Prefix:      "TIEBREAKER_",
	})
}
