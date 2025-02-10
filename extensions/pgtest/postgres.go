package pgtest

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // postgres driver
	_ "github.com/golang-migrate/migrate/v4/source/file"       // file driver
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
)

const _defaultContainerExpirationTime uint = 60 // seconds

type Config struct {
	// Expire container that takes more than `Expire` seconds running.
	// It is very important to use when you are debugging and killed the
	// process before the call the teardown or if a panic happens.
	Expire        uint
	MigrationPath string
}

func StartDockerContainer(cfg Config) (teardownFn func(), err error) {
	dockerPool, err := dockertest.NewPool("")
	if err != nil {
		return nil, fmt.Errorf("creating docker pool: %w", err)
	}

	if err = dockerPool.Client.Ping(); err != nil {
		return nil, fmt.Errorf("pinging to docker: %w", err)
	}

	resource, err := dockerPool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env:        []string{"POSTGRES_USER=postgres", "POSTGRES_PASSWORD=postgres"},
	}, func(hc *docker.HostConfig) {
		hc.AutoRemove = true
		hc.RestartPolicy = docker.NeverRestart()
	})
	if err != nil {
		return nil, fmt.Errorf("starting docker container: %w", err)
	}

	containerExpirationTime := _defaultContainerExpirationTime
	if cfg.Expire != 0 {
		containerExpirationTime = cfg.Expire
	}

	if err = resource.Expire(containerExpirationTime); err != nil {
		return nil, fmt.Errorf("setting container expire timeout: %w", err)
	}

	port := resource.GetPort("5432/tcp")

	var pgxPool *pgxpool.Pool

	if err = dockerPool.Retry(func() error {
		pgxPool, err = pgxpool.New(context.Background(), getConnString(port, "postgres"))
		if err != nil {
			return fmt.Errorf("creating pgx pool: %w", err)
		}

		if err = pgxPool.Ping(context.Background()); err != nil {
			return fmt.Errorf("pinging to postgres: %w", err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("connecting to postgres: %w", err)
	}

	// The migrations will be performed on the template1 database, so when a new database is created it will
	// already have the migrations applied.
	if err = runMigrations(getConnString(port, "template1"), cfg.MigrationPath); err != nil {
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	concurrentPool = pgxPool

	teardownFn = func() {
		pgxPool.Close()

		_ = dockerPool.Purge(resource)
	}

	return teardownFn, nil
}

func getConnString(port, dbName string) string {
	return fmt.Sprintf("postgres://postgres:postgres@localhost:%s/%s?sslmode=disable", port, dbName)
}

func getGoModuleRoot() (string, error) {
	cmd := exec.Command("go", "env", "GOMOD")
	cmd.Env = os.Environ()

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("getting go env GOMOD output: %w", err)
	}

	return filepath.Dir(string(output)), nil
}

func runMigrations(connString, mpath string) error {
	migrationPath := "gateways/postgres/migrations"
	if mpath != "" {
		migrationPath = mpath
	}
	rootPath, err := getGoModuleRoot()
	if err != nil {
		return fmt.Errorf("getting go module root: %w", err)
	}

	path := filepath.Join(rootPath, migrationPath)
	m, err := migrate.New("file://"+path, connString)
	if err != nil {
		return fmt.Errorf("creating migrate instance: %w", err)
	}

	if err = m.Up(); err != nil {
		return fmt.Errorf("running up migrations: %w", err)
	}

	serr, err := m.Close()
	if serr != nil {
		return fmt.Errorf("closing the source: %w", serr)
	}

	if err != nil {
		return fmt.Errorf("closing pg connection: %w", err)
	}

	return nil
}
