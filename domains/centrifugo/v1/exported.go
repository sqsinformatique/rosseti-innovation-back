package centrifugov1

import (
	"bytes"
	"container/list"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	sessionv1 "github.com/sqsinformatique/rosseti-innovation-back/domains/session/v1"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/cfg"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/context"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/httpsrv"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/orm"
	"github.com/sqsinformatique/rosseti-innovation-back/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type empty struct{}

type CentrifugoV1 struct {
	log                 zerolog.Logger
	privateV1           *echo.Group
	publicV1            *echo.Group
	sessionV1           *sessionv1.SessionV1
	config              *cfg.AppCfg
	orm                 *orm.ORM
	mongoDB             **mongo.Client
	db                  **sqlx.DB
	lastActiveThemes    list.List
	lastActiveThemesMap map[string]string
}

func NewCentrifugoV1(ctx *context.Context, orm *orm.ORM, sessionV1 *sessionv1.SessionV1) (*CentrifugoV1, error) {
	if ctx == nil {
		return nil, errors.New("empty context or orm client")
	}

	c := &CentrifugoV1{}
	c.log = ctx.GetPackageLogger(empty{})
	c.privateV1 = ctx.GetHTTPGroup(httpsrv.PrivateSrv, httpsrv.V1)
	c.publicV1 = ctx.GetHTTPGroup(httpsrv.PublicSrv, httpsrv.V1)
	c.sessionV1 = sessionV1
	c.config = ctx.Config
	c.mongoDB = ctx.GetMongoDB()
	c.orm = orm
	c.db = ctx.GetDatabase()

	c.lastActiveThemesMap = make(map[string]string)

	c.privateV1.POST("/centrifugo/connect", c.AuthConnectHandler)
	c.publicV1.POST("/centrifugo/publish", c.PublishHandler)
	c.publicV1.GET("/centrifugo/chat/:id", c.GetHistoryHandler)
	c.publicV1.POST("/themes", c.CreateThemeHandler)
	c.publicV1.GET("/directionsdetailed", c.GetDirectionsDetailedHandler)
	c.publicV1.GET("/directions", c.GetDirectionsHandler)
	c.publicV1.GET("/lastactivethems", c.GetLastActiveThemes)
	c.publicV1.PUT("/themes/:id/like", c.PutLikeThemes)

	return c, nil
}

func (c *CentrifugoV1) AddLastActive(pub *models.Publish) error {
	if pub.Type != "theme" {
		return nil
	}

	channelID, err := strconv.Atoi(pub.Channel)
	if err != nil {
		return err
	}

	c.log.Debug().Msgf("last active theme: %s", pub.Channel)

	if c.lastActiveThemes.Len() < 10 {
		c.lastActiveThemes.PushFront(channelID)
		c.lastActiveThemesMap[pub.Channel] = "saved"
	} else {
		e := c.lastActiveThemes.Front()
		c.lastActiveThemes.Remove(e)
		delete(c.lastActiveThemesMap, e.Value.(string))

		c.lastActiveThemes.PushFront(channelID)
		c.lastActiveThemesMap[pub.Channel] = "saved"
	}

	return nil
}

func (c *CentrifugoV1) Publish(pub *models.Publish, userID int) error {
	client := &http.Client{}

	centrifugoReq := &models.CentrifugoAPIRequest{
		Method: "publish",
		Params: &models.CentrifugoParams{
			Channel: pub.Channel,
		},
	}
	centrifugoReqData := make(map[string]interface{})
	centrifugoReqData["message"] = pub.Message
	centrifugoReqData["sender"] = userID
	centrifugoReq.Params.Data = centrifugoReqData

	jsonData, err := json.Marshal(centrifugoReq)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, c.config.Centrifugo.DSN, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New("failed handle request on centrinfugo side")
	}
	channelID, err := strconv.Atoi(pub.Channel)
	if err != nil {
		return err
	}

	err = c.SaveToDB(channelID, userID, "", pub.Message)
	if err != nil {
		return err
	}

	return nil
}
