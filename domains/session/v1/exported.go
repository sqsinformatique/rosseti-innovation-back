package sessionv1

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/context"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/httpsrv"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/orm"
)

type empty struct{}

type SessionV1 struct {
	log      zerolog.Logger
	db       **sqlx.DB
	orm      *orm.ORM
	publicV1 *echo.Group
}

func NewSessionV1(ctx *context.Context, orm *orm.ORM) (*SessionV1, error) {
	if ctx == nil || orm == nil {
		return nil, errors.New("empty context or orm client")
	}

	s := &SessionV1{}
	s.log = ctx.GetPackageLogger(empty{})
	s.publicV1 = ctx.GetHTTPGroup(httpsrv.PublicSrv, httpsrv.V1)
	s.db = ctx.GetDatabase()
	s.orm = orm

	s.publicV1.DELETE("/sessions", s.SessionDeleteHandler)

	return s, nil
}
