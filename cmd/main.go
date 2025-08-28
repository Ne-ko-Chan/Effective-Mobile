package main

import (
	"database/sql"
	"log"
	"rest-service/cmd/api"
	"rest-service/config"
	"rest-service/db"
)

func main() {
	log.Println("DB: Setting up database connection")
	db, err := db.NewPSQLStorage(config.Envs.ConnString)
	if err != nil {
		log.Fatal(err)
	}

	if err = initStorage(db); err != nil {
		log.Fatalf("DB: Error connecting to database %v\n", err)
	}

	server := api.NewAPIServer(":"+config.Envs.Port, db)
	if err := server.Run(); err != nil {
		log.Fatalf("Critical error: %v\n", err)
	}
}

func initStorage(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}

	log.Println("DB: Successfully connected")
	return nil
}
