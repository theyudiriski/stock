package postgres

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"net/url"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"

	"stock/config"
)

type Client struct {
	Leader *sql.DB
}

func NewClient(config config.Database, migrate bool) (*Client, error) {
	leaderDB, err := openDB(config.Leader)
	if err != nil {
		return nil, fmt.Errorf("failed to open leader DB: %w", err)
	}

	if migrate {
		if err = migrateDB(leaderDB); err != nil {
			return nil, fmt.Errorf("migrateDB(): %w", err)
		}
	}

	return &Client{
		Leader: leaderDB,
	}, nil
}

func openDB(config config.DatabaseConfig) (*sql.DB, error) {
	var (
		db  *sql.DB
		err error
	)

	dbURI := &url.URL{
		Scheme: config.Scheme,
		User:   url.UserPassword(config.Username, config.Password),
		Host:   fmt.Sprintf("%s:%s", config.Host, config.Port),
		Path:   config.DB,
	}

	db, err = sql.Open("pgx", dbURI.String())
	if err != nil {
		return nil, fmt.Errorf("sql.Open(): %w", err)
	}

	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	return db, nil
}

//go:embed migrations/*.sql
var emfs embed.FS

func migrateDB(db *sql.DB) error {
	d, err := iofs.New(emfs, "migrations")
	if err != nil {
		return err
	}

	driver, err := pgx.WithInstance(db, &pgx.Config{})
	if err != nil {
		return fmt.Errorf("pgx.WithInstance(): %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "pgx", driver)
	if err != nil {
		return fmt.Errorf("migrate.New(): %w", err)
	}

	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}
