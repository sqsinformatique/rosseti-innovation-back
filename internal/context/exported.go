package context

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/cfg"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

type Context struct {
	Log         zerolog.Logger
	DB          **sqlx.DB
	MongoDB     **mongo.Client
	Config      *cfg.AppCfg
	HTTPServers map[string]*echo.Echo
	HTTPGroups  map[string]*echo.Group
}

func NewContext() *Context {
	return &Context{}
}

func (ctx *Context) RegisterConfig(config *cfg.AppCfg) {
	ctx.Config = config
}

func (ctx *Context) RegisterLogger() {
	ctx.Log = logger.NewLogger(ctx.Config)
}

func (ctx *Context) GetPackageLogger(emptyStruct interface{}) (log zerolog.Logger) {
	return logger.InitializeLogger(&ctx.Log, emptyStruct)
}

func (ctx *Context) RegisterDatabase(db **sqlx.DB) {
	ctx.DB = db
}

func (ctx *Context) RegisterMongoDB(mongodb **mongo.Client) {
	ctx.MongoDB = mongodb
}

func (ctx *Context) GetDatabase() **sqlx.DB {
	return ctx.DB
}

func (ctx *Context) GetMongoDB() **mongo.Client {
	return ctx.MongoDB
}

func (ctx *Context) RegisterHTTPServer(name string, srv *echo.Echo) {
	if ctx.HTTPServers == nil {
		ctx.HTTPServers = make(map[string]*echo.Echo)
	}

	ctx.HTTPServers[name] = srv
}

func (ctx *Context) GetHTTPServer(name string) *echo.Echo {
	return ctx.HTTPServers[name]
}

func (ctx *Context) RegisterHTTPGroup(srvName, grName string, gr *echo.Group) {
	if ctx.HTTPGroups == nil {
		ctx.HTTPGroups = make(map[string]*echo.Group)
	}

	ctx.HTTPGroups[srvName+"+"+grName] = gr
}

func (ctx *Context) GetHTTPGroup(srvName, grName string) *echo.Group {
	return ctx.HTTPGroups[srvName+"+"+grName]
}
