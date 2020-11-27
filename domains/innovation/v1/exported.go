package innovationv1

import (
	"errors"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	profilev1 "github.com/sqsinformatique/rosseti-innovation-back/domains/profile/v1"
	userv1 "github.com/sqsinformatique/rosseti-innovation-back/domains/user/v1"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/cfg"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/context"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/httpsrv"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/orm"
	"github.com/sqsinformatique/rosseti-innovation-back/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type empty struct{}

type InnovationV1 struct {
	log       zerolog.Logger
	cfg       *cfg.AppCfg
	db        **sqlx.DB
	mongodb   **mongo.Client
	elasticDB **elasticsearch.Client
	orm       *orm.ORM
	profilev1 *profilev1.ProfileV1
	publicV1  *echo.Group
	userV1    *userv1.UserV1
}

func NewInnovationV1(ctx *context.Context,
	profilev1 *profilev1.ProfileV1,
	orm *orm.ORM,
	userV1 *userv1.UserV1,
) (*InnovationV1, error) {
	if ctx == nil || profilev1 == nil || orm == nil {
		return nil, errors.New("empty context or profilev1 client or orm client")
	}

	inn := &InnovationV1{}
	inn.log = ctx.GetPackageLogger(empty{})
	inn.publicV1 = ctx.GetHTTPGroup(httpsrv.PublicSrv, httpsrv.V1)
	inn.profilev1 = profilev1
	inn.cfg = ctx.Config
	inn.db = ctx.GetDatabase()
	inn.mongodb = ctx.GetMongoDB()
	inn.elasticDB = ctx.GetElasticDB()
	inn.userV1 = userV1
	inn.orm = orm

	inn.publicV1.POST("/innovations", inn.userV1.Introspect(inn.innovationPostHandler, types.User))
	inn.publicV1.POST("/innovations/search", inn.userV1.Introspect(inn.searchPostHandler, types.User))
	// a.publicV1.GET("/innovations/:actid", a.userV1.Introspect(a.actGetHandler, types.User))
	// a.publicV1.GET("/innovations/staff/:id", a.userV1.Introspect(a.innovationsByStaffIDGetHandler, types.User))
	// a.publicV1.GET("/innovations/superviser/:id", a.userV1.Introspect(a.innovationsBySuperviserIDGetHandler, types.User))
	// a.publicV1.PUT("/innovations/:actid", a.userV1.Introspect(a.ActPutHandler, types.User))
	// a.publicV1.POST("/innovations/:actid/images", a.userV1.Introspect(a.actPostImagesHandler, types.User))
	// a.publicV1.GET("/innovations/:actid/images/:id", a.userV1.Introspect(a.actGetImageHandler, types.User))
	// a.publicV1.DELETE("/innovations/:actid", a.userV1.Introspect(a.ActDeleteHandler, types.Moderator))
	// a.publicV1.POST("/innovations/:actid/signsupervisor", a.userV1.Introspect(a.innovationsignSuperviserPostHandler, types.Moderator))
	// a.publicV1.POST("/innovations/:actid/signstaff", a.userV1.Introspect(a.innovationsignStaffPostHandler, types.User))
	// a.publicV1.POST("/innovations/:actid/revoke", a.userV1.Introspect(a.ActRevertedSuperviserPostHandler, types.Moderator))

	return inn, nil
}
