package main

import "testing"

// 0     1
// 1    1 1
// 2   1 2 1
// 3  1 3 3 1
// 4 1 4 6 4 1
func Test_paskal(t *testing.T) {
	type args struct {
		i int
		j int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"0 0",
			args{0, 0},
			1,
		},
		{
			"1 0",
			args{1, 0},
			1,
		},
		{
			"1 1",
			args{1, 1},
			1,
		},
		{
			"2 0",
			args{2, 0},
			1,
		},
		{
			"2 1",
			args{2, 1},
			2,
		},
		{
			"2 2",
			args{2, 2},
			1,
		},
		{
			"4 1",
			args{4, 1},
			4,
		},
		{
			"4 2",
			args{4, 2},
			6,
		},

		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := paskal(tt.args.i, tt.args.j); got != tt.want {
				t.Errorf("paskal() = %v, want %v", got, tt.want)
			}
		})
	}
}
