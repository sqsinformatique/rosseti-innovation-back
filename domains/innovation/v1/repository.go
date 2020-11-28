package innovationv1

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"strconv"
	"strings"
	"time"

	jsonpatch "github.com/evanphx/json-patch"
	"github.com/sqsinformatique/rosseti-innovation-back/internal/db"
	"github.com/sqsinformatique/rosseti-innovation-back/models"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
)

const (
	innovationIndex = "innovation"
)

func (inn *InnovationV1) CreateInnovation(request *models.Innovation) (*models.Innovation, error) {

	request.CreateTimestamp()

	result, err := inn.orm.InsertInto("innovation", request)
	if err != nil {
		return nil, err
	}

	return result.(*models.Innovation), nil
}

// SearchResults wraps the Elasticsearch search response.
//
type SearchResults struct {
	Total int    `json:"total"`
	Hits  []*Hit `json:"hits"`
}

// Hit wraps the document returned in search response.
//
type Hit struct {
	models.Innovation
	// URL        string        `json:"url"`
	Sort       []interface{} `json:"sort"`
	Highlights *struct {
		Title       []string `json:"title"`
		Description []string `json:"descriptions"`
		Alt         []string `json:"alt"`
		Transcript  []string `json:"transcript"`
	} `json:"highlights,omitempty"`
}

// Search returns results matching a query, paginated by after.
//
func (inn *InnovationV1) Search(query string, after ...string) (*SearchResults, error) {
	var results SearchResults
	es := *inn.elasticDB
	res, err := es.Search(
		es.Search.WithIndex("innovation"),
		es.Search.WithBody(inn.buildQuery(query, after...)),
	)
	if err != nil {
		return &results, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return &results, err
		}
		return &results, fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
	}

	type envelopeResponse struct {
		Took int
		Hits struct {
			Total struct {
				Value int
			}
			Hits []struct {
				ID         string          `json:"_id"`
				Source     json.RawMessage `json:"_source"`
				Highlights json.RawMessage `json:"highlight"`
				Sort       []interface{}   `json:"sort"`
			}
		}
	}

	var r envelopeResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return &results, err
	}

	results.Total = r.Hits.Total.Value

	if len(r.Hits.Hits) < 1 {
		results.Hits = []*Hit{}
		return &results, nil
	}

	for _, hit := range r.Hits.Hits {
		var h Hit
		id, err := strconv.Atoi(hit.ID)
		if err != nil {
			inn.log.Warn().Msgf("failed convert ID: %s", hit.ID)
		}

		h.ID = id
		h.Sort = hit.Sort
		// h.URL = strings.Join([]string{baseURL, hit.ID, ""}, "/")

		if err := json.Unmarshal(hit.Source, &h); err != nil {
			return &results, err
		}

		if len(hit.Highlights) > 0 {
			if err := json.Unmarshal(hit.Highlights, &h.Highlights); err != nil {
				return &results, err
			}
		}

		results.Hits = append(results.Hits, &h)
	}

	return &results, nil
}

// Search returns results matching a query, paginated by after.
//
func (inn *InnovationV1) SearchTitle(query string, after ...string) (*SearchResults, error) {
	var results SearchResults
	es := *inn.elasticDB
	res, err := es.Search(
		es.Search.WithIndex("innovation"),
		es.Search.WithBody(inn.buildQueryTitle(query, after...)),
	)
	if err != nil {
		return &results, err
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return &results, err
		}
		return &results, fmt.Errorf("[%s] %s: %s", res.Status(), e["error"].(map[string]interface{})["type"], e["error"].(map[string]interface{})["reason"])
	}

	type envelopeResponse struct {
		Took int
		Hits struct {
			Total struct {
				Value int
			}
			Hits []struct {
				ID         string          `json:"_id"`
				Source     json.RawMessage `json:"_source"`
				Highlights json.RawMessage `json:"highlight"`
				Sort       []interface{}   `json:"sort"`
			}
		}
	}

	var r envelopeResponse
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return &results, err
	}

	results.Total = r.Hits.Total.Value

	if len(r.Hits.Hits) < 1 {
		results.Hits = []*Hit{}
		return &results, nil
	}

	for _, hit := range r.Hits.Hits {
		var h Hit
		id, err := strconv.Atoi(hit.ID)
		if err != nil {
			inn.log.Warn().Msgf("failed convert ID: %s", hit.ID)
		}

		h.ID = id
		h.Sort = hit.Sort
		// h.URL = strings.Join([]string{baseURL, hit.ID, ""}, "/")

		if err := json.Unmarshal(hit.Source, &h); err != nil {
			return &results, err
		}

		if len(hit.Highlights) > 0 {
			if err := json.Unmarshal(hit.Highlights, &h.Highlights); err != nil {
				return &results, err
			}
		}

		results.Hits = append(results.Hits, &h)
	}

	return &results, nil
}

func (inn *InnovationV1) buildQuery(query string, after ...string) io.Reader {
	var b strings.Builder

	b.WriteString("{\n")

	if query == "" {
		b.WriteString(searchAll)
	} else {
		b.WriteString(fmt.Sprintf(searchMatch, query))
	}

	if len(after) > 0 && after[0] != "" && after[0] != "null" {
		b.WriteString(",\n")
		b.WriteString(fmt.Sprintf(`	"search_after": %s`, after))
	}

	b.WriteString("\n}")

	// fmt.Printf("%s\n", b.String())
	return strings.NewReader(b.String())
}

func (inn *InnovationV1) buildQueryTitle(query string, after ...string) io.Reader {
	var b strings.Builder

	b.WriteString("{\n")

	if query == "" {
		b.WriteString(searchAll)
	} else {
		b.WriteString(fmt.Sprintf(searchMatchTitle, query))
	}

	if len(after) > 0 && after[0] != "" && after[0] != "null" {
		b.WriteString(",\n")
		b.WriteString(fmt.Sprintf(`	"search_after": %s`, after))
	}

	b.WriteString("\n}")

	// fmt.Printf("%s\n", b.String())
	return strings.NewReader(b.String())
}

const searchAll = `
	"query" : { "match_all" : {} },
	"size" : 25,
	"sort" : { "published" : "desc", "_doc" : "asc" }`

const searchMatch = `
	"query" : {
		"multi_match" : {
			"query" : %q,
			"fields" : ["title^100", "descriptions^50", "alt^10", "transcript"],
			"operator" : "and"
		}
	},
	"highlight" : {
		"fields" : {
			"title" : { "number_of_fragments" : 0 },
			"descriptions" : { "number_of_fragments" : 0 },
			"alt" : { "number_of_fragments" : 0 },
			"transcript" : { "number_of_fragments" : 5, "fragment_size" : 25 }
		}
	},
	"size" : 25,
	"sort" : [ { "_score" : "desc" }, { "_doc" : "asc" } ]`

const searchMatchTitle = `
	"query" : {
		"multi_match" : {
			"query" : %q,
			"fields" : ["title^100"],
			"operator" : "and"
		}
	},
	"highlight" : {
		"fields" : {
			"title" : { "number_of_fragments" : 0 }
		}
	},
	"size" : 25,
	"sort" : [ { "_score" : "desc" }, { "_doc" : "asc" } ]`

func writeToGridFile(fileName string, file multipart.File, gridFile *gridfs.UploadStream) (int, error) {
	reader := bufio.NewReader(file)
	defer func() { file.Close() }()
	// make a buffer to keep chunks that are read
	buf := make([]byte, 1024)
	fileSize := 0
	for {
		// read a chunk
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return 0, errors.New("could not read the input file")
		}
		if n == 0 {
			break
		}
		// write a chunk
		if size, err := gridFile.Write(buf[:n]); err != nil {
			return 0, errors.New("could not write to GridFs for " + fileName)
		} else {
			fileSize += size
		}
	}
	gridFile.Close()
	return fileSize, nil
}

func (inn *InnovationV1) CreateImages(actID string, multipartForm *multipart.Form) error {
	for _, fileHeaders := range multipartForm.File {
		for _, fileHeader := range fileHeaders {
			file, _ := fileHeader.Open()
			mongoconn := *inn.mongodb
			bucket, err := gridfs.NewBucket(
				mongoconn.Database(inn.cfg.Mongo.ImageDB),
			)
			if err != nil {
				return err
			}

			gridFile, err := bucket.OpenUploadStream(
				actID + "_" + fileHeader.Filename, // this is the name of the file which will be saved in the database
			)
			if err != nil {
				return err
			}

			fileSize, err := writeToGridFile(fileHeader.Filename, file, gridFile)
			if err != nil {
				return err
			}

			inn.log.Debug().Msgf("Write file to DB was successful. File size: %d \n", fileSize)
		}
	}

	return nil
}

func (inn *InnovationV1) GetImage(actID, imageID string) (*bytes.Buffer, int64, error) {
	mongoconn := *inn.mongodb
	bucket, err := gridfs.NewBucket(
		mongoconn.Database(inn.cfg.Mongo.ImageDB),
	)
	if err != nil {
		return nil, 0, err
	}

	var buf bytes.Buffer
	dStream, err := bucket.DownloadToStreamByName(actID+"_"+imageID, &buf)
	if err != nil {
		return nil, 0, err
	}

	inn.log.Debug().Msgf("File size to download: %v\n", dStream)
	return &buf, dStream, nil
}

func (inn *InnovationV1) GetInnovationByID(id int64) (data *models.Innovation, err error) {
	data = &models.Innovation{}

	conn := *inn.db
	if inn.db == nil {
		return nil, db.ErrDBConnNotEstablished
	}

	err = conn.Get(data, "select * from production.innovation where id=$1", id)
	if err != nil {
		return nil, err
	}

	inn.log.Debug().Msgf("user %+v", data)

	return
}

func mergeInnovationData(oldData *models.Innovation, patch *[]byte) (newData *models.Innovation, err error) {
	id := oldData.ID

	original, err := json.Marshal(oldData)
	if err != nil {
		return
	}

	merged, err := jsonpatch.MergePatch(original, *patch)
	if err != nil {
		return
	}

	err = json.Unmarshal(merged, &newData)
	if err != nil {
		return
	}

	// Protect ID from changes
	newData.ID = id

	newData.UpdatedAt.Time = time.Now()
	newData.UpdatedAt.Valid = true

	return newData, nil
}

func (inn *InnovationV1) UpdateInnovationByID(id int64, patch *[]byte) (writeData *models.Innovation, err error) {
	data, err := inn.GetInnovationByID(id)
	if err != nil {
		return
	}

	writeData, err = mergeInnovationData(data, patch)
	if err != nil {
		return
	}

	if inn.db == nil {
		return nil, db.ErrDBConnNotEstablished
	}

	_, err = inn.orm.Update("innovation", writeData)
	if err != nil {
		return nil, err
	}

	return writeData, err
}

func (inn *InnovationV1) GetInnovationByUserID(id int64) (data *ArrayOfInnovationData, err error) {
	conn := *inn.db
	if inn.db == nil {
		return nil, db.ErrDBConnNotEstablished
	}

	rows, err := conn.Queryx(conn.Rebind("select * from production.innovation where author_id=$1"), id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data = &ArrayOfInnovationData{}

	for rows.Next() {
		var item models.Innovation

		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		*data = append(*data, item)
	}

	return data, nil
}

func (inn *InnovationV1) SelectAllInnovation() (data *ArrayOfInnovationData, err error) {
	conn := *inn.db
	if inn.db == nil {
		return nil, db.ErrDBConnNotEstablished
	}

	rows, err := conn.Queryx("select * from production.innovation")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data = &ArrayOfInnovationData{}

	for rows.Next() {
		var item models.Innovation

		err = rows.StructScan(&item)
		if err != nil {
			return nil, err
		}

		*data = append(*data, item)
	}

	return data, nil
}

func (inn *InnovationV1) GetExpertByInnovationID(id int64) (data *models.InnovationExperts, err error) {
	data = &models.InnovationExperts{}

	conn := *inn.db
	if inn.db == nil {
		return nil, db.ErrDBConnNotEstablished
	}

	err = conn.Get(data, "select * from production.experts where id=$1", id)
	if err != nil {
		return nil, err
	}

	inn.log.Debug().Msgf("user %+v", data)

	return
}
