package config

import (
	"reflect"
	"testing"
)

func TestSkipKindMaps(t *testing.T) {
	for _, testCase := range []struct {
		name         string
		csvSkipKinds string
		expect       map[string]bool
	}{
		{
			"nothing to skip",
			"",
			map[string]bool{},
		},
		{
			"a single kind to skip",
			"somekind",
			map[string]bool{
				"somekind": true,
			},
		},
		{
			"multiple kinds to skip",
			"somekind,anotherkind,yetsomeotherkind",
			map[string]bool{
				"somekind":         true,
				"anotherkind":      true,
				"yetsomeotherkind": true,
			},
		},
	} {
		got := splitCSV(testCase.csvSkipKinds)
		if !reflect.DeepEqual(got, testCase.expect) {
			t.Errorf("%s - got %+v, expected %+v", testCase.name, got, testCase.expect)
		}
	}
}

func TestFromFlags(t *testing.T) {
	testCases := []struct {
		args []string
		conf Config
	}{
		{
			[]string{},
			Config{
				Files:             []string{},
				KubernetesVersion: "1.18.0",
				NumberOfWorkers:   4,
				OutputFormat:      "text",
				SchemaLocations:   nil,
				SkipKinds:         map[string]bool{},
				RejectKinds:       map[string]bool{},
			},
		},
		{
			[]string{"-h"},
			Config{
				Files:             []string{},
				Help:              true,
				KubernetesVersion: "1.18.0",
				NumberOfWorkers:   4,
				OutputFormat:      "text",
				SchemaLocations:   nil,
				SkipKinds:         map[string]bool{},
				RejectKinds:       map[string]bool{},
			},
		},
		{
			[]string{"-skip", "a,b,c"},
			Config{
				Files:             []string{},
				KubernetesVersion: "1.18.0",
				NumberOfWorkers:   4,
				OutputFormat:      "text",
				SchemaLocations:   nil,
				SkipKinds:         map[string]bool{"a": true, "b": true, "c": true},
				RejectKinds:       map[string]bool{},
			},
		},
		{
			[]string{"-summary", "-verbose", "file1", "file2"},
			Config{
				Files:             []string{"file1", "file2"},
				KubernetesVersion: "1.18.0",
				NumberOfWorkers:   4,
				OutputFormat:      "text",
				SchemaLocations:   nil,
				SkipKinds:         map[string]bool{},
				RejectKinds:       map[string]bool{},
				Summary:           true,
				Verbose:           true,
			},
		},
		{
			[]string{"-ignore-missing-schemas", "-kubernetes-version", "1.16.0", "-n", "2", "-output", "json",
				"-schema-location", "folder", "-schema-location", "anotherfolder", "-skip", "kinda,kindb", "-strict",
				"-reject", "kindc,kindd", "-summary", "-verbose", "file1", "file2"},
			Config{
				Files:                []string{"file1", "file2"},
				IgnoreMissingSchemas: true,
				KubernetesVersion:    "1.16.0",
				NumberOfWorkers:      2,
				OutputFormat:         "json",
				SchemaLocations:      []string{"folder", "anotherfolder"},
				SkipKinds:            map[string]bool{"kinda": true, "kindb": true},
				RejectKinds:          map[string]bool{"kindc": true, "kindd": true},
				Strict:               true,
				Summary:              true,
				Verbose:              true,
			},
		},
	}

	for i, testCase := range testCases {
		cfg, _, _ := FromFlags("kubeconform", testCase.args)
		if reflect.DeepEqual(cfg, testCase.conf) != true {
			t.Errorf("test %d: failed parsing config - expected , got: \n%+v\n%+v", i, testCase.conf, cfg)
		}
	}
}
