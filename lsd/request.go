package main

//go:generate go tool yacc -o parser.go -p lsd parser.y

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io"
)

type RequestHeaders []RequestHeader
type RequestHeader struct {
	name string
	arg  HeaderArgument
}

type HeaderArgument interface{}

type Request struct {
	backend Backend
	logger  *log.Entry

	table   string
	headers RequestHeaders
}

type RequestError struct {
	request  *Request
	response *Response
	code     int
	message  string
}

func (err *RequestError) Format(out io.Writer) {
	var rsphdr string
	if err.response != nil {
		rsphdr = err.response.rsphdr
	} else if err.request != nil {
		for _, h := range err.request.headers {
			if h.name == "ResponseHeader" {
				rsphdr = h.toString()
				break
			}
		}
	}

	if rsphdr == "fixed16" {
		fmt.Fprintf(out, "%3d %11d\n%s\n", err.code, len(err.message), err.message)
	}
}

func (err *RequestError) String() string {
	return fmt.Sprintf("RequestError: %d %q", err.code, err.message)
}

func NewRequest(backend Backend, logger *log.Entry) (req *Request) {
	req = &Request{
		backend: backend,
		logger:  logger,
	}
	return
}

func (req *Request) ProcessRequest(request string) *Response {
	req.logger.Debug("parse request:\n", request)
	req.Parse(request)

	req.logger.Debug("table: ", req.table)
	req.logger.Debug("headers: ", req.headers)

	return req.SetupResponse()
}

func (req *Request) Parse(data string) {
	lexer := &lsdLex{
		Request: req,
		data:    []byte(data),
		log:     req.logger,
	}

	lsdDebug = 1
	lsdErrorVerbose = true

	if lsdParse(lexer) != 0 {
		req.Error(400, "Bad request")
	}
}

func (req *Request) SetupResponse() (rsp *Response) {
	rsp = NewResponse(req)
	rsp.data = req.backend.GetData(req.table)
	if rsp.data == nil {
		req.Error(404, fmt.Sprintf("Invalid GET request, no such table '%s'", req.table))
	}
	rsp.processHeaders(req.headers)
	return rsp
}

func (req *Request) Error(code int, message string) {
	panic(&RequestError{
		request: req,
		code:    code,
		message: message,
	})
}

func (req *Request) HandleErrors() {
	err := recover()
	if err == nil {
		return
	}

	switch err := err.(type) {
	case *RequestError:
		panic(err)
	default:
		req.Error(500, fmt.Sprintf("%v", err))
	}
}

func (rh *RequestHeader) toBool() bool {
	switch arg := rh.arg.(type) {
	case string:
		return arg == "on"
	case bool:
		return arg
	case float64:
		return arg != 0
	default:
		return false
	}
}

func (rh *RequestHeader) toString() string {
	switch arg := rh.arg.(type) {
	case string:
		return arg
	default:
		panic(fmt.Sprintf("header value not string: %T (%v)", arg, arg))
	}
}

func (rh *RequestHeader) toStringSlice() (res []string) {
	switch arg := rh.arg.(type) {
	case []HeaderArgument:
		res = make([]string, len(arg))
		for i, value := range arg {
			switch value := value.(type) {
			case string:
				res[i] = value
			}
		}
	case string:
		res = []string{arg}
	default:
		res = []string{}
	}
	return
}

func (rh *RequestHeader) toInt() int {
	switch arg := rh.arg.(type) {
	case float64:
		return int(arg)
	default:
		panic(fmt.Sprintf("header value not int: %T (%v)", arg, arg))
	}
}
