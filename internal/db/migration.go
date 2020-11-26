package db

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pressly/goose"
	"github.com/rs/zerolog"
)

const (
	migrationActionEnv     = "PG_MIGRATE_ACTION"
	migrationsCountEnv     = "PG_MIGRATE_COUNT"
	migrationsDirectoryEnv = "PG_MIGRATE_DIR"
	migrateOnlyEnv         = "PG_MIGRATE_ONLY"

	actionNothing = "nothing"
	actionUp      = "up"
	actionDown    = "down"
)

// MigrationsType represents enumerator for acceptable migration types.
type MigrationsType int

const (
	MigrationTypeSQLFiles MigrationsType = iota
	MigrationTypeGoCode
)

func isMigrateOnly() bool {
	// Figure out was migrate-only mode requested?
	migrateOnlyRaw, migrateOnlyFound := os.LookupEnv(migrateOnlyEnv)
	if migrateOnlyFound && migrateOnlyRaw != "" {
		return true
	}

	return false
}

func getMigrationDir(mType MigrationsType) (string, error) {
	migrationsDirectory, migrationsDirectoryFound := os.LookupEnv(migrationsDirectoryEnv)
	if !migrationsDirectoryFound && mType == MigrationTypeSQLFiles {
		return "", fmt.Errorf("migrations directory isn't defined, please define it in '%s' environment variable",
			migrationsDirectoryEnv)
	}

	// If we're migrate database using Go functions - force migrations
	// directory to be "." regardless of environment variable value.
	if mType == MigrationTypeGoCode {
		migrationsDirectory = "."
	}

	return migrationsDirectory, nil
}

func getMigrationCount() (int64, error) {
	countRaw, countFound := os.LookupEnv(migrationsCountEnv)
	if !countFound {
		return 0, nil
	}

	count, err := strconv.ParseInt(countRaw, 10, 64)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func getMigrationAction() string {
	actionRaw, actionFound := os.LookupEnv(migrationActionEnv)
	if !actionFound {
		actionRaw = actionNothing
	}

	return strings.ToLower(actionRaw)
}

func (db *DB) setMigrationFlag() {
	os.Setenv(migrationActionEnv, actionNothing)
}

func (db *DB) getCurrentDBVersion(log *zerolog.Logger) int64 {
	currentDBVersion, gooseerr := goose.GetDBVersion(db.conn.DB)
	if gooseerr != nil {
		log.Fatal().Err(gooseerr).Msg("Failed to get database version")
	}

	return currentDBVersion
}

func (db *DB) migrateSchema() error {
	if db.schemaName == "" {
		return nil
	}

	_, err := db.conn.Exec(db.conn.Rebind("CREATE SCHEMA IF NOT EXISTS " + db.schemaName))

	return err
}

// Migrates database.
func (db *DB) Migrate() {
	migrationsLog := db.log.With().Str("subsystem", "database migrations").Logger()

	err := db.migrateSchema()
	if err != nil {
		migrationsLog.Error().Err(err)
		db.setMigrationFlag()

		return
	}

	action := getMigrationAction()
	if action == actionNothing {
		migrationsLog.Info().Msg("Database migration wasn't requested")
		db.setMigrationFlag()

		return
	}

	count, err := getMigrationCount()
	if err != nil {
		migrationsLog.Error().Err(err).Msg("Failed to parse migrations count to apply/rollback. Doing nothing.")
		db.setMigrationFlag()

		return
	}

	migrationsDirectory, err := getMigrationDir(MigrationTypeSQLFiles)
	if err != nil {
		migrationsLog.Error().Err(err)
		db.setMigrationFlag()

		return
	}

	_ = goose.SetDialect("postgres")

	currentDBVersion := db.getCurrentDBVersion(&migrationsLog)
	migrationsLog.Debug().Int64("database version", currentDBVersion).Msg("Current database version obtained")

	switch {
	case action == actionUp && count == 0:
		migrationsLog.Info().Msg("Applying all unapplied migrations...")

		err = goose.Up(db.conn.DB, migrationsDirectory)
	case action == actionUp && count != 0:
		newVersion := currentDBVersion + count

		migrationsLog.Info().Int64("new version", newVersion).Msg("Migrating database to specific version")

		err = goose.UpTo(db.conn.DB, migrationsDirectory, newVersion)
	case action == actionDown && count == 0:
		migrationsLog.Info().Msg("Downgrading database to zero state, you'll need to re-apply migrations!")

		err = goose.DownTo(db.conn.DB, migrationsDirectory, 0)

		migrationsLog.Fatal().Msg("Database downgraded to zero state. You have to re-apply migrations")
	case action == actionDown && count != 0:
		newVersion := currentDBVersion - count

		migrationsLog.Info().Int64("new version", newVersion).Msg("Downgrading database to specific version")

		err = goose.DownTo(db.conn.DB, migrationsDirectory, newVersion)
	default:
		migrationsLog.Fatal().Str("action", action).Int64("count", count).Msg("Unsupported set of migration parameters, cannot continue")
	}

	if err != nil {
		migrationsLog.Fatal().Err(err).Msg("Failed to execute migration sequence")
	}

	migrationsLog.Info().Msg("Database migrated successfully")
	db.setMigrationFlag()

	// Figure out was migrate-only mode requested?
	if isMigrateOnly() {
		migrationsLog.Warn().Msg("Only database migrations was requested, shutting down")
		os.Exit(0)
	}

	migrationsLog.Info().Msg("Migrate-only mode wasn't requested")
}

func (db *DB) SetSchemaName(name string) {
	db.schemaName = name
	goose.SetTableName(db.schemaName + ".goose_db_version")
}
