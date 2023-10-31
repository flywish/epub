package epub

import (
	"testing"
)

func TestRead(t *testing.T) {
	type args struct {
		src        string
		isMakeFile bool
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "case-normal",
			args: args{
				src:        "test_data/gaodengshuxue.epub",
				isMakeFile: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Read(tt.args.src, tt.args.isMakeFile)
		})
	}
}

func TestWrite(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name string
	}{
		{
			name: "case-normal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Write()
		})
	}
}
