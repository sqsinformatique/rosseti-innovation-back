package innovationv1

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"mime"
	"net/http"
	"path/filepath"
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

func (inn *InnovationV1) innovationPutHandler(ec echo.Context) (err error) {
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("innovationPutHandler").
			SetSummary("Update Innovation").
			AddInBodyParameter("innovation", "Request for update innovation", &models.Innovation{}, false).
			AddInHeaderParameter("Authorization", "Authorization header", reflect.String, true).
			AddResponse(http.StatusOK, "OK", &InnovationDataResult{Body: &models.Innovation{}})
		return nil
	}

	// Main code of handler
	hndlLog := logger.HandlerLogger(&inn.log, ec)

	innovationID, err := strconv.ParseInt(ec.Param("id"), 10, 64)
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
			hndlLog.Err(err).Msgf("ORDER DATA NOT UPDATED, id %d", innovationID)

			return ec.JSON(
				http.StatusBadRequest,
				httpsrv.BadRequest(err),
			)
		}
	}

	innovationData, err := inn.UpdateInnovationByID(innovationID, &bodyBytes)
	if err != nil {
		hndlLog.Err(err).Msgf("BAD REQUEST, id %d, body %s", innovationID, string(bodyBytes))

		return ec.JSON(
			http.StatusConflict,
			httpsrv.NotUpdated(err),
		)
	}

	jsonData, err := json.Marshal(innovationData)
	if err != nil {
		hndlLog.Err(err).Msgf("MASHAL INNOVATION FAILED %d", innovationID)

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
		hndlLog.Err(err).Msgf("POST TO ELASTIC FAILED %d", innovationID)
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		hndlLog.Err(err).Msgf("[%s] Error indexing document ID=%d", res.Status(), innovationData.ID)
	} else {
		// Deserialize the response into a map.
		var r map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
			hndlLog.Err(err).Msgf("DECODE ELASTIC RESPONSE FAILED %d", innovationID)
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
			SetDescription("searchPostHandler").
			SetSummary("Search Innovation").
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

func (inn *InnovationV1) searchTitlePostHandler(ec echo.Context) error {
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("searchTittlePostHandler").
			SetSummary("Create Innovation by Title").
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

	searchResult, err := inn.SearchTitle(q, a)
	if err != nil {
		hndlLog.Err(err).Msgf("SEARCH INNOVATION BY TITLE FAILED: query %s", q)

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

func (inn *InnovationV1) innovationPostImagesHandler(ec echo.Context) (err error) {
	// Main code of handler
	hndlLog := logger.HandlerLogger(&inn.log, ec)

	actID := ec.Param("actid")

	multipartForm, err := ec.MultipartForm()
	if err != nil {
		hndlLog.Err(err).Msgf("failed to read multipartform")
		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	err = inn.CreateImages(actID, multipartForm)
	if err != nil {
		hndlLog.Err(err).Msgf("failed to upload images")
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

func (inn *InnovationV1) innovationGetImageHandler(ec echo.Context) (err error) {
	// Main code of handler
	hndlLog := logger.HandlerLogger(&inn.log, ec)

	innID := ec.Param("innid")
	imageID := ec.Param("id")

	gridFile, size, err := inn.GetImage(innID, imageID)
	if err != nil {
		hndlLog.Err(err).Msgf("failed to download image")
		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	ec.Response().Header().Set("Content-Length", strconv.Itoa(int(size)))
	ec.Response().Header().Set("Content-Disposition", "inline; filename=\""+imageID+"\"")

	return ec.Stream(http.StatusOK, mime.TypeByExtension(filepath.Ext(imageID)), gridFile)
}

func (inn *InnovationV1) innovationGetByUserIDHandler(ec echo.Context) (err error) {
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("GetDirectionsHandler").
			SetSummary("Get directions").
			AddInHeaderParameter("Authorization", "Authorization header", reflect.String, true).
			AddResponse(http.StatusOK, "OK", &InnovationDataArrayResult{Body: &ArrayOfInnovationData{}})
		return nil
	}

	// Main code of handler
	hndlLog := logger.HandlerLogger(&inn.log, ec)

	userID, err := strconv.ParseInt(ec.Param("id"), 10, 64)
	if err != nil {
		hndlLog.Err(err).Msgf("BAD REQUEST, id %s", ec.Param("id"))

		return ec.JSON(
			http.StatusBadRequest,
			httpsrv.BadRequest(err),
		)
	}

	directionsData, err := inn.GetInnovationByUserID(userID)
	if err != nil {
		hndlLog.Err(err).Msg("SELECT ALL DIRECTIONS FAILED")

		return ec.JSON(
			http.StatusConflict,
			httpsrv.CreateFailed(err),
		)
	}

	return ec.JSON(
		http.StatusOK,
		InnovationDataArrayResult{Body: directionsData},
	)
}

func (inn *InnovationV1) innovationGetAllDetailedHandler(ec echo.Context) (err error) {
	if echoSwagger.IsBuildingSwagger(ec) {
		echoSwagger.AddToSwagger(ec).
			SetProduces("application/json").
			SetDescription("innovationGetAllDetailedHandler").
			SetSummary("Get all innovation detail").
			AddInHeaderParameter("Authorization", "Authorization header", reflect.String, true).
			AddResponse(http.StatusOK, "OK", &InnovationDataArrayResult{Body: &ArrayOfInnovationData{}})
		return nil
	}

	// Main code of handler
	hndlLog := logger.HandlerLogger(&inn.log, ec)

	innovationData, err := inn.SelectAllInnovation()
	if err != nil {
		hndlLog.Err(err).Msg("SELECT ALL INNOVATION FAILED")

		return ec.JSON(
			http.StatusConflict,
			httpsrv.CreateFailed(err),
		)
	}

	innovationDetailData := ArrayOfInnovationDetailData{}

	for _, v := range *innovationData {
		item := models.InnovationDetail{}
		item.Innovation = v

		author, err := inn.profilev1.GetProfileByID(int64(item.AuthorID))
		if err != nil {
			hndlLog.Err(err).Msgf("SELECT AUTHOR INNOVATION FAILED %+v", v)

			return ec.JSON(
				http.StatusConflict,
				httpsrv.CreateFailed(err),
			)
		}

		item.Author = author

		expertForInnovation, err := inn.GetExpertByInnovationID(int64(item.ID))
		if err != nil {
			hndlLog.Err(err).Msgf("SELECT AUTHOR INNOVATION FAILED %+v", v)

			return ec.JSON(
				http.StatusConflict,
				httpsrv.CreateFailed(err),
			)
		}

		expert, err := inn.profilev1.GetProfileByID(int64(expertForInnovation.ID))
		if err != nil {
			hndlLog.Err(err).Msgf("SELECT EXPERT INNOVATION FAILED %+v", v)

			return ec.JSON(
				http.StatusConflict,
				httpsrv.CreateFailed(err),
			)
		}

		item.Expert = expert

		innovationDetailData = append(innovationDetailData, item)
	}

	return ec.JSON(
		http.StatusOK,
		InnovationDataArrayResult{Body: &innovationDetailData},
	)
}
