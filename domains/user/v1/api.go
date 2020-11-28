package userv1

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/sqsinformatique/rosseti-innovation-back/internal/echo-swagger"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/httpsrv"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/logger"
	"github.com/sqsinformatique/rosseti-innovation-back/models"
	"github.com/sqsinformatique/rosseti-innovation-back/types"
)

func (u *UserV1) userPostHandler(ec echo.Context) (err error) {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		return nil
	}

	// Main code of handler
	hndlLog := logger.HandlerLogger(&u.log, ec)

	var userCreds models.NewCredentials

	var bodyBytes []byte
	if ec.Request().Body != nil {
		bodyBytes, err = ioutil.ReadAll(ec.Request().Body)

		ec.Request().Body.Close()

		if err != nil {
			hndlLog.Err(err).Msg("BAD REQUEST")

			return ec.JSON(
				http.StatusBadRequest,
				httpsrv.BadRequest(err),
			)
		}
	}

	err = json.Unmarshal(bodyBytes, &userCreds)
	if err != nil {
		hndlLog.Err(err).Msg("BAD REQUEST")

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	err = userCreds.Validate()
	if err != nil {
		hndlLog.Err(err).Msg("BAD REQUEST")

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	userData, err := u.CreateUser(&userCreds)
	if err != nil {
		hndlLog.Err(err).Msgf("CREATE USER FAILED %s", &userCreds)

		return ec.JSON(
			http.StatusConflict,
			httpsrv.CreateFailed(err),
		)
	}

	return ec.JSON(
		http.StatusOK,
		UserDataResult{Body: userData},
	)
}

func (u *UserV1) userGetHandler(ec echo.Context) (err error) {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		return nil
	}

	// Main code of handler
	hndlLog := logger.HandlerLogger(&u.log, ec)

	userID, err := strconv.ParseInt(ec.Param("id"), 10, 64)
	if err != nil {
		hndlLog.Err(err).Msgf("BAD REQUEST, id %s", ec.Param("id"))

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	userData, err := u.GetUserByID(userID)
	if err != nil {
		hndlLog.Err(err).Msgf("NOT FOUND, id %d", userID)

		return ec.JSON(
			http.StatusNotFound,
			httpsrv.NotFound(err),
		)
	}

	return ec.JSON(
		http.StatusOK,
		UserDataResult{Body: userData},
	)
}

func (u *UserV1) UserPutHandler(ec echo.Context) (err error) {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		return nil
	}

	// Main code of handler
	hndlLog := logger.HandlerLogger(&u.log, ec)

	userID, err := strconv.ParseInt(ec.Param("id"), 10, 64)
	if err != nil {
		hndlLog.Err(err).Msgf("BAD REQUEST, id %s", ec.Param("id"))

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	var bodyBytes []byte
	if ec.Request().Body != nil {
		bodyBytes, err = ioutil.ReadAll(ec.Request().Body)

		ec.Request().Body.Close()

		if err != nil {
			hndlLog.Err(err).Msgf("USER DATA NOT UPDATED, id %d", userID)

			return ec.JSON(
				http.StatusBadRequest,
				httpsrv.BadRequest(err),
			)
		}
	}

	userData, err := u.UpdateUserByID(userID, &bodyBytes)
	if err != nil {
		hndlLog.Err(err).Msgf("BAD REQUEST, id %d, body %s", userID, string(bodyBytes))

		return ec.JSON(
			http.StatusConflict,
			httpsrv.NotUpdated(err),
		)
	}

	return ec.JSON(
		http.StatusOK,
		UserDataResult{Body: userData},
	)
}

func (u *UserV1) CredsPutHandler(ec echo.Context) (err error) {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		return nil
	}

	// Main code of handler
	hndlLog := logger.HandlerLogger(&u.log, ec)

	userID, err := strconv.ParseInt(ec.Param("id"), 10, 64)
	if err != nil {
		hndlLog.Err(err).Msgf("BAD REQUEST, id %s", ec.Param("id"))

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	var userCreds models.UpdateCredentials

	err = ec.Bind(&userCreds)
	if err != nil {
		hndlLog.Err(err).Msgf("BAD REQUEST, id %d", userID)

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	err = userCreds.Validate()
	if err != nil {
		hndlLog.Err(err).Msgf("BAD REQUEST, id %d, userCreds %s", userID, &userCreds)

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	userData, err := u.UpdateUserCredsByID(userID, &userCreds)
	if err != nil {
		hndlLog.Err(err).Msgf("DATA NOT UPDATED, id %d, userCreds %s", userID, &userCreds)

		return ec.JSON(
			http.StatusConflict,
			httpsrv.NotUpdated(err),
		)
	}

	return ec.JSON(
		http.StatusOK,
		UserDataResult{Body: userData},
	)
}

func (u *UserV1) CredsPostHandler(ec echo.Context) (err error) {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		return nil
	}

	// Main code of handler
	hndlLog := logger.HandlerLogger(&u.log, ec)

	var userCreds models.Credentials

	err = ec.Bind(&userCreds)
	if err != nil {
		hndlLog.Err(err).Msg("BAD REQUEST")

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	err = userCreds.Validate()
	if err != nil {
		hndlLog.Err(err).Msgf("BAD REQUEST, userCreds %s", &userCreds)

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	userData, err := u.GetUserDataByCreds(&userCreds)
	if err != nil {
		hndlLog.Err(err).Msgf("UNAUTHORIZED, userCreds %s", &userCreds)

		return ec.JSON(
			http.StatusUnauthorized,
			httpsrv.Unauthorized(err),
		)
	}

	return ec.JSON(
		http.StatusOK,
		UserDataResult{Body: userData},
	)
}

func (u *UserV1) UserDeleteHandler(ec echo.Context) (err error) {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		return nil
	}

	// Main code of handler
	hndlLog := logger.HandlerLogger(&u.log, ec)

	userID, err := strconv.ParseInt(ec.Param("id"), 10, 64)
	if err != nil {
		hndlLog.Err(err).Msgf("BAD REQUEST, id %s", ec.Param("id"))

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	hard := ec.QueryParam("hard")
	if hard == "true" {
		err = u.HardDeleteUserByID(userID)
	} else {
		err = u.SoftDeleteUserByID(userID)
	}

	if err != nil {
		hndlLog.Err(err).Msgf("DATA NOT DELETED, id %d", userID)

		return ec.JSON(
			http.StatusConflict,
			httpsrv.NotDeleted(err),
		)
	}

	return ec.JSON(
		http.StatusOK,
		httpsrv.OkResult(),
	)
}

func (u *UserV1) authPostHandler(ec echo.Context) (err error) {
	// Swagger
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("authPostHandler").
			SetSummary("Authorization user").
			AddInBodyParameter("credentials", "User credentials", &models.Credentials{}, false).
			AddResponse(http.StatusOK, "Test", &SessionDataResult{Body: &models.Session{}})
		return nil
	}

	// Main code of handler
	hndlLog := logger.HandlerLogger(&u.log, ec)

	var cred models.Credentials
	err = ec.Bind(&cred)
	if err != nil {
		hndlLog.Err(err).Msgf("GET USER FAILED %+v", &cred)

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	err = cred.Validate()
	if err != nil {
		hndlLog.Err(err).Msgf("GET USER FAILED %+v", &cred)

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	data, err := u.GetUserDataByCreds(&cred)
	if err != nil {
		hndlLog.Err(err).Msgf("GET USER FAILED %+v", &cred)

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	session, err := u.sessionV1.CreateSession(data.ID)
	if err != nil {
		hndlLog.Err(err).Msgf("CREATE SESSION FAILED %+v", &cred)

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	cookie := new(http.Cookie)
	cookie.Name = "rosseti-session"
	cookie.Value = session.ID
	cookie.Expires = time.Now().Add(24 * time.Hour)

	return ec.JSON(
		http.StatusOK,
		SessionDataResult{Body: session},
	)
}

func (u *UserV1) Introspect(next echo.HandlerFunc, minRole types.Role) echo.HandlerFunc {
	return func(ec echo.Context) error {
		if !u.enableintrospect {
			return next(ec)
		}

		// Main code of handler
		hndlLog := logger.HandlerLogger(&u.log, ec)

		idCookie, err := ec.Cookie("rosseti-session")
		if err != nil {
			hndlLog.Err(err).Msg("GET SESSION FAILED")

			return ec.JSON(
				http.StatusBadRequest,
				httpsrv.BadRequest(err),
			)
		}

		session, err := u.sessionV1.GetSession(idCookie.Value)
		if err != nil {
			hndlLog.Err(err).Msgf("GET SESSION FAILED %s", idCookie.Value)

			return ec.JSON(
				http.StatusBadRequest,
				httpsrv.BadRequest(err),
			)
		}

		hndlLog.Debug().Msgf("good session for userID %d", session.UserID)

		user, err := u.GetUserByID(int64(session.UserID))
		if err != nil {
			hndlLog.Err(err).Msgf("GET USER FAILED %s", idCookie.Value)

			return ec.JSON(
				http.StatusBadRequest,
				httpsrv.BadRequest(err),
			)
		}

		if user.Role < minRole {
			err = errors.New("restricted access to user")
			hndlLog.Err(err).Msgf("RESTRICTED ACCESS to USER %d", session.UserID)

			return ec.JSON(
				http.StatusForbidden,
				httpsrv.Forbidden(err),
			)
		}

		return next(ec)
	}
}
