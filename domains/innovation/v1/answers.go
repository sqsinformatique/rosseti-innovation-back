package innovationv1

import (
	"github.com/sqsinformatique/rosseti-innovation-back/internal/httpsrv"
	"github.com/sqsinformatique/rosseti-innovation-back/models"
)

type InnovationDataResult httpsrv.ResultAnsw

type SearchDataResult httpsrv.ResultAnsw

type InnovationDataArrayResult httpsrv.ResultAnsw

type ArrayOfInnovationData []models.Innovation

type ArrayOfInnovationDetailData []models.InnovationDetail
