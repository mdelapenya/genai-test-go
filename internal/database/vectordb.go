package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	dbName     = "documents"
	dbUser     = "user"
	dbPassword = "password"
)

var dbPort uint16 = 5432

func MustVectorDatabase(ctx context.Context) string {
	postgresContainer, err := postgres.Run(ctx,
		"docker.io/pgvector/pgvector:pg16",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
		testcontainers.CustomizeRequest(testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Name: "pgvector",
			},
			Reuse: true,
		}),
	)
	if err != nil {
		panic(err)
	}

	p, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		panic(err)
	}

	dbPort = uint16(p.Int())

	// Example: "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"
	return postgresContainer.MustConnectionString(ctx, "sslmode=disable")
}

// returns callback to close the connection and an error if any
func VectorDBClient(connString string) (*pgx.ConnPool, error) {
	pgxConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     "localhost",
			Database: dbName,
			User:     dbUser,
			Password: dbPassword,
			Port:     dbPort,
		},
	}

	conn, err := pgx.NewConnPool(pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("connect to pgx: %w", err)
	}

	return conn, nil
}

func CheckInitialEmbeddings(connString string) bool {
	c, err := VectorDBClient(connString)
	if err != nil {
		return false
	}
	defer c.Close()

	r, err := c.Query("SELECT COUNT(1) FROM documentsembeddings")
	if err != nil {
		return false
	}

	if r.Err() != nil {
		return false
	}

	// check if count is greater than 0
	if r.Next() {
		var count int
		err = r.Scan(&count)
		if err != nil {
			return false
		}

		if count > 0 {
			return false
		}
	}

	return true
}
