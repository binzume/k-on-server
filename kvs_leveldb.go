package main

import (
	// "github.com/gin-gonic/contrib/static"
	"encoding/json"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/util"
	"reflect"
)

type LevelDbKVS struct {
	leveldb *leveldb.DB
}

func NewLevelDbKVS(path string) *LevelDbKVS {
	db, _ := leveldb.OpenFile(path, nil)
	return &LevelDbKVS{db}
}

func (c *LevelDbKVS) Close() error {
	return c.leveldb.Close()
}

func (c *LevelDbKVS) Initialize() error {
	return nil
}

func (c *LevelDbKVS) get(typ string, id string, ent interface{}) (found bool, err error) {
	data, err := c.leveldb.Get([]byte(typ+":"+id), nil)
	if err == leveldb.ErrNotFound {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	json.Unmarshal(data, ent)
	return true, nil
}

func (c *LevelDbKVS) store(typ string, id string, ent interface{}) (created bool, err error) {
	s, _ := json.Marshal(ent)
	err = c.leveldb.Put([]byte(typ+":"+id), s, nil)
	return false, err
}

func (c *LevelDbKVS) del(typ string, id string) (found bool, err error) {
	err = c.leveldb.Delete([]byte(typ+":"+id), nil)
	return true, err
}

func (c *LevelDbKVS) query(typ string, slice interface{}, name, term string, offset, limit int) (list interface{}, err error) {
	tt := reflect.ValueOf(slice).Elem().Type().Elem().Elem()
	iter := c.leveldb.NewIterator(util.BytesPrefix([]byte(typ+":")), nil)
	listv := reflect.ValueOf(slice).Elem()
	if listv.Kind() != reflect.Slice {
		panic("not slice")
	}
	count := 0
	iter.Last()
	for iter.Prev() {
		if limit >= 0 && count >= limit {
			break
		}
		value := iter.Value()
		if name != "" {
			m := map[string]string{}
			_ = json.Unmarshal(value, &m)
			if m[name] != term {
				continue
			}
		}
		count++
		if count <= offset {
			continue
		}
		if limit >= 0 && count > offset+limit {
			break
		}
		v := reflect.New(tt).Elem()
		_ = json.Unmarshal(value, v.Addr().Interface())
		listv = reflect.Append(listv, v.Addr())
	}
	iter.Release()
	err = iter.Error()
	reflect.ValueOf(slice).Elem().Set(listv)
	return listv.Interface(), err
}
