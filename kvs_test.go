package main

import "testing"

func TestKVS(t *testing.T) {
	type TestData struct {
		Data string
	}
	id := "123"

	// kvs, _ := NewKVS("elastic", "http://localhost:9200,http://localhost:9201")
	kvs, _ := NewKVS("leveldb", "testdata")
	kvs.ClearAll()
	kvs.Initialize()

	// store
	created, err := kvs.store("test", id, &TestData{"hoge"})
	if err != nil {
		t.Errorf("store failed. %v", err)
	}
	if !created {
		t.Error("already exits")
	}

	// get
	var result TestData
	found, err := kvs.get("test", id, &result)
	if err != nil {
		t.Errorf("%v != nil", err)
	}
	if !found {
		t.Error("not found")
	}
	if result.Data != "hoge" {
		t.Errorf("%v != hoge", result)
	}

	// get as map
	resultmap := make(map[string]interface{})
	found, err = kvs.get("test", id, &resultmap)
	if err != nil {
		t.Errorf("get as map failed. %v", err)
	}
	if resultmap["Data"] != "hoge" {
		t.Errorf("%v != hoge", resultmap)
	}

	// query
	// time.Sleep(3 * time.Second) // wait for elasticsearch
	var results []*TestData
	queryResult, err := kvs.query("test", &results, "Data", "hoge", 0, -1)
	if err != nil {
		t.Errorf("query error. %v", err)
	}
	if len(results) != 1 {
		t.Errorf("empty? %v", len(results))
		t.Errorf("queryResult: %v", queryResult)
	}

	// delete
	found, err = kvs.del("test", id)
	if err != nil {
		t.Errorf("%v != nil", err)
	}
	if !found {
		t.Error("not found")
	}

	// get deleted data
	found, err = kvs.get("test", id, &result)
	if found {
		t.Error("not deleted")
	}

	err = kvs.Close()
	if err != nil {
		t.Errorf("Close failed. %v", err)
	}
}
