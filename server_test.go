package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestServer(t *testing.T) {
	kvs, _ := NewKVS("leveldb", "testdata")
	kvs.Initialize()
	defer kvs.Close()

	// gin.SetMode(gin.ReleaseMode)
	ts := httptest.NewServer(initHttpd(kvs))
	defer ts.Close()

	res, err := http.Get(ts.URL)
	if err != nil {
		t.Errorf("get failed. %v (%v)", err, ts.URL)
	}
	if res.StatusCode != 200 {
		t.Errorf("invalid response. %v", res)
	}

	res, err = http.Get(ts.URL + "/stats")
	if err != nil {
		t.Errorf("get failed. %v (%v)", err, ts.URL)
	}
	if res.StatusCode != 200 {
		t.Errorf("invalid response. %v", res)
	}
	if !strings.HasPrefix(res.Header.Get("Content-Type"), "application/json") {
		t.Errorf("type: %v", res.Header.Get("Content-Type"))
	}
}
