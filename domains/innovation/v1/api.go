package innovationv1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/sqsinformatique/rosseti-innovation-back/internal/echo-swagger"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/httpsrv"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/logger"
	"github.com/sqsinformatique/rosseti-innovation-back/models"
)

func (inn *InnovationV1) innovationPostHandler(ec echo.Context) (err error) {
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("innovationPostHandler").
			SetSummary("Create Innovation").
			AddInBodyParameter("innovation", "Request for create innovation", &models.Innovation{}, false).
			AddInHeaderParameter("Authorization", "Authorization header", reflect.String, true).
			AddResponse(http.StatusOK, "OK", &InnovationDataResult{Body: &models.Innovation{}})
		return nil
	}

	// Main code of handler
	hndlLog := logger.HandlerLogger(&inn.log, ec)

	var innovation models.Innovation
	err = ec.Bind(&innovation)
	if err != nil {
		hndlLog.Err(err).Msgf("CREATE INNOVATION FAILED %+v", &innovation)

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	innovationData, err := inn.CreateInnovation(&innovation)
	if err != nil {
		hndlLog.Err(err).Msgf("CREATE INNOVATION FAILED %+v", &innovation)

		return ec.JSON(
			http.StatusConflict,
			httpsrv.CreateFailed(err),
		)
	}

	jsonData, err := json.Marshal(innovationData)
	if err != nil {
		hndlLog.Err(err).Msgf("MASHAL INNOVATION FAILED %+v", &innovation)

		return ec.JSON(
			http.StatusConflict,
			httpsrv.CreateFailed(err),
		)
	}

	req := esapi.IndexRequest{
		Index:      "innovation",
		DocumentID: strconv.Itoa(innovationData.ID),
		Body:       bytes.NewReader(jsonData),
		Refresh:    "true",
	}

	// Perform the request with the client.
	res, err := req.Do(context.Background(), *inn.elasticDB)
	if err != nil {
		hndlLog.Err(err).Msgf("POST TO ELASTIC FAILED %+v", &innovation)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		hndlLog.Err(err).Msgf("[%s] Error indexing document ID=%d", res.Status(), innovationData.ID)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			hndlLog.Err(err).Msgf("DECODE ELASTIC RESPONSE FAILED %+v", &innovation)
			return err
		} else {
			// Print the response status and indexed document version.
			hndlLog.Debug().Msgf("[%s] %s; version=%d", res.Status(), r["result"], int(r["_version"].(float64)))
		}
	}

	return ec.JSON(
		http.StatusOK,
		InnovationDataResult{Body: innovationData},
	)
}

func (inn *InnovationV1) searchPostHandler(ec echo.Context) error {
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("innovationPostHandler").
			SetSummary("Create Innovation").
			AddInQueryParameter("q", "query", reflect.String, true).
			AddInQueryParameter("a", "after", reflect.String, false).
			AddInHeaderParameter("Authorization", "Authorization header", reflect.String, true).
			AddResponse(http.StatusOK, "OK", &SearchDataResult{Body: &SearchDataResult{}})
		return nil
	}

	// Main code of handler
	hndlLog := logger.HandlerLogger(&inn.log, ec)

	q := ec.QueryParam("q")
	a := ec.QueryParam("a")

	searchResult, err := inn.Search(q, a)
	if err != nil {
		hndlLog.Err(err).Msgf("SEARCH INNOVATION FAILED: query %s", q)

		return ec.JSON(
			http.StatusConflict,
			httpsrv.CreateFailed(err),
		)
	}

	return ec.JSON(
		http.StatusOK,
		SearchDataResult{Body: searchResult},
	)
}
