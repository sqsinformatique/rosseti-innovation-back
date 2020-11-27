package userv1

import (
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	sessionv1 "github.com/sqsinformatique/rosseti-innovation-back/domains/session/v1"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/context"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/httpsrv"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/orm"
	"github.com/sqsinformatique/rosseti-innovation-back/types"
)

type empty struct{}

type UserV1 struct {
	log              zerolog.Logger
	db               **sqlx.DB
	orm              *orm.ORM
	publicV1         *echo.Group
	sessionV1        *sessionv1.SessionV1
	enableintrospect bool
}

func NewUserV1(ctx *context.Context, orm *orm.ORM, sessionV1 *sessionv1.SessionV1) (*UserV1, error) {
	if ctx == nil || orm == nil {
		return nil, errors.New("empty context or orm client")
	}

	u := &UserV1{}
	u.log = ctx.GetPackageLogger(empty{})
	u.publicV1 = ctx.GetHTTPGroup(httpsrv.PublicSrv, httpsrv.V1)
	u.db = ctx.GetDatabase()
	u.orm = orm
	u.sessionV1 = sessionV1
	u.enableintrospect = ctx.Config.Introspection.Enable

	u.publicV1.POST("/auth", u.authPostHandler)

	u.publicV1.POST("/user", u.userPostHandler)
	u.publicV1.GET("/users/:id", u.Introspect(u.userGetHandler, types.User))
	u.publicV1.PUT("/users/:id", u.Introspect(u.UserPutHandler, types.Admin))
	u.publicV1.PUT("/credentials/:id", u.CredsPutHandler)
	u.publicV1.POST("/credentials", u.CredsPostHandler)
	u.publicV1.DELETE("/users/:id", u.Introspect(u.UserDeleteHandler, types.Admin))

	return u, nil
}
