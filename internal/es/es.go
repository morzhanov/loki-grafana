package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/morzhanov/go-elk-example/internal/doc"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
)

type client struct {
	esearch    *elasticsearch.Client
	esDocIndex string
	log        *zap.Logger
}

type DocResponse struct {
	ID     string `json:"id"`
	Author string `json:"author"`
	Title  string `json:"title"`
	Text   string `json:"text"`
}

type UpdateDocRequest struct {
	Doc *doc.Document `json:"doc"`
}

type ElasticSearch interface {
	Save(d *doc.Document) error
	Find(field string, value string) ([]*DocResponse, error)
	Update(id string, d *doc.Document) error
	Delete(id string) error
}

func (c *client) checkResponseError(r *esapi.Response) error {
	if !r.IsError() {
		return nil
	}
	var e map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil {
		return fmt.Errorf("error parsing the response body: %s", err)
	}
	return fmt.Errorf("[%s] %s: %s",
		r.Status(),
		e["error"].(map[string]interface{})["type"],
		e["error"].(map[string]interface{})["reason"],
	)
}

func (c *client) Save(d *doc.Document) error {
	b, err := json.Marshal(d)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(b)

	id := uuid.NewV4().String()
	c.log.Info("saving new doc...", zap.String("id", id), zap.String("title", d.Title))
	r, err := c.esearch.Create(c.esDocIndex, id, reader)
	if err != nil {
		return err
	}
	return c.checkResponseError(r)
}

func (c *client) Find(field string, value string) ([]*DocResponse, error) {
	c.log.Info("looking for docs...", zap.String("field", field), zap.String("value", value))

	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match": map[string]interface{}{
				field: value,
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	r, err := c.esearch.Search(
		c.esearch.Search.WithContext(context.Background()),
		c.esearch.Search.WithIndex(c.esDocIndex),
		c.esearch.Search.WithBody(&buf),
		c.esearch.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	if err := c.checkResponseError(r); err != nil {
		return nil, err
	}
	var dec map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&dec); err != nil {
		return nil, err
	}

	var res []*DocResponse
	for _, mhit := range dec["hits"].(map[string]interface{})["hits"].([]interface{}) {
		hit := mhit.(map[string]interface{})
		data := hit["_source"].(map[string]interface{})
		d := DocResponse{
			ID:     hit["_id"].(string),
			Author: data["author"].(string),
			Title:  data["title"].(string),
			Text:   data["text"].(string),
		}
		res = append(res, &d)
	}
	return res, r.Body.Close()
}

func (c *client) Update(id string, d *doc.Document) error {
	u := UpdateDocRequest{Doc: d}
	b, err := json.Marshal(u)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(b)

	c.log.Info("updating a doc...", zap.String("id", id))
	r, err := c.esearch.Update(c.esDocIndex, id, reader)
	if err != nil {
		return err
	}
	return c.checkResponseError(r)
}

func (c *client) Delete(id string) error {
	c.log.Info("deleting a doc...", zap.String("id", id))
	r, err := c.esearch.Delete(c.esDocIndex, id)
	if err != nil {
		return err
	}
	return c.checkResponseError(r)
}

func NewES(esUri string, esDocIndex string, log *zap.Logger) (ElasticSearch, error) {
	esearch, err := elasticsearch.NewClient(elasticsearch.Config{Addresses: []string{esUri}})
	if err != nil {
		return nil, err
	}
	return &client{esearch, esDocIndex, log}, nil
}
