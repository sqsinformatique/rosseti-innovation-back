package db

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/cfg"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/context"

	// PostgreSQL driver.
	_ "github.com/lib/pq"
)

type empty struct{}

type DB struct {
	log        zerolog.Logger
	conn       *sqlx.DB
	config     *cfg.AppCfg
	schemaName string
}

// Initialize initialize database
func NewDB(ctx *context.Context) (*DB, error) {
	if ctx == nil || ctx.Config == nil {
		return nil, errors.New("empty context or config")
	}
	db := &DB{}
	db.config = ctx.Config
	db.log = ctx.GetPackageLogger(empty{})
	ctx.RegisterDatabase(&db.conn)

	db.SetSchemaName("production")

	return db, nil
}

// Start start db
func (db *DB) Start() error {
	// Compose DSN.
	dsn := db.config.Database.DSN
	if db.config.Database.Params != "" {
		dsn += "?" + db.config.Database.Params
	}

	dbConn, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		db.log.Fatal().Err(err).Msgf("Failed to connect to PostgreSQL database")
	}

	db.log.Info().Msg("Database connection established")

	// Set connection pooling options.
	maxConnectionLifetime, err := time.ParseDuration(db.config.Database.MaxConnectionLifetime)
	if err != nil {
		db.log.Fatal().Err(err).Msg("Failed to parse MaxConnectionLifetime")
	}

	dbConn.SetConnMaxLifetime(maxConnectionLifetime)
	dbConn.SetMaxIdleConns(db.config.Database.MaxIdleConnections)
	dbConn.SetMaxOpenConns(db.config.Database.MaxOpenedConnections)

	db.conn = dbConn
	// Migrate database.
	db.Migrate()

	return nil
}
