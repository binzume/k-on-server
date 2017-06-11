package main

import (
	"encoding/json"
	"reflect"

	elastic "gopkg.in/olivere/elastic.v2"
)

// ElasticKVS
type ElasticKVS struct {
	IndexName string
	Client    *elastic.Client
}

func NewElasticKVS(index string, urls ...string) *ElasticKVS {
	client, err := elastic.NewClient(elastic.SetURL(urls...))
	if err != nil {
		panic(err)
	}
	return &ElasticKVS{index, client}
}

func (c *ElasticKVS) Close() error {
	return nil
}

func (c *ElasticKVS) Initialize() error {
	type H map[string]interface{}
	notAnalyzed := H{"type": "string", "index": "not_analyzed"}
	indexTempalte := H{
		"template": "*",
		"mappings": H{
			"_default_": H{
				"_source": H{"compress": true},
				"properties": H{
					"id": notAnalyzed,
				},
			},
		},
	}

	exists, err := c.Client.IndexExists(c.IndexName).Do()
	if err != nil {
		return err
	}
	if !exists {
		_, err1 := c.Client.CreateIndex(c.IndexName).BodyJson(indexTempalte).Do()
		return err1
	}
	return nil
}

func (c *ElasticKVS) ClearAll() error {
	_, err := c.Client.DeleteIndex(c.IndexName).Do()
	return err
}

func (c *ElasticKVS) get(typ string, id string, ent interface{}) (found bool, err error) {
	result, err := c.Client.Get().
		Index(c.IndexName).Type(typ).
		Id(id).
		Do()
	if err != nil {
		return false, err
	}
	if result.Found {
		json.Unmarshal(*result.Source, ent)
		return true, nil
	}
	return false, nil
}

func (c *ElasticKVS) store(typ string, id string, ent interface{}) (created bool, err error) {
	result, err := c.Client.Index().
		Index(c.IndexName).Type(typ).
		Id(id).
		BodyJson(ent).
		Do()
	if err != nil {
		return false, err
	}
	return result.Created, err
}

func (c *ElasticKVS) del(typ string, id string) (found bool, err error) {
	result, err := c.Client.Delete().
		Index(c.IndexName).Type(typ).
		Id(id).
		Do()
	if err != nil {
		return false, err
	}
	return result.Found, err
}

func appendSlice(slice interface{}, searchResult *elastic.SearchResult) {
	s := reflect.ValueOf(slice).Elem()
	if s.Kind() != reflect.Slice {
		panic("not slice")
	}
	for _, item := range searchResult.Each(s.Type().Elem()) {
		s = reflect.Append(s, reflect.ValueOf(item))
	}
	reflect.ValueOf(slice).Elem().Set(s)
}

func (c *ElasticKVS) query(typ string, slice interface{}, name, term string, offset, limit int) (result interface{}, err error) {
	var query elastic.Query
	if name != "" {
		query = elastic.NewTermQuery(name, term)
	}
	return c.queryInternal(typ, slice, query, "", offset, limit)
}

func (c *ElasticKVS) queryInternal(typ string, slice interface{}, q elastic.Query, sortField string, offset, limit int) (result *elastic.SearchResult, err error) {
	search := c.Client.Search().
		Index(c.IndexName).Type(typ).
		Query(q).
		Pretty(true)
	if sortField != "" {
		search = search.Sort(sortField, false) // desc
	}
	if offset > 0 || limit > 0 {
		search = search.From(offset).Size(limit)
	}
	searchResult, err := search.Do()
	if err == nil {
		appendSlice(slice, searchResult)
	}
	return searchResult, err
}
