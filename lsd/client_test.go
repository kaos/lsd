package main

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"io"
	"strings"
)

type TestClient struct {
	logger *log.Entry
	reader io.Reader
	output *bytes.Buffer
	closed bool
}

func NewTestClient(data string) *TestClient {
	return &TestClient{
		reader: strings.NewReader(data),
		output: new(bytes.Buffer),
		logger: log.WithFields(
			log.Fields{
				"test": true,
			}),
	}
}

func (tc *TestClient) Log() *log.Entry {
	return tc.logger
}

func (tc *TestClient) Read(p []byte) (n int, err error) {
	return tc.reader.Read(p)
}

func (tc *TestClient) Write(p []byte) (n int, err error) {
	return tc.output.Write(p)
}

func (tc *TestClient) Close() {
	tc.closed = true
}

func (tc *TestClient) HandleErrors() {
	err := recover()
	if err == nil {
		return
	}

	switch err := err.(type) {
	case *RequestError:
		err.Format(tc)
		if err.code >= 500 {
			tc.logger.Error(err)
			panic(err)
		}
	default:
		panic(err)
	}
}
