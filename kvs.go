package main

type KVS interface {
	Close() error
	Initialize() error
	get(typ string, id string, ent interface{}) (found bool, err error)
	store(typ string, id string, ent interface{}) (created bool, err error)
	del(typ string, id string) (found bool, err error)
	query(typ string, slice interface{}, name, term string, offset, limit int) (list interface{}, err error)
}
