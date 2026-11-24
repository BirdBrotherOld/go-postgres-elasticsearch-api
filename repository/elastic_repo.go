package repository

import (
	"belajar-go-pg-es/models"
	"bytes"
	"context"
	"encoding/json"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
)

type ElasticRepo struct {
	es *elasticsearch.Client
}

func NewElasticRepo(es *elasticsearch.Client) *ElasticRepo {
	return &ElasticRepo{es: es}
}

// Upsert dipakai untuk Create maupun Update di Elasticsearch
func (r *ElasticRepo) Upsert(p models.Product) error {
	dataJSON, _ := json.Marshal(p)
	res, err := r.es.Index(
		"products",
		bytes.NewReader(dataJSON),
		r.es.Index.WithDocumentID(strconv.Itoa(p.ID)),
		r.es.Index.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (r *ElasticRepo) Delete(id int) error {
	res, err := r.es.Delete(
		"products",
		strconv.Itoa(id),
		r.es.Delete.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (r *ElasticRepo) Search(ctx context.Context, keyword string, minPrice int, maxPrice int) ([]interface{}, error) {
	// Siapkan kondisi wajib (kata kunci)
	mustQueries := []map[string]interface{}{
		{
			"multi_match": map[string]interface{}{
				"query":     keyword,
				"fields":    []string{"name", "description"},
				"fuzziness": "AUTO",
			},
		},
	}

	// Siapkan kondisi filter rentang harga
	var filterQueries []map[string]interface{}
	priceRange := map[string]interface{}{}
	if minPrice > 0 {
		priceRange["gte"] = minPrice
	}
	if maxPrice > 0 {
		priceRange["lte"] = maxPrice
	}
	if len(priceRange) > 0 {
		filterQueries = append(filterQueries, map[string]interface{}{
			"range": map[string]interface{}{
				"price": priceRange,
			},
		})
	}

	// Gabungkan ke dalam Bool Query Elasticsearch
	esQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   mustQueries,
				"filter": filterQueries,
			},
		},
	}

	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(esQuery)

	res, err := r.es.Search(
		r.es.Search.WithContext(ctx),
		r.es.Search.WithIndex("products"),
		r.es.Search.WithBody(&buf),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var rawResult map[string]interface{}
	json.NewDecoder(res.Body).Decode(&rawResult)

	var products []interface{}
	if hitsObj, ok := rawResult["hits"].(map[string]interface{}); ok {
		if hitsList, ok := hitsObj["hits"].([]interface{}); ok {
			for _, hit := range hitsList {
				hitMap := hit.(map[string]interface{})
				products = append(products, hitMap["_source"])
			}
		}
	}
	return products, nil
}