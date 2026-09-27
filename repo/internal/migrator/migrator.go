package migrator

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"log"
	"repo/internal/config"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

type Migrator struct {
	migrationsPath string
}

func NewMigrator() *Migrator {
	return &Migrator{
		migrationsPath: "migrations",
	}
}

func (m *Migrator) Migrate(cfg config.Config) error {

	db, err := sql.Open("postgres", cfg.DSN())
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		return fmt.Errorf("unable to ping database: %v", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("unable to create db instance: %v", err)
	}

	sourceDriver, err := iofs.New(migrationsFS, m.migrationsPath)
	if err != nil {
		return fmt.Errorf("unable to create source driver: %v", err)
	}

	migrator, err := migrate.NewWithInstance("iofs", sourceDriver, "postgres", driver)
	if err != nil {
		return fmt.Errorf("unable to create migrator: %v", err)
	}
	defer migrator.Close()

	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("unable to apply migrations: %v", err)
	}

	version, dirty, err := migrator.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("unable to get migration version: %v", err)
	}

	log.Printf("Migrations applied successfully. Version: %d, Dirty: %t", version, dirty)
	return nil
}

func GetDSN(cfg config.Config) string {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)
	return dsn
}
