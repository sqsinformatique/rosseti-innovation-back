package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/rs/zerolog"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/cfg"
	intctx "github.com/sqsinformatique/rosseti-innovation-back/internal/context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type empty struct{}

type DB struct {
	log    zerolog.Logger
	conn   *mongo.Client
	config *cfg.AppCfg
	cancel context.CancelFunc
}

// Initialize initialize database
func NewMongoDB(ctx *intctx.Context) (*DB, error) {
	if ctx == nil || ctx.Config == nil {
		return nil, errors.New("empty context or config")
	}
	db := &DB{}
	db.config = ctx.Config
	db.log = ctx.GetPackageLogger(empty{})
	ctx.RegisterMongoDB(&db.conn)

	return db, nil
}

// Start start db
func (db *DB) Start() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	dbConn, err := mongo.Connect(ctx, options.Client().ApplyURI(db.config.Mongo.DSN))
	if err != nil {
		db.log.Fatal().Err(err).Msgf("Failed to connect to Mongo database")
	}

	db.log.Info().Msg("Mongo Database connection established")

	db.conn = dbConn
	db.cancel = cancel

	return nil
}

func (db *DB) Stop() {
	db.cancel()
}
