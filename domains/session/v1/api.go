package sessionv1

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/httpsrv"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/logger"
)

func (s *SessionV1) SessionDeleteHandler(ec echo.Context) (err error) {
	// Main code of handler
	hndlLog := logger.HandlerLogger(&s.log, ec)

	idCookie, err := ec.Cookie("rosseti-session")
	if err != nil {
		hndlLog.Err(err).Msg("DELETE SESSION FAILED")

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	err = s.DeleteSession(idCookie.Value)
	if err != nil {
		hndlLog.Err(err).Msgf("DELETE SESSION FAILED %s", idCookie.Value)

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
