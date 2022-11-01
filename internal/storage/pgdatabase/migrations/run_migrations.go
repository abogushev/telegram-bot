package migrations

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	. "gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/logger"
	"go.uber.org/zap"
)

func run(url, path string, f func(m *migrate.Migrate) error) {
	m, err := migrate.New(fmt.Sprintf("file://%v", path), url)
	if err != nil {
		Log.Fatal("failed to read migrations", zap.Error(err))
	}
	if err := f(m); err != nil && err != migrate.ErrNoChange {
		Log.Fatal("failed to execute migrations", zap.Error(err))
	}
}

func Up(url, path string) {
	run(url, path, func(m *migrate.Migrate) error { return m.Up() })
}

func Down(url, path string) {
	run(url, path, func(m *migrate.Migrate) error { return m.Down() })
}
