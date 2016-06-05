package main

import (
	log "github.com/Sirupsen/logrus"
	"reflect"
	"testing"
)

func TestSimpleRequest(t *testing.T) {
	logger := log.WithField("test", true)
	req := NewRequest(&TestBackend{
		data: map[string]Rows{
			"hosts": Rows{},
		},
	}, logger)
	rsp := req.ProcessRequest("GET hosts\nColumns: foo bar\n\n")
	if req.table != "hosts" {
		t.Errorf("got table %q, want %q", req.table, "hosts")
	}
	want := []string{"foo", "bar"}
	if !reflect.DeepEqual(rsp.cols, want) {
		t.Errorf("got cols %q, want %q", rsp.cols, want)
	}
}
