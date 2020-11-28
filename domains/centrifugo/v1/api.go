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

	err = c.AddLastActive(&pub)
	if err != nil {
		hndlLog.Err(err).Msg("FAILED ADD TO LAST ACTIVE")

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

func (c *CentrifugoV1) GetHistoryHandler(ec echo.Context) error {
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("getHistoryHandler").
			SetSummary("Get chat history").
			AddInPathParameter("id", "Chat id", reflect.Int64).
			AddInHeaderParameter("Authorization", "Authorization header", reflect.String, true).
			AddResponse(http.StatusOK, "OK", &ChatDataResult{Body: &models.ChatChannel{}})
		return nil
	}

	// Main code of handler
	hndlLog := logger.HandlerLogger(&c.log, ec)
	chatID, err := strconv.ParseInt(ec.Param("id"), 10, 64)
	if err != nil {
		hndlLog.Err(err).Msgf("BAD REQUEST, id %s", ec.Param("id"))

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	chatData, err := c.GetChat(int(chatID))
	if err != nil {
		hndlLog.Err(err).Msgf("Failed get chat, id %s", ec.Param("id"))

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	return ec.JSON(
		http.StatusOK,
		ChatDataResult{Body: chatData},
	)
}

func (c *CentrifugoV1) CreateThemeHandler(ec echo.Context) (err error) {
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("CreateThemeHandler").
			SetSummary("Create Theme").
			AddInBodyParameter("theme", "New Theme", &models.Theme{}, true).
			AddInHeaderParameter("Authorization", "Authorization header", reflect.String, true).
			AddResponse(http.StatusOK, "OK", &ThemeDataResult{Body: &models.Theme{}})
		return nil
	}

	// Main code of handler
	hndlLog := logger.HandlerLogger(&c.log, ec)

	var theme models.Theme
	err = ec.Bind(&theme)
	if err != nil {
		hndlLog.Err(err).Msgf("CREATE THEME FAILED %+v", &theme)

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	themeData, err := c.CreateTheme(&theme)
	if err != nil {
		hndlLog.Err(err).Msgf("CREATE THEME FAILED %+v", &theme)

		return ec.JSON(
			http.StatusConflict,
			httpsrv.CreateFailed(err),
		)
	}

	return ec.JSON(
		http.StatusOK,
		ThemeDataResult{Body: themeData},
	)
}

func (c *CentrifugoV1) GetDirectionsHandler(ec echo.Context) (err error) {
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("GetDirectionsHandler").
			SetSummary("Get directions").
			AddInHeaderParameter("Authorization", "Authorization header", reflect.String, true).
			AddResponse(http.StatusOK, "OK", &DirectionsDataResult{Body: &ArrayOfDirectionData{}})
		return nil
	}

	// Main code of handler
	hndlLog := logger.HandlerLogger(&c.log, ec)

	directionsData, err := c.SelectAllDirections()
	if err != nil {
		hndlLog.Err(err).Msg("SELECT ALL DIRECTIONS FAILED")

		return ec.JSON(
			http.StatusConflict,
			httpsrv.CreateFailed(err),
		)
	}

	return ec.JSON(
		http.StatusOK,
		DirectionsDataResult{Body: directionsData},
	)
}

func (c *CentrifugoV1) GetDirectionsDetailedHandler(ec echo.Context) (err error) {
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("GetDirectionsDetailedHandler").
			SetSummary("Get directions").
			AddInHeaderParameter("Authorization", "Authorization header", reflect.String, true).
			AddResponse(http.StatusOK, "OK", &DirectionsDataResult{Body: &ArrayOfDirectionDetailedData{}})
		return nil
	}

	// Main code of handler
	hndlLog := logger.HandlerLogger(&c.log, ec)

	directionsData, err := c.SelectAllDirections()
	if err != nil {
		hndlLog.Err(err).Msg("SELECT ALL DIRECTIONS FAILED")

		return ec.JSON(
			http.StatusConflict,
			httpsrv.CreateFailed(err),
		)
	}

	detailed := ArrayOfDirectionDetailedData{}
	for _, v := range *directionsData {
		item := models.DirectionDetailed{}
		themes, err := c.SelectThemesByDirection(v.ID)
		if err != nil {
			hndlLog.Err(err).Msgf("SELECT THEMES FAILED, direction %d", v.ID)

			return ec.JSON(
				http.StatusConflict,
				httpsrv.CreateFailed(err),
			)
		}

		item.Themes = *themes
		item.Direction = v

		detailed = append(detailed, item)
	}

	return ec.JSON(
		http.StatusOK,
		DirectionsDataResult{Body: detailed},
	)
}

func (c *CentrifugoV1) GetLastActiveThemes(ec echo.Context) (err error) {
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("GetLastActiveThemes").
			SetSummary("Get last active themes").
			AddInHeaderParameter("Authorization", "Authorization header", reflect.String, true).
			AddResponse(http.StatusOK, "OK", &ArrayThemeDataResult{Body: &ArrayOfThemesData{}})
		return nil
	}

	// Main code of handler
	hndlLog := logger.HandlerLogger(&c.log, ec)

	themesData, err := c.SelectLastActiveThemes()
	if err != nil {
		hndlLog.Err(err).Msgf("SELECT LAST ACTIVE THEMES FAILED, direction")

		return ec.JSON(
			http.StatusConflict,
			httpsrv.CreateFailed(err),
		)
	}

	return ec.JSON(
		http.StatusOK,
		ArrayThemeDataResult{Body: themesData},
	)
}

func (c *CentrifugoV1) PutLikeThemes(ec echo.Context) (err error) {
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("PutLikeThemes").
			SetSummary("Like themes").
			AddInPathParameter("id", "Theme id", reflect.Int64).
			AddResponse(http.StatusOK, "OK", httpsrv.OkResult())
		return nil
	}

	// Main code of handler
	hndlLog := logger.HandlerLogger(&c.log, ec)
	themeID, err := strconv.ParseInt(ec.Param("id"), 10, 64)
	if err != nil {
		hndlLog.Err(err).Msgf("BAD REQUEST, id %s", ec.Param("id"))

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	err = c.LikeTheme(themeID)
	if err != nil {
		hndlLog.Err(err).Msgf("FAILED LIKE THEME, id %s", ec.Param("id"))

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	return ec.JSON(
		http.StatusOK,
		httpsrv.OkResult(),
	)
}
