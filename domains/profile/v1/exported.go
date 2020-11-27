package profilev1

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	userv1 "github.com/sqsinformatique/rosseti-innovation-back/domains/user/v1"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/context"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/crypto"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/httpsrv"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/orm"
	"github.com/sqsinformatique/rosseti-innovation-back/types"
)

type empty struct{}

type ProfileV1 struct {
	log      zerolog.Logger
	db       **sqlx.DB
	orm      *orm.ORM
	publicV1 *echo.Group
	userV1   *userv1.UserV1
}

func NewProfileV1(ctx *context.Context, orm *orm.ORM, userV1 *userv1.UserV1) (*ProfileV1, error) {
	if ctx == nil || orm == nil {
		return nil, errors.New("empty context or orm client")
	}

	p := &ProfileV1{}
	p.log = ctx.GetPackageLogger(empty{})
	p.publicV1 = ctx.GetHTTPGroup(httpsrv.PublicSrv, httpsrv.V1)
	p.userV1 = userV1
	p.db = ctx.GetDatabase()
	p.orm = orm

	p.publicV1.POST("/profiles", p.userV1.Introspect(p.ProfilePostHandler, types.User))
	p.publicV1.GET("/profiles/:id", p.userV1.Introspect(p.ProfileGetHandler, types.User))
	p.publicV1.GET("/profilessearch", p.userV1.Introspect(p.ProfileSearchGetHandler, types.User))
	p.publicV1.PUT("/profiles/:id", p.userV1.Introspect(p.ProfilePutHandler, types.User))
	p.publicV1.DELETE("/profiles/:id", p.userV1.Introspect(p.ProfileDeleteHandler, types.Admin))

	return p, nil
}

func (o *ProfileV1) SignDataByID(id int64, data interface{}) (string, error) {
	profile, err := o.GetProfileByID(id)
	if err != nil {
		return "", err
	}

	key, err := crypto.UnmarshalPrivate(profile.PrivateKey)
	if err != nil {
		return "", fmt.Errorf("failed unmarhal privatekey: %w", err)
	}

	return crypto.DataSign(data, key)
}
