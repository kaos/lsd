package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io"
	"strconv"
	"strings"
)

type Response struct {
	request *Request
	logger  *log.Entry

	seps      []string
	cols      []string
	data      Rows
	keepalive bool
	format    string
	rsphdr    string
	offset    int
	limit     int
	colhdrs   interface{}
	filters   []Filterer
	stats     []Filterer
	debug     bool
}

const SEP_DATASET = 0
const SEP_COLUMN = 1
const SEP_LIST = 2
const SEP_JOIN = 3

func NewResponse(req *Request) (rsp *Response) {
	rsp = &Response{
		request: req,
		logger:  req.logger,
		format:  "csv",
		limit:   -1,
		seps:    []string{"\n", ";", ",", "|"},
		filters: []Filterer{},
		stats:   []Filterer{},
	}
	return
}

func (rsp *Response) processHeaders(headers RequestHeaders) {
	rsp.logger.Debug("processing headers...")
	for _, header := range headers {
		switch header.name {
		case "Debug":
			rsp.debug = header.toBool()
		case "Columns":
			rsp.cols = header.toStringSlice()
			if rsp.colhdrs == nil {
				rsp.colhdrs = false
			}
		case "ColumnHeaders":
			rsp.colhdrs = header.toBool()
		case "KeepAlive":
			rsp.keepalive = header.toBool()
		case "OutputFormat":
			rsp.format = header.toString()
		case "ResponseHeader":
			rsp.rsphdr = header.toString()
		case "Limit":
			rsp.limit = header.toInt()
		case "Offset":
			rsp.offset = header.toInt()
		case "Filter":
			rsp.filters = append(rsp.filters, header.arg.(Filterer))
		case "Or":
			rsp.filters = rsp.groupFilters(rsp.filters, &OrFilter{}, header.toInt())
		case "And":
			rsp.filters = rsp.groupFilters(rsp.filters, &AndFilter{}, header.toInt())
		case "Negate":
			rsp.filters = rsp.groupFilters(rsp.filters, &NegateFilter{}, 1)
		case "Stats":
			switch arg := header.arg.(type) {
			case Filterer:
				rsp.stats = append(rsp.stats, &Stats{
					FilterExpression: arg.(*FilterExpression),
					update:           StatsCount,
				})
			case []HeaderArgument:
				// [sum | min | max | avg | ...] field_name
				args := header.toStringSlice()
				if len(args) == 2 {
					f := NewStatser(args[0], args[1])
					if f != nil {
						rsp.stats = append(rsp.stats, f)
					}
				}
			}
		case "StatsOr":
			rsp.stats = rsp.groupFilters(rsp.stats, &OrFilter{}, header.toInt())
		case "StatsAnd":
			rsp.stats = rsp.groupFilters(rsp.stats, &AndFilter{}, header.toInt())
		case "StatsNegate":
			rsp.stats = rsp.groupFilters(rsp.stats, &NegateFilter{}, 1)
		default:
			rsp.Error(400, fmt.Sprintf("Undefined request header '%s'", header.name))
		}
	}

	rsp.data = rsp.processStats(
		rsp.processLimit(
			rsp.processFilters(rsp.data),
		),
	)
	rsp.processColumns()
}

func (rsp *Response) groupFilters(filters []Filterer, grouper FilterGrouper, count int) []Filterer {
	tot := len(filters)
	if count > tot {
		count = tot
	}
	if count <= 0 {
		return filters
	}

	bp := tot - count
	grouper.Group(filters[bp : bp+count])
	return append(filters[0:bp], grouper.(Filterer))
}

func (rsp *Response) processFilters(rows Rows) Rows {
	out := make(Rows, 0, len(rows))
	rsp.logger.Debug("processing filters...")
loop:
	for _, row := range rows {
		for _, f := range rsp.filters {
			if !f.Filter(row) {
				continue loop
			}
		}
		out = append(out, row)
	}
	return out
}

func (rsp *Response) processLimit(rows Rows) Rows {
	if rsp.offset >= len(rows) {
		return nil
	} else if rsp.offset > 0 {
		if rsp.limit < 0 || rsp.offset+rsp.limit >= len(rows) {
			return rows[rsp.offset:]
		} else {
			return rows[rsp.offset:rsp.limit]
		}
	} else if rsp.limit >= 0 && rsp.limit < len(rows) {
		return rows[0:rsp.limit]
	}

	return rows
}

func (rsp *Response) processStats(rows Rows) Rows {
	if len(rsp.stats) == 0 {
		return rows
	}

	rsp.logger.Debug("processing stats...")
	out := Object{}
	for _, row := range rows {
		obj := out
		if rsp.cols != nil {
			for _, col := range rsp.cols {
				k := rsp.csvValue(row[col])
				next, exists := obj[k]
				if !exists {
					next = Object{}
					obj[k] = next
				}
				obj = next.(Object)
			}
		}

		for i, f := range rsp.stats {
			if f.Filter(row) {
				f.(Statser).Update(row, obj, "s"+strconv.FormatInt(int64(i), 10))
			}
		}
	}

	if rsp.debug {
		fmt.Printf("collect stats: %v\nfor cols: %v\n\n", out, rsp.cols)
	}

	if rsp.cols == nil {
		if rsp.colhdrs == nil {
			rsp.colhdrs = false
		}
	}

	rsp.cols = nil
	return rsp.collectStats(out, 0)
}

func (rsp *Response) collectStats(src Object, group int) Rows {
	rows := Rows{}
	var dst Object

	for i, f := range rsp.stats {
		f.(Statser).Update(nil, src, "s"+strconv.FormatInt(int64(i), 10))
	}

	// use src.Keys(), to ensure proper sort order
	keys := src.Keys()
	for _, k := range keys {
		if len(k) > 0 && k[0] == '_' {
			continue
		}
		v := src[k]
		switch v := v.(type) {
		case Object:
			for _, obj := range rsp.collectStats(v, group+1) {
				obj["g"+strconv.FormatInt(int64(group), 10)] = k
				rows = append(rows, obj)
			}
		case int, float64:
			if dst == nil {
				dst = Object{}
				rows = append(rows, dst)
			}
			dst[k] = v
		default:
			panic(fmt.Sprintf("collectStats does not support %T", v))
		}
	}
	return rows
}

func (rsp *Response) processColumns() {
	if rsp.cols == nil {
		var keys []string
		if len(rsp.data) > 0 {
			keys = rsp.data[0].Keys()
		} else {
			keys = NewObject(rsp.request.table).Keys()
		}

		for _, k := range keys {
			if k[0] != '_' {
				rsp.cols = append(rsp.cols, k)
			}
		}
		if rsp.colhdrs == nil {
			rsp.colhdrs = true
		}
	}
}

func (rsp *Response) SendResponse(client io.Writer) {
	if rsp.debug {
		fmt.Printf("cols: %v, data: %v\n", rsp.cols, rsp.data)
	}

	b := &bytes.Buffer{}

	switch rsp.format {
	case "csv":
		rsp.sendCSV(b)
	case "json":
		rsp.sendJSON(b)
	case "python":
		// Not yet implemented
	}

	w := bufio.NewWriter(client)
	switch rsp.rsphdr {
	case "fixed16":
		fmt.Fprintf(w, "200 %11d\n", b.Len())
	}

	rsp.logger.Infof("200 OK %d", b.Len())
	rsp.logger.Debugf("response data:\n%s", b.String())

	b.WriteTo(w)
	w.Flush()
}

func (rsp *Response) csvValue(value interface{}) string {
	return rsp.csvValueSep(value, SEP_LIST)
}

func (rsp *Response) csvValueSep(value interface{}, sep int) string {
	switch value := value.(type) {
	case string:
		return value
	case int:
		return strconv.FormatInt(int64(value), 10)
	case float64:
		return strconv.FormatFloat(value, 'f', -1, 64)
	case []interface{}:
		s := make([]string, len(value))
		for i, v := range value {
			s[i] = rsp.csvValue(v)
		}
		return rsp.csvValueSep(s, sep)
	case []string:
		return strings.Join(value, rsp.seps[sep])
	case Object:
		return rsp.csvValueSep(rsp.marshalObject(value), SEP_JOIN)
	case nil:
		return ""
	default:
		rsp.Error(500, fmt.Sprintf("<csvValue: type '%T' not implemented>", value))
		return ""
	}
}

func (rsp *Response) jsonValue(value interface{}) interface{} {
	switch value := value.(type) {
	case string, int, float64, []string:
		return value
	case []interface{}:
		s := make([]interface{}, len(value))
		for i, v := range value {
			s[i] = rsp.jsonValue(v)
		}
		return s
	case Object:
		return rsp.marshalObject(value)
	case nil:
		return ""
	default:
		rsp.Error(500, fmt.Sprintf("<jsonValue: type '%T' not implemented>", value))
		return nil
	}
}

func (rsp *Response) marshalObject(obj Object) (cols []interface{}) {
	return rsp.marshalObjectKeys(obj, obj.Keys())
}

func (rsp *Response) marshalRootObject(obj Object) (cols []interface{}) {
	return rsp.marshalObjectKeys(obj, rsp.cols)
}

func (rsp *Response) marshalObjectKeys(obj Object, keys []string) (cols []interface{}) {
	cols = make([]interface{}, len(keys))
	for i, col := range keys {
		value, ok := obj[col]
		if !ok {
			rsp.Error(400, fmt.Sprintf("Table '%s' has no column '%s'", rsp.request.table, col))
		}
		switch rsp.format {
		case "csv":
			cols[i] = rsp.csvValue(value)
		case "json":
			cols[i] = rsp.jsonValue(value)
		}
	}
	return
}

func (rsp *Response) sendCSV(out io.Writer) {
	if rsp.colhdrs.(bool) && len(rsp.cols) > 0 {
		fmt.Fprintf(
			out, "%s%s",
			strings.Join(rsp.cols, rsp.seps[SEP_COLUMN]),
			rsp.seps[SEP_DATASET],
		)
	}

	for _, obj := range rsp.data {
		cols := rsp.marshalRootObject(obj)
		rsp.csvPrintRow(out, cols)
	}
}

func (rsp *Response) csvPrintRow(out io.Writer, row []interface{}) {
	for i, col := range row {
		sep := rsp.seps[SEP_COLUMN]
		if i == 0 {
			sep = ""
		}
		fmt.Fprintf(out, "%s%s", sep, col)
	}
	fmt.Fprint(out, rsp.seps[SEP_DATASET])
}

func (rsp *Response) sendJSON(out io.Writer) {
	rows := make([]interface{}, 0)
	if rsp.colhdrs.(bool) && len(rsp.cols) > 0 {
		rows = append(rows, rsp.cols)
	}

	for _, obj := range rsp.data {
		row := rsp.marshalRootObject(obj)
		rows = append(rows, row)
	}

	bytes, err := json.Marshal(rows)
	if err != nil {
		rsp.Error(500, err.Error())
	}
	out.Write(bytes)
}

func (rsp *Response) Keepalive() bool {
	return rsp.keepalive
}

func (rsp *Response) Error(code int, message string) {
	panic(&RequestError{
		response: rsp,
		code:     code,
		message:  message,
	})
}
