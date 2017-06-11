package main

import (
	"fmt"
	"strings"
)

// Simple key value store interface
type KVS interface {
	Close() error
	Initialize() error
	get(typ string, id string, ent interface{}) (found bool, err error)
	store(typ string, id string, ent interface{}) (created bool, err error)
	del(typ string, id string) (found bool, err error)
	query(typ string, slice interface{}, name, term string, offset, limit int) (list interface{}, err error)
}

// New KVS client instance
func NewKVS(dbtype, path string) (KVS, error) {
	if dbtype == "leveldb" || dbtype == "" {
		return NewLevelDbKVS(path), nil
	} else if dbtype == "elastic" {
		return NewElasticKVS("k-on", strings.Split(path, ",")...), nil
	}
	return nil, fmt.Errorf("unknown dbtype %s", dbtype)
}
