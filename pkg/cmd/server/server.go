package cmd

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/AYaro/go-test-task/pkg/protocol/grpc"
	v1 "github.com/AYaro/go-test-task/pkg/service/v1"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func RunServer() error {
	ctx := context.Background()

	db, err := sql.Open("postgres", "host=postgres port=5433 user=admin password=password dbname=parts sslmode=disable")
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()
	m, err := migrate.New(
		"file://db/migrations",
		"postgres://admin:password@postgres:5433/parts?sslmode=disable")
	if err != nil {
		return fmt.Errorf("failed to get migrations: %v", err)
	}
	if err := m.Up(); err != nil {
		fmt.Printf("failed to up migration: %v", err)
	}

	API := v1.NewPartServiceServer(db)

	return grpc.RunServer(ctx, API, "9090")
}
