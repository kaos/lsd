package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
)

type sensuOpts struct {
	api      *string
	mappings *string
}

type sensuBackend struct {
	BackendBase
	opts      sensuOpts
	mappings  map[string]SensuMapping
	templates *template.Template
	tplCtx    Object
}

var sensu sensuBackend

var tplFuncs = template.FuncMap{
	"split":   strings.Split,
	"join":    strings.Join,
	"add":     tplAdd,
	"service": sensu.getService,
}

func init() {
	RegisterBackend("sensu", sensuBackendFactory)
	sensu.opts.api = flag.String("sensu-api", "http://localhost:4567", "Sensu API endpoint URL")
	sensu.opts.mappings = flag.String("sensu-mappings", "sensu-mappings.json", "Map sensu objects to livestatus format")
	sensu.templates = template.New("sensu-mappings")
	sensu.templates.Funcs(tplFuncs)
	sensu.tplCtx = Object{
		"type": "",
		"data": nil,
		"src":  nil,
	}
}

func sensuBackendFactory() (backend Backend) {
	backend = &sensu
	backend.Init()
	return
}

func (sb *sensuBackend) Init() {
	sb.BackendBase.Init()
	sb.loadMappings()
	addInfo("backend", "sensu")
	sb.logger.Infof("initialized backend sensu, api endpoint: %s", *sb.opts.api)
}

func (sb *sensuBackend) GetData(table string) Rows {
	switch table {
	case "hosts":
		return sb.fetchData("/clients", "hosts", true)
	case "services":
		return sb.fetchData("/results", "services", true)
	case "info":
		// addInfo("template", sb.templates.DefinedTemplates())
		return lsdInfo()
	default:
		return nil
	}
}

func (sb *sensuBackend) fetchData(path, typ string, multi bool) Rows {
	var data Rows
	var err error
	if multi {
		err = getJson(*sb.opts.api+path, &data)
	} else {
		var obj Object
		err = getJson(*sb.opts.api+path, &obj)
		data = Rows{obj}
	}
	if err != nil {
		sb.logger.Error(err)
		return nil
	}
	return sb.translate(typ, data)
}

// SensuMapping ...
type SensuMapping map[string]interface{}
type applyMappings []mappingKind
type mappingFun func(dst, src Object, mapping SensuMapping)
type mappingKind struct {
	fun     mappingFun
	mapping SensuMapping
}

func (sb *sensuBackend) translate(typ string, rows Rows) Rows {
	defer func(t string) {
		sb.tplCtx["type"] = t
	}(sb.tplCtx["type"].(string))
	sb.tplCtx["type"] = typ

	mapping := sb.mappings[typ]
	ms := applyMappings{
		mappingKind{
			fun:     sb.mapDefault,
			mapping: getMapping("default", mapping),
		},
		mappingKind{
			fun:     sb.mapObject,
			mapping: getMapping("map", mapping),
		},
	}
	out := make(Rows, len(rows))
	for i, row := range rows {
		out[i] = sb.applyMapping(NewObject(typ), row, ms)
	}
	return out
}

func getMapping(typ string, mapping SensuMapping) SensuMapping {
	res, ok := mapping[typ]
	if ok {
		switch res := res.(type) {
		case SensuMapping:
			return res
		case map[string]interface{}:
			return res
		}
	}
	return nil
}

func (sb *sensuBackend) applyMapping(dst, src Object, mappings applyMappings) Object {
	defer func(s interface{}) {
		sb.tplCtx["src"] = s
	}(sb.tplCtx["src"])
	sb.tplCtx["src"] = src

	for _, m := range mappings {
		if m.mapping != nil {
			m.fun(dst, src, m.mapping)
		}
	}
	return dst
}

func (sb *sensuBackend) mapDefault(dst, src Object, mapping SensuMapping) {
	for dkey, val := range mapping {
		sb.mapValue(dst, dkey, val)
	}
}

func (sb *sensuBackend) mapObject(dst, src Object, mapping SensuMapping) {
	for skey, dkey := range mapping {
		if val, ok := src[skey]; ok {
			sb.mapValue(dst, dkey, val)
		} else {
			t := strings.SplitN(skey, ":", 2)
			if len(t) == 2 {
				sb.tplCtx["value"] = sb.applyTemplate(t[0], t[1], src)
			}
			sb.mapValue(dst, dkey, sb.tplCtx["data"])
			sb.tplCtx["data"] = nil
		}
	}
}

func (sb *sensuBackend) mapValue(dst Object, key, val interface{}) {
	switch key := key.(type) {
	case string:
		dst[key] = val
	case []interface{}:
		for _, akey := range key {
			switch akey := akey.(type) {
			case string:
				dst[akey] = val
			case Object:
				if tpl, ok := akey["template"]; ok {
					sb.mapTemplate(dst, tpl.(string), akey, val)
				}
			case map[string]interface{}:
				if tpl, ok := akey["template"]; ok {
					sb.mapTemplate(dst, tpl.(string), Object(akey), val)
				}
			default:
				panic(fmt.Sprintf("unhandled mapping key type: %T", akey))
			}
		}
	case map[string]interface{}:
		sb.mapObject(dst, toObject(val), key)
	case Object:
		sb.mapObject(dst, toObject(val), SensuMapping(key))
	default:
		panic(fmt.Sprintf("unhandled mapping type: %T", key))
	}

}

func (sb *sensuBackend) mapTemplate(dst Object, tpl string, opts Object, val interface{}) {
	f := opts["field"].(string)
	dst[f] = sb.applyTemplate(f, tpl, val)
}

func (sb *sensuBackend) applyTemplate(name, tpl string, val interface{}) string {
	t := sb.getTemplate(sb.tplCtx["type"].(string)+":"+name, tpl)
	b := &bytes.Buffer{}
	sb.tplCtx["value"] = val
	if err := t.Execute(b, sb.tplCtx); err != nil {
		panic(fmt.Sprintf("%s: render template error: %s", name, err))
	}
	return b.String()
}

func (sb *sensuBackend) getTemplate(name, tpl string) *template.Template {
	if t := sb.templates.Lookup(name); t != nil {
		return t
	}
	return template.Must(sb.templates.New(name).Parse(tpl))
}

func (sb *sensuBackend) loadMappings() {
	file, err := os.Open(*sb.opts.mappings)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&sb.mappings)
	if err != nil {
		panic(err)
	}
}

func (sb *sensuBackend) getService(name string) Object {
	obj := sb.fetchData(
		"/results/"+sb.tplCtx["src"].(Object)["name"].(string)+"/"+name,
		"services", false)[0]
	sb.tplCtx["data"] = obj
	return obj
}

func tplAdd(a, b interface{}) interface{} {
	if a == nil {
		return b
	}
	if b == nil {
		return a
	}
	switch a := a.(type) {
	case float64:
		return a + b.(float64)
	case string:
		return a + b.(string)
	default:
		return fmt.Sprintf("<ADD %#v WITH %#v>", a, b)
	}
}
