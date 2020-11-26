package centrifugov1

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/sqsinformatique/rosseti-innovation-back/internal/echo-swagger"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/httpsrv"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/logger"
	"github.com/sqsinformatique/rosseti-innovation-back/models"
)

var (
	ErrBadAuthRequest = errors.New("bad authorization request")
)

// ExtractToken extract token from request
func ExtractToken(r *http.Request) (string, error) {
	if r == nil {
		return "", ErrBadAuthRequest
	}

	// Get token from query
	if r.URL == nil {
		return "", ErrBadAuthRequest
	}
	queryValues := r.URL.Query()
	token := queryValues.Get("session")

	if token != "" {
		return token, nil
	}

	// Get token from cookie
	tokenCookie, err := r.Cookie("rosseti-session")
	if err == nil {
		token = tokenCookie.Value
	}

	if token != "" {
		return token, nil
	}

	// Token may be in body in data object
	contents, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	js := make(map[string]interface{})
	err = json.Unmarshal(contents, &js)
	if err != nil {
		return "", ErrBadAuthRequest
	}

	var (
		ok   bool
		data map[string]interface{}
		val  interface{}
	)
	if val, ok = js["data"]; ok {
		if data, ok = val.(map[string]interface{}); ok {
			if token, ok = data["session"].(string); ok && token != "" {
				return token, nil
			}
		}
	}

	// If not token not found in query, try get from Authorization header
	token = r.Header.Get("Authorization")

	splitToken := strings.Split(token, " ")
	if len(splitToken) < 2 {
		return "", ErrBadAuthRequest
	}

	switch strings.ToLower(strings.TrimSpace(splitToken[0])) {
	case "bearer", "session":
		token = strings.TrimSpace(splitToken[1])
		return token, nil
	}

	return "", ErrBadAuthRequest
}

// ExtractToken extract token from request
func ExtractToken2(r *http.Request) (string, error) {
	if r == nil {
		return "", ErrBadAuthRequest
	}

	// Get token from query
	if r.URL == nil {
		return "", ErrBadAuthRequest
	}
	queryValues := r.URL.Query()
	token := queryValues.Get("session")

	if token != "" {
		return token, nil
	}

	// Get token from cookie
	tokenCookie, err := r.Cookie("rosseti-session")
	if err == nil {
		token = tokenCookie.Value
	}

	if token != "" {
		return token, nil
	}

	// If not token not found in query, try get from Authorization header
	token = r.Header.Get("Authorization")

	splitToken := strings.Split(token, " ")
	if len(splitToken) < 2 {
		return "", ErrBadAuthRequest
	}

	switch strings.ToLower(strings.TrimSpace(splitToken[0])) {
	case "bearer", "session":
		token = strings.TrimSpace(splitToken[1])
		return token, nil
	}

	return "", ErrBadAuthRequest
}

func (c *CentrifugoV1) AuthConnectHandler(ec echo.Context) (err error) {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		// echoSwagger.AddToSwagger(ec).
		// 	SetProduces("application/json").
		// 	SetDescription("authСщттусеHandler").
		// 	SetSummary("Authorization for connect to centrifugo").
		// 	AddInBodyParameter("credentials", "User credentials", &models.Credentials{}, false).
		// 	AddResponse(http.StatusOK, "Test", &SessionDataResult{Body: &models.Session{}})
		return nil
	}

	// Main code of handler
	hndlLog := logger.HandlerLogger(&c.log, ec)

	sessionID, err := ExtractToken(ec.Request())
	if err != nil {
		hndlLog.Err(err).Msg("GET SESSION FAILED")

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	session, err := c.sessionV1.GetSession(sessionID)
	if err != nil {
		hndlLog.Err(err).Msgf("GET SESSION FAILED %s", sessionID)

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	hndlLog.Debug().Msgf("good session for userID %d", session.UserID)

	return ec.JSON(
		http.StatusOK,
		models.CentrifugoIntrospectionResult{
			Result: &models.CentrifugoIntrospection{
				User: strconv.Itoa(session.UserID),
				// Data: introspection,
			},
		},
	)
}

func (c *CentrifugoV1) PublishHandler(ec echo.Context) (err error) {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("publishHandler").
			SetSummary("Publish to centrifugo").
			AddInBodyParameter("publish", "Request for publish", &models.Publish{}, false).
			AddInHeaderParameter("Authorization", "Authorization header", reflect.String, true).
			AddResponse(http.StatusOK, "OK", httpsrv.OkResult())
		return nil
	}

	// Main code of handler
	hndlLog := logger.HandlerLogger(&c.log, ec)

	sessionID, err := ExtractToken2(ec.Request())
	if err != nil {
		hndlLog.Err(err).Msg("GET SESSION FAILED")

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	session, err := c.sessionV1.GetSession(sessionID)
	if err != nil {
		hndlLog.Err(err).Msgf("GET SESSION FAILED %s", sessionID)

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	hndlLog.Debug().Msgf("good session for userID %d", session.UserID)

	var pub models.Publish

	err = ec.Bind(&pub)
	if err != nil {
		hndlLog.Err(err).Msg("BAD REQUEST")

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	err = c.Publish(&pub, session.UserID)
	if err != nil {
		hndlLog.Err(err).Msg("BAD REQUEST")

		return ec.JSON(
			http.StatusInternalServerError,
			httpsrv.InternalServerError(err),
		)
	}

	return ec.JSON(
		http.StatusOK,
		httpsrv.OkResult(),
	)
}