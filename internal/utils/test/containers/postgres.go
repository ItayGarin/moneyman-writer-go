package containers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type PostgresTestContainer struct {
	Terminate func(context.Context) error
	Host      string
	Port      string
	Username  string
	Password  string
	Database  string
	URL       string
	DSN       string
}

func NewPostgresTestContainer() *PostgresTestContainer {
	ctx := context.Background()

	username := "user"
	password := "password"
	db := "test_db"
	req := testcontainers.ContainerRequest{
		Image:        "postgres:latest",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       db,
			"POSTGRES_USER":     username,
			"POSTGRES_PASSWORD": password,
		},
		WaitingFor: wait.ForAll(
			wait.ForListeningPort("5432/tcp"),
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(5*time.Second),
		),
		HostConfigModifier: func(config *container.HostConfig) {
			config.AutoRemove = true
		},
	}
	postgresContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("Failed to start container: %s", err)
	}

	mappedPort, err := postgresContainer.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("Failed to get mapped port: %s", err)
	}

	dbHost, err := postgresContainer.Host(ctx)
	if err != nil {
		log.Fatalf("Failed to get container IP: %s", err)
	}
	dbPort := mappedPort.Port()
	dsn := fmt.Sprintf("host=%s port=%s user=user password=password dbname=test_db sslmode=disable", dbHost, dbPort)
	url := fmt.Sprintf("postgres://user:password@%s:%s/%s?sslmode=%s", dbHost, dbPort, "test_db", "disable")

	return &PostgresTestContainer{
		Terminate: postgresContainer.Terminate,
		Host:      dbHost,
		Port:      dbPort,
		Username:  username,
		Password:  password,
		Database:  db,
		URL:       url,
		DSN:       dsn,
	}
}
