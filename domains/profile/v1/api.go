package profilev1

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/httpsrv"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/logger"
	"github.com/sqsinformatique/rosseti-innovation-back/models"
)

func (o *ProfileV1) ProfilePostHandler(ec echo.Context) (err error) {
	// Main code of handler
	hndlLog := logger.HandlerLogger(&o.log, ec)

	var profile models.Profile
	err = ec.Bind(&profile)
	if err != nil {
		hndlLog.Err(err).Msgf("CREATE USER FAILED %+v", &profile)

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	profileData, err := o.CreateProfile(&profile)
	if err != nil {
		hndlLog.Err(err).Msgf("CREATE ORDER FAILED %+v", &profile)

		return ec.JSON(
			http.StatusConflict,
			httpsrv.CreateFailed(err),
		)
	}

	return ec.JSON(
		http.StatusOK,
		ProfileDataResult{Body: profileData},
	)
}

func (o *ProfileV1) ProfileGetHandler(ec echo.Context) (err error) {
	// Main code of handler
	hndlLog := logger.HandlerLogger(&o.log, ec)

	profileID, err := strconv.ParseInt(ec.Param("id"), 10, 64)
	if err != nil {
		hndlLog.Err(err).Msgf("BAD REQUEST, id %s", ec.Param("id"))

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	profileData, err := o.GetProfileByID(profileID)
	if err != nil {
		hndlLog.Err(err).Msgf("NOT FOUND, id %d", profileID)

		return ec.JSON(
			http.StatusNotFound,
			httpsrv.NotFound(err),
		)
	}

	return ec.JSON(
		http.StatusOK,
		ProfileDataResult{Body: profileData},
	)
}

func (o ProfileV1) ProfileSearchGetHandler(ec echo.Context) (err error) {
	// Main code of handler
	hndlLog := logger.HandlerLogger(&o.log, ec)

	var value models.Search

	err = ec.Bind(&value)
	if err != nil {
		hndlLog.Err(err).Msg("BAD REQUEST")

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	profileData, err := o.SearchProfileByLastName(&value)
	if err != nil {
		hndlLog.Err(err).Msgf("failed to get profileData %+v", &value)

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	return ec.JSON(
		http.StatusOK,
		ProfileDataResult{Body: profileData},
	)
}

func (o *ProfileV1) ProfilePutHandler(ec echo.Context) (err error) {
	// Main code of handler
	hndlLog := logger.HandlerLogger(&o.log, ec)

	profileID, err := strconv.ParseInt(ec.Param("id"), 10, 64)
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
			hndlLog.Err(err).Msgf("ORDER DATA NOT UPDATED, id %d", profileID)

			return ec.JSON(
				http.StatusBadRequest,
				httpsrv.BadRequest(err),
			)
		}
	}

	profileData, err := o.UpdateProfileByID(profileID, &bodyBytes)
	if err != nil {
		hndlLog.Err(err).Msgf("BAD REQUEST, id %d, body %s", profileID, string(bodyBytes))

		return ec.JSON(
			http.StatusConflict,
			httpsrv.NotUpdated(err),
		)
	}

	return ec.JSON(
		http.StatusOK,
		ProfileDataResult{Body: profileData},
	)
}

func (o *ProfileV1) ProfileDeleteHandler(ec echo.Context) (err error) {
	// Main code of handler
	hndlLog := logger.HandlerLogger(&o.log, ec)

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
		err = o.HardDeleteProfileByID(userID)
	} else {
		err = o.SoftDeleteProfileByID(userID)
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
