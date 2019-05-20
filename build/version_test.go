package build

import "testing"

func TestInfo(t *testing.T) {
	type args struct {
		v string
		m string
		d string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"info",
			args{
				v: "v0.1.0",
				m: "commit",
				d: "20190520",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Info(tt.args.v, tt.args.m, tt.args.d)
			Print()
		})
	}
}
