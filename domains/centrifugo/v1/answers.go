package centrifugov1

import (
	"github.com/sqsinformatique/rosseti-innovation-back/internal/httpsrv"
	"github.com/sqsinformatique/rosseti-innovation-back/models"
)

type ChatDataResult httpsrv.ResultAnsw

type ThemeDataResult httpsrv.ResultAnsw

type DirectionsDataResult httpsrv.ResultAnsw

type ArrayOfDirectionData []models.Direction

type ArrayOfThemesData []models.Theme

type ArrayThemeDataResult httpsrv.ResultAnsw

type ArrayOfDirectionDetailedData []models.DirectionDetailed
