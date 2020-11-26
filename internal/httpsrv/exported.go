package httpsrv

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/cfg"
	internalctx "github.com/sqsinformatique/rosseti-innovation-back/internal/context"
	echoSwagger "github.com/sqsinformatique/rosseti-innovation-back/internal/echo-swagger"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/swagger"
)

const (
	defaultBodyReadTimeout   = 10
	defaultBodyWriteTimeout  = 10
	defaultHeaderReadTimeout = 5

	PublicSrv  = "public"
	PrivateSrv = "private"
	V1         = "1"
)

type empty struct{}

type HTTPSrv struct {
	PublicSrv  *echo.Echo
	PrivateSrv *echo.Echo

	PublicV1  *echo.Group
	PrivateV1 *echo.Group

	log    zerolog.Logger
	config *cfg.AppCfg
}

func initializeSrv(ctx *internalctx.Context, name string) *echo.Echo {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Debug = false
	e.DisableHTTP2 = true
	e.HideBanner = true
	e.HidePort = true

	e.Server.ReadHeaderTimeout = time.Second * time.Duration(defaultHeaderReadTimeout)
	e.Server.ReadTimeout = time.Second * time.Duration(defaultBodyReadTimeout)
	e.Server.WriteTimeout = time.Second * time.Duration(defaultBodyWriteTimeout)

	// CORS default
	// Allows requests from any origin wth GET, HEAD, PUT, POST or DELETE method.
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "X-Request-Id"},
	}))

	ctx.RegisterHTTPServer(name, e)
	return e
}

func initializeGroup(ctx *internalctx.Context, srvName, group string) *echo.Group {
	srv := ctx.HTTPServers[srvName]
	gr := srv.Group("/api/v" + group)

	ctx.RegisterHTTPGroup(srvName, group, gr)
	return gr
}

func (h *HTTPSrv) BuildSwagger() {
	// Build public swagger
	err := echoSwagger.BuildSwagger(
		h.PublicSrv,
		"/swagger/*",
		h.config.PublicHTTP.Listen,
		swagger.NewSwagger().
			SetBasePath("/api/v1/").
			SetInfo(swagger.NewInfo().
				SetTitle("Public API V1").
				SetDescription("This is a documentation of Rosseti Innovation Back Public API V1").
				SetTermOfService("http://swagger.io/terms/").
				SetContact(swagger.NewContact()).
				SetLicense(swagger.NewLicense())),
		nil,
	)

	if err != nil {
		h.log.Fatal().Err(err).Msg("Failed to build Rosseti Innovation Back Public API V1")
	}

	// Build private swagger
	err = echoSwagger.BuildSwagger(
		h.PrivateSrv,
		"/swagger/*",
		h.config.PrivateHTTP.Listen,
		swagger.NewSwagger().
			SetBasePath("/api/v1/").
			SetInfo(swagger.NewInfo().
				SetTitle("Private API V1").
				SetDescription("This is a documentation of Rosseti Innovation Back Private API V1").
				SetTermOfService("http://swagger.io/terms/").
				SetContact(swagger.NewContact()).
				SetLicense(swagger.NewLicense())),
		nil,
	)

	if err != nil {
		h.log.Fatal().Err(err).Msg("Failed to build Rosseti Innovation Back Private API V1")
	}
}

func NewHTTPSrv(ctx *internalctx.Context) (*HTTPSrv, error) {
	if ctx == nil || ctx.Config == nil {
		return nil, errors.New("empty context or config")
	}

	h := &HTTPSrv{}
	h.config = ctx.Config
	h.log = ctx.GetPackageLogger(empty{})

	h.PublicSrv = initializeSrv(ctx, PublicSrv)
	h.PrivateSrv = initializeSrv(ctx, PrivateSrv)

	// API V1
	h.PublicV1 = initializeGroup(ctx, PublicSrv, V1)
	h.PrivateV1 = initializeGroup(ctx, PrivateSrv, V1)

	return h, nil
}

func (h *HTTPSrv) start(srv *echo.Echo, address string) {
	go func() {
		h.log.Info().Str("address", address).Msg("Starting server...")

		err := srv.Start(address)
		if !strings.Contains(err.Error(), "Server closed") {
			h.log.Error().Err(err).Msg("HTTP server critical error occurred")
		}
	}()
}

func (h *HTTPSrv) Start() {
	h.start(h.PublicSrv, h.config.PublicHTTP.Listen)
	h.start(h.PrivateSrv, h.config.PrivateHTTP.Listen)
}

func shutdown(srv *echo.Echo) error {
	err := srv.Shutdown(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (h *HTTPSrv) Shutdown() error {
	if err := shutdown(h.PublicSrv); err != nil {
		return err
	}

	if err := shutdown(h.PrivateSrv); err != nil {
		return err
	}

	return nil
}
