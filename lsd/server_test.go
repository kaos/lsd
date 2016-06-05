package main

import (
	"strings"
	"testing"
)

func TestResponseHeaders(t *testing.T) {
	backend := &TestBackend{
		data: map[string]Rows{
			"hosts": Rows{
				Object{
					"name":    "foo",
					"address": "1.2.3.4",
				},
				Object{
					"name":    "bar",
					"address": "1.2.3.5",
				},
			},
		},
	}

	tests := map[string]string{
		"GET hosts\nColumns: name\nResponseHeader: fixed16":              "200           8\nfoo\nbar\n",
		"GET otherr\nResponseHeader: fixed16":                            "404          43\nInvalid GET request, no such table 'otherr'\n",
		"GET hosts\nColumns: name barf address\nResponseHeader: fixed16": "400          34\nTable 'hosts' has no column 'barf'\n",
		"GET hosts\nResponseHeader: fixed16\nFoo: bar":                   "400          30\nUndefined request header 'Foo'\n",
	}

	RunRequestTest(t, tests, backend)
}

func TestServeClient(t *testing.T) {
	backend := &TestBackend{
		data: map[string]Rows{
			"hosts": Rows{
				Object{
					"name":    "foo",
					"address": "1.2.3.4",
				},
				Object{
					"name":    "bar.baz",
					"address": "1.2.3.5",
				},
			},
		},
	}

	tests := map[string]string{
		"GET hosts":                                            "address;name\n1.2.3.4;foo\n1.2.3.5;bar.baz\n",
		"GET hosts\n\n":                                        "address;name\n1.2.3.4;foo\n1.2.3.5;bar.baz\n",
		"GET hosts\nOutputFormat: json":                        "[[\"address\",\"name\"],[\"1.2.3.4\",\"foo\"],[\"1.2.3.5\",\"bar.baz\"]]",
		"GET hosts\nColumns: name":                             "foo\nbar.baz\n",
		"GET hosts\nKeepAlive: on\n\nGET hosts\nColumns: name": "address;name\n1.2.3.4;foo\n1.2.3.5;bar.baz\nfoo\nbar.baz\n",
		"GET hosts\nLimit: 1\nColumns: name":                   "foo\n",
		"GET hosts\nLimit: 10\nColumns: name":                  "foo\nbar.baz\n",
		"GET hosts\nOffset: 1\nColumns: name":                  "bar.baz\n",
		"GET hosts\nOffset: 2\nColumns: name":                  "",
		"GET hosts\nOffset: 20\nColumns: name":                 "",
	}

	RunRequestTest(t, tests, backend)
}

func TestFilters(t *testing.T) {
	backend := &TestBackend{
		data: map[string]Rows{
			"services": Rows{
				Object{
					"host_name":   "foo",
					"description": "bar",
					"state":       0,
				},
				Object{
					"host_name":   "foo",
					"description": "baz",
					"state":       0,
				},
				Object{
					"host_name":   "foo",
					"description": "quux",
					"state":       1,
				},
				Object{
					"host_name":   "foo",
					"description": "corge",
					"state":       2,
				},
				Object{
					"host_name":   "roof",
					"description": "grault",
					"state":       2,
				},
				Object{
					"host_name":   "roof",
					"description": "garply",
					"state":       3,
				},
			},
			"hosts": Rows{
				Object{
					"name":    "foo.baz",
					"address": "192.168.0.2",
					"groups": []interface{}{
						"chef-client",
						"sensu-server",
					},
					"last_check": 1324674972,
				},
				Object{
					"name":    "bar.baz",
					"address": "192.168.0.3",
					"groups": []interface{}{
						"chef-client",
						"webserver",
						"memcached",
					},
					"last_check": 1324674956,
				},
				Object{
					"name":    "quux.baz",
					"address": "192.168.0.4",
					"groups": []interface{}{
						"chef-client",
						"mailserver",
					},
					"last_check": 1324674966,
				},
				Object{
					"name":       "orphan.child",
					"address":    "192.168.0.5",
					"groups":     []interface{}{},
					"last_check": 1324674979,
				},
			},
		},
	}

	filter_tests := map[string]string{
		"GET services\nColumnHeaders: off\nFilter: state = 2":                                                               "corge;foo;2\ngrault;roof;2\n",
		"GET services\nFilter: host_name = roof\nColumnHeaders: off\n":                                                      "grault;roof;2\ngarply;roof;3\n",
		"GET services\nFilter: state = 2\nFilter: host_name = roof\nColumnHeaders: off\n":                                   "grault;roof;2\n",
		"GET services\nFilter: state = 2\nFilter: host_name = roof\nOr: 2\nColumnHeaders: off\n":                            "corge;foo;2\ngrault;roof;2\ngarply;roof;3\n",
		"GET services\nFilter: state = 2\nFilter: host_name = roof\nAnd: 2\nFilter: state = 0\nOr: 2\nColumnHeaders: off\n": "bar;foo;0\nbaz;foo;0\ngrault;roof;2\n",
		"GET services\nFilter: state = 0\nNegate:\nColumns: description\n":                                                  "quux\ncorge\ngrault\ngarply\n",
		"GET hosts\nFilter: name = foo.baz\nFilter: name = bar.baz\nOr: 2\nColumns: address":                                "192.168.0.2\n192.168.0.3\n",
		"GET hosts\nFilter: name != foo.bar.baz\nFilter: name !~ child\nAnd: 2\nColumns: name":                              "foo.baz\nbar.baz\nquux.baz\n",
		"GET hosts\nFilter: name != 3foo.3bar.b3az3.roobo\nFilter: name !~ child\nAnd: 2\nColumns: name":                    "foo.baz\nbar.baz\nquux.baz\n",
		"GET hosts\nFilter: name = 3parda01.abc.example.com\nColumnHeaders: off":                                            "",
		"GET hosts\nFilter: name = Foo.Baz\nColumnHeaders: off":                                                             "",
		"GET hosts\nFilter: name =~ Foo.BAZ\nColumns: name":                                                                 "foo.baz\n",
		"GET hosts\nFilter: name !=~ Foo.BAZ\nFilter: name !~ child\nAnd: 2\nColumns: name":                                 "bar.baz\nquux.baz\n",
		"GET hosts\nFilter: name ~ x\nColumns: name":                                                                        "quux.baz\n",
		"GET hosts\nFilter: name ~ X\nColumns: name":                                                                        "",
		"GET hosts\nFilter: name ~~ X\nColumns: name":                                                                       "quux.baz\n",
		"GET hosts\nFilter: groups >= memcached\nColumns: name":                                                             "bar.baz\n",
		"GET hosts\nFilter: groups = \nColumns: name":                                                                       "orphan.child\n",
		"GET hosts\nFilter: groups =\nColumns: name":                                                                        "orphan.child\n",
	}

	RunRequestTest(t, filter_tests, backend)
}

func TestStats(t *testing.T) {
	backend := &TestBackend{
		data: map[string]Rows{
			"services": Rows{
				Object{
					"host_name":      "foo",
					"description":    "bar",
					"state":          0,
					"execution_time": 5.0,
					"check_command":  "check_testthat",
				},
				Object{
					"host_name":      "foo",
					"description":    "baz",
					"state":          0,
					"execution_time": 15,
					"check_command":  "check_testit",
				},
				Object{
					"host_name":      "foo",
					"description":    "quux",
					"state":          1,
					"execution_time": 4,
					"check_command":  "check_testthat",
				},
				Object{
					"host_name":      "foo",
					"description":    "corge",
					"state":          2,
					"execution_time": 7,
					"check_command":  "check_testit",
				},
				Object{
					"host_name":      "roof",
					"description":    "grault",
					"state":          2,
					"execution_time": 8,
					"check_command":  "check_testit",
				},
				Object{
					"host_name":      "roof",
					"description":    "garply",
					"state":          3,
					"execution_time": 3,
					"check_command":  "check_testit",
				},
			},
		},
	}

	tests := map[string]string{
		"GET services\nStats: state = 2":                                                                                                              "2\n",
		"GET services\nFilter: host_name = roof\nStats: state = 3":                                                                                    "1\n",
		"GET services\nStats: state < 3\nStats: state = 3":                                                                                            "5;1\n",
		"GET services\nStats: state != 9999\nColumns: state":                                                                                          "0;2\n1;1\n2;2\n3;1\n",
		"GET services\nStats: min execution_time":                                                                                                     "3\n",
		"GET services\nStats: max execution_time":                                                                                                     "15\n",
		"GET services\nStats: min execution_time\nColumns: host_name":                                                                                 "foo;4\nroof;3\n",
		"GET services\nStats: max execution_time\nColumns: host_name":                                                                                 "foo;15\nroof;8\n",
		"GET services\nStats: sum execution_time":                                                                                                     "42\n",
		"GET services\nStats: avg execution_time":                                                                                                     "7\n",
		"GET services\nStats: avg execution_time\nColumns: host_name":                                                                                 "foo;7.75\nroof;5.5\n",
		"GET services\nStats: avg execution_time\nColumns: host_name check_command":                                                                   "foo;check_testit;11\nfoo;check_testthat;4.5\nroof;check_testit;5.5\n",
		"GET services\nStats: execution_time < 7.5":                                                                                                   "4\n",
		"GET services\nStats: execution_time >= 7.5":                                                                                                  "2\n",
		"GET services\nStats: execution_time <= 6.6":                                                                                                  "3\n",
		"GET services\nStats: state = 0\nStats: state = 2\nStatsOr: 2\nStats: execution_time < 6\nStats: execution_time > 8\nStatsOr: 2\nStatsAnd: 2": "2\n",
		"GET services\nFilter: state = 0\nFilter: state = 2\nOr: 2\nStats: execution_time < 6\nStats: execution_time > 8\nStatsOr: 2":                 "2\n",
	}

	RunRequestTest(t, tests, backend)
}

func TestFields(t *testing.T) {
	backend := &TestBackend{
		data: map[string]Rows{
			"hosts": Rows{
				Object{
					"name":    "i-334455",
					"address": "192.168.0.2",
					"groups": []string{
						"chef-client",
						"sensu-server",
					},
					"last_check":       1324674972,
					"custom_variables": Object{},
				},
				Object{
					"name":    "i-424242",
					"address": "192.168.0.3",
					"groups": []interface{}{
						"chef-client",
						"webserver",
						"memcached",
					},
					"last_check": 1324674956,
					"custom_variables": Object{
						"foo_var": "baz_val",
					},
				},
			},
		},
	}

	tests := map[string]string{
		"GET hosts\nOutputFormat: csv":  "address;custom_variables;groups;last_check;name\n192.168.0.2;;chef-client,sensu-server;1324674972;i-334455\n192.168.0.3;baz_val;chef-client,webserver,memcached;1324674956;i-424242\n",
		"GET hosts\nOutputFormat: json": "[[\"address\",\"custom_variables\",\"groups\",\"last_check\",\"name\"],[\"192.168.0.2\",[],[\"chef-client\",\"sensu-server\"],1324674972,\"i-334455\"],[\"192.168.0.3\",[\"baz_val\"],[\"chef-client\",\"webserver\",\"memcached\"],1324674956,\"i-424242\"]]",
	}

	RunRequestTest(t, tests, backend)

}

func RunRequestTest(t *testing.T, tests map[string]string, backend Backend) {
	for request, want := range tests {
		client := NewTestClient(request)
		ServeClient(client, backend)

		if client.output.String() != want {
			t.Errorf("got response %q, want %q for request %q", client.output.String(), want, request)
		}

		if !client.closed {
			t.Errorf("client.Close() not called")
		}
	}
}

func TestMapRequests(t *testing.T) {
	map_tests := map[string][]string{
		"GET foo":              []string{"GET foo"},
		"GET foo\n":            []string{"GET foo"},
		"GET foo\n\n":          []string{"GET foo"},
		"GET foo\nHeader: bar": []string{"GET foo\nHeader: bar"},
		"GET foo\n\nGET bar":   []string{"GET foo", "GET bar"},
	}

	for data, want := range map_tests {
		idx := 0
		MapRequests(strings.NewReader(data),
			func(got string) bool {
				if got != want[idx] {
					t.Errorf("got request %d %q, want %q, for data %q", idx, got, want[idx], data)
				}
				idx++
				return true
			})
		if idx != len(want) {
			t.Errorf("got %d requests out of data, want %d, for data %q", idx, len(want), data)
		}
	}
}
