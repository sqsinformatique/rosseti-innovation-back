package elastic

import (
	"errors"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/rs/zerolog"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/cfg"
	intctx "github.com/sqsinformatique/rosseti-innovation-back/internal/context"
)

type empty struct{}

type DB struct {
	log    zerolog.Logger
	conn   *elasticsearch.Client
	config *cfg.AppCfg
}

// Initialize initialize database
func NewElasticDB(ctx *intctx.Context) (*DB, error) {
	if ctx == nil || ctx.Config == nil {
		return nil, errors.New("empty context or config")
	}
	db := &DB{}
	db.config = ctx.Config
	db.log = ctx.GetPackageLogger(empty{})
	ctx.RegisterElasticDB(&db.conn)

	return db, nil
}

func (db *DB) Start() error {
	cfg := elasticsearch.Config{
		Addresses: []string{
			db.config.Elastic.DSN,
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return err
	}

	db.conn = es

	return nil
}
