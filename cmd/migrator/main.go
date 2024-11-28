package main

import (
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"log"
)

func main() {
	var migrationsPath string

	flag.StringVar(&migrationsPath, "migrations-path", "", "Path to a directory containing the migration files")
	flag.Parse()

	if migrationsPath == "" {
		panic("migrations-path is required")
	}

	dbURL := "postgres://postgres:root@localhost:5432/grpc-auth?sslmode=disable"

	// Open a connection to the database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic("unable to connect to postgres")
	}

	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			panic("unable to close postgres connection")
		}
	}(db)

	// Ping to ensure the connection is established
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create a Postgres driver instance for migrate
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic("unable to configure database")
	}

	// Initialize migrate with the migration source and database driver
	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationsPath),
		"grpc-auth",
		driver,
	)
	if err != nil {
		panic("Unable to Initialize migrate")
	}

	// Run the migrations up
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("No migrations to apply")

			return
		}

		panic("Unable to Run migration")
	}

	fmt.Println("Migrations ran successfully!")
}
