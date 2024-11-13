package main

import (
	"reflect"
	"testing"
)

func Test_parseProblem(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want task
	}{
		{
			"+ 15:17",
			args{"+ 15:17"},
			task{true, 0},
		},
		{
			"+10 15:17",
			args{"+10 15:17"},
			task{true, 10},
		},
		{
			"-10 15:17",
			args{"+ 15:17"},
			task{true, 0},
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseTask(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseProblem() = %v, want %v", got, tt.want)
			}
		})
	}
}
