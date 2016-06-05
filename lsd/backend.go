package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"sort"
)

type NewBackendFactory func() Backend

type Backend interface {
	Init()
	GetData(table string) Rows
}

type Rows []Object
type Object map[string]interface{}

func toObject(val interface{}) Object {
	switch val := val.(type) {
	case Object:
		return val
	case map[string]interface{}:
		return Object(val)
	default:
		return nil
	}
}

func (o Object) Keys() []string {
	i, keys := 0, make([]string, len(o))
	for c := range o {
		keys[i] = c
		i++
	}
	sort.Strings(keys)
	return keys
}

type BackendBase struct {
	logger *log.Logger
}

func (base *BackendBase) Init() {
	base.logger = log.New()
}

// from answer on http://stackoverflow.com/a/31129967/444060
// by http://stackoverflow.com/users/1125956/connor-peet
func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
