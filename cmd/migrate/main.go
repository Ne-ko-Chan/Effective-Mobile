package main

import (
	"log"
	"os"
	"rest-service/config"
	"rest-service/db"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func main() {
	log.Print("Performing migrations...")
	db, err := db.NewPSQLStorage(config.Envs.ConnString)
	if err != nil {
		log.Fatal(err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	cmd := os.Args[len(os.Args)-1]
	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}

	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}

	if cmd == "newest" {
		for {
			if err := m.Up(); err != nil {
				if err == migrate.ErrNoChange {
					break
				}
				log.Fatal(err)
			}
		}
	}

	if cmd == "oldest" {
		for {
			if err := m.Down(); err != nil {
				if err == migrate.ErrNoChange {
					break
				}
				log.Fatal(err)
			}
		}
	}
	log.Println("Migrations applied successfully")
}
