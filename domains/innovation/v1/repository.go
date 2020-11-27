package innovationv1

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/sqsinformatique/rosseti-innovation-back/models"
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
