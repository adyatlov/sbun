package tools

import (
	"reflect"
	"regexp"
	"testing"
)

func Test_fileNumber(t *testing.T) {
	type args struct {
		fileName string
		r        *regexp.Regexp
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"returns 0 when there is no number and no extension",
			args{"stdout", stdoutRegexp},
			0,
		},
		{
			"returns a correct number when there is a number and no extension",
			args{"stdout.1", stdoutRegexp},
			1,
		},
		{
			"returns a correct number when there is a number > 9, and no extension",
			args{"stdout.999", stdoutRegexp},
			999,
		},
		{
			"returns 0 when there is no number with extension",
			args{"stdout.gz", stdoutRegexp},
			0,
		},
		{
			"returns a correct number when there is a number with extension",
			args{"stdout.1.gz", stdoutRegexp},
			1,
		},
		{
			"returns a correct number when there is a number > 9, with extension",
			args{"stdout.999.gz", stdoutRegexp},
			999,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := fileNumber(tt.args.fileName, tt.args.r); got != tt.want {
				t.Errorf("fileNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterLogFilePaths(t *testing.T) {
	type args struct {
		paths []string
		r     *regexp.Regexp
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			"filters filters unexpected files",
			args{
				[]string{"/path/to/logs/stdout.gz", "stdout", "stdout.2.gz", "stdout.3",
					"unexpected_file", "another_unexpected_file."},
				stdoutRegexp,
			},
			[]string{"/path/to/logs/stdout.gz", "stdout", "stdout.2.gz", "stdout.3"},
		},
		{
			"does not filter expected files",
			args{
				[]string{"/path/to/logs/stdout.gz", "stdout", "stdout.2.gz", "stdout.3"},
				stdoutRegexp,
			},
			[]string{"/path/to/logs/stdout.gz", "stdout", "stdout.2.gz", "stdout.3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterPathsByFileName(tt.args.paths, tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterPathsByFileName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sortFiles(t *testing.T) {
	type args struct {
		paths []string
		r     *regexp.Regexp
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{"sorts files correctly",
			args{
				[]string{
					"path/to/task/logs/stdout.gz",
					"path/to/task/logs/stdout.1",
					"path/to/task/logs/stdout.4.gz",
					"path/to/task/logs/stdout.2.gz",
					"path/to/task/logs/stdout.123",
				},
				stdoutRegexp,
			},
			[]string{
				"path/to/task/logs/stdout.123",
				"path/to/task/logs/stdout.4.gz",
				"path/to/task/logs/stdout.2.gz",
				"path/to/task/logs/stdout.1",
				"path/to/task/logs/stdout.gz",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sortPathsByFileName(tt.args.paths, tt.args.r)
			got := tt.args.paths
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("filterPathsByFileName() sorts %v, want %v", got, tt.want)
			}
		})
	}
}
