package migrations

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func run(url, path string, f func(m *migrate.Migrate) error) {
	m, err := migrate.New(fmt.Sprintf("file://%v", path), url)
	if err != nil {
		log.Fatal(err)
	}
	if err := f(m); err != nil {
		log.Fatal(err)
	}
}

func Up(url,path string) {
	run(url, path, func(m *migrate.Migrate) error { return m.Up() })
}

func Down(url,path string) {
	run(url, path, func(m *migrate.Migrate) error { return m.Down() })
}
