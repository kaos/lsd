package main

import (
	"reflect"
	"testing"
)

func TestHostMapping(t *testing.T) {
	data := Rows{
		Object{
			"name":    "i-334455",
			"address": "192.168.0.2",
			"subscriptions": []interface{}{
				"chef-client",
				"sensu-server",
			},
			"timestamp": 1324674972,
		},
		Object{
			"name":    "i-424242",
			"address": "192.168.0.3",
			"subscriptions": []interface{}{
				"chef-client",
				"webserver",
				"memcached",
			},
			"timestamp": 1324674956,
		},
	}
	want := Rows{
		Object{
			"name":    "i-334455",
			"address": "192.168.0.2",
			"groups": []interface{}{
				"chef-client",
				"sensu-server",
			},
			"last_check": 1324674972,
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
		},
	}

	mapping := SensuMapping{
		"map": SensuMapping{
			"name":          "name",
			"address":       "address",
			"subscriptions": "groups",
			"timestamp":     "last_check",
		},
	}

	TranslateSensuData(t, data, want, mapping)
}

func TestServiceMapping(t *testing.T) {
	data := Rows{
		Object{
			"client": "i-424242",
			"check": Object{
				"name":    "chef_client_process",
				"command": "check-procs.rb -p chef-client -W 1",
				"subscribers": []interface{}{
					"production",
				},
				"interval": 60,
				"issued":   1389374667,
				"executed": 1389374667,
				"output":   "WARNING Found 0 matching processes\n",
				"status":   1,
				"duration": 0.005,
			},
		},
	}
	want := Rows{
		Object{
			"host_name":     "i-424242",
			"description":   "chef_client_process",
			"check_command": "check-procs.rb -p chef-client -W 1",
			"groups": []interface{}{
				"production",
			},
			"check_interval":     60,
			"last_check":         1389374667,
			"long_plugin_output": "WARNING Found 0 matching processes\n",
			"state":              1,
			"execution_time":     0.005,
		},
	}

	mapping := SensuMapping{
		"map": SensuMapping{
			"client": "host_name",
			"check": map[string]interface{}{
				"command":     "check_command",
				"interval":    "check_interval",
				"name":        "description",
				"duration":    "execution_time",
				"subscribers": "groups",
				"executed":    "last_check",
				"output":      "long_plugin_output",
				"status":      "state",
			},
		},
	}

	TranslateSensuData(t, data, want, mapping)
}

func TestTemplateMapping(t *testing.T) {
	data := Rows{
		Object{
			"simple": "abcd",
			"data":   "foo,bar,baz",
		},
	}

	want := Rows{
		Object{
			"fixed":         "<abcd>",
			"first_data":    "foo",
			"rejoined_data": "foo|bar|baz",
			"no-source":     "value: <nil>",
		},
	}

	mapping := SensuMapping{
		"map": SensuMapping{
			"simple": []interface{}{
				Object{"field": "fixed", "template": "<{{.value}}>"},
			},
			"data": []interface{}{
				Object{
					"field":    "first_data",
					"template": "{{index (split .value `,`) 0}}",
				},
				map[string]interface{}{
					"field":    "rejoined_data",
					"template": "{{join (split .value `,`) `|`}}",
				},
			},
			"_": []interface{}{
				Object{
					"field":    "no-source",
					"template": "value: {{printf `%#v` .value}}",
				},
			},
		},
	}

	TranslateSensuData(t, data, want, mapping)
}

func TranslateSensuData(t *testing.T, sensuData, want Rows, mapping SensuMapping) {
	sensu.mappings = map[string]SensuMapping{"test": mapping}
	got := sensu.translate("test", sensuData)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v using mapping %v", got, want, mapping)
	}
}
