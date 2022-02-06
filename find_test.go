package gofind

import (
	"os"
	"testing"
	"testing/fstest"
)

var fileSystem = fstest.MapFS{
	"test-dir": &fstest.MapFile{
		Mode: os.ModeDir,
	},
	"test_dir/test_file": &fstest.MapFile{},
	"test_dir/test_second_dir": &fstest.MapFile{
		Mode: os.ModeDir,
	},
	"test_dir/test_second_dir/test_file_first": &fstest.MapFile{},
	"test_dir/test_second_dir/test_third_dir": &fstest.MapFile{
		Mode: os.ModeDir,
	},
	"test_dir/test_second_dir/test_third_dir/file_in_third_dir": &fstest.MapFile{},
	"test_dir/test_second_dir/add_me_to_map_too":                &fstest.MapFile{},
	"test_dir/test_second_dir/add_me":                           &fstest.MapFile{},
}

func TestNewFind(t *testing.T) {

	var testData = []struct {
		fs          fstest.MapFS
		expression  string
		matchNumber int
		matches     map[string]bool
	}{
		{
			fs:          fileSystem,
			expression:  "test_dir/test_file",
			matchNumber: 1,
			matches: map[string]bool{
				"test_dir/test_file": true,
			},
		},
		{
			fs:          fileSystem,
			expression:  "test_dir/test_second_dir/test_third_dir/",
			matchNumber: 1,
			matches: map[string]bool{
				"test_dir/test_second_dir/test_third_dir/file_in_third_dir": true,
			},
		},
		{
			fs:          fileSystem,
			expression:  "test_dir/test_second_dir/test_third_dir",
			matchNumber: 2,
			matches: map[string]bool{
				"test_dir/test_second_dir/test_third_dir":                   true,
				"test_dir/test_second_dir/test_third_dir/file_in_third_dir": true,
			},
		},
		{
			fs:          fileSystem,
			expression:  "test_dir",
			matchNumber: 8,
			matches: map[string]bool{
				"test_dir":                                                  true,
				"test_dir/test_file":                                        true,
				"test_dir/test_second_dir":                                  true,
				"test_dir/test_second_dir/test_file_first":                  true,
				"test_dir/test_second_dir/test_third_dir":                   true,
				"test_dir/test_second_dir/test_third_dir/file_in_third_dir": true,
				"test_dir/test_second_dir/add_me_to_map_too":                true,
				"test_dir/test_second_dir/add_me":                           true,
			},
		},
		{
			fs:          fileSystem,
			expression:  "test_dir/test_second_dir/*",
			matchNumber: 5,
			matches: map[string]bool{
				"test_dir/test_second_dir/test_file_first":                  true,
				"test_dir/test_second_dir/test_third_dir":                   true,
				"test_dir/test_second_dir/test_third_dir/file_in_third_dir": true,
				"test_dir/test_second_dir/add_me_to_map_too":                true,
				"test_dir/test_second_dir/add_me":                           true,
			},
		},
		{
			fs:          fileSystem,
			expression:  "test_dir/test_second_dir/test_*",
			matchNumber: 3,
			matches: map[string]bool{
				"test_dir/test_second_dir/test_file_first":                  true,
				"test_dir/test_second_dir/test_third_dir":                   true,
				"test_dir/test_second_dir/test_third_dir/file_in_third_dir": true,
			},
		},
		{
			fs:          fileSystem,
			expression:  "test_dir/test_second_dir/test_third*",
			matchNumber: 2,
			matches: map[string]bool{
				"test_dir/test_second_dir/test_third_dir":                   true,
				"test_dir/test_second_dir/test_third_dir/file_in_third_dir": true,
			},
		},
	}

	for _, value := range testData {
		find := NewFind(value.expression, value.fs)
		find.Run()
		result := find.Results()
		if len(result) != value.matchNumber {
			t.Errorf("there should be exact %d matches, but found %d", value.matchNumber, len(result))
		}
		for _, n := range result {
			if _, ok := value.matches[n]; !ok {
				t.Errorf("the path should be found in the filesystem: %s", n)
			}
		}

	}
}
