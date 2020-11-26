package centrifugov1

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
	sessionv1 "github.com/sqsinformatique/rosseti-innovation-back/domains/session/v1"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/cfg"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/context"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/httpsrv"
	"github.com/sqsinformatique/rosseti-innovation-back/models"
)

type empty struct{}

type CentrifugoV1 struct {
	log       zerolog.Logger
	privateV1 *echo.Group
	sessionV1 *sessionv1.SessionV1
	config    *cfg.AppCfg
}

func NewCentrifugoV1(ctx *context.Context, sessionV1 *sessionv1.SessionV1) (*CentrifugoV1, error) {
	if ctx == nil {
		return nil, errors.New("empty context or orm client")
	}

	c := &CentrifugoV1{}
	c.log = ctx.GetPackageLogger(empty{})
	c.privateV1 = ctx.GetHTTPGroup(httpsrv.PrivateSrv, httpsrv.V1)
	c.sessionV1 = sessionV1
	c.config = ctx.Config

	c.privateV1.POST("/centrifugo/connect", c.AuthConnectHandler)
	c.privateV1.POST("/centrifugo/publish", c.PublishHandler)

	return c, nil
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

	return nil
}
