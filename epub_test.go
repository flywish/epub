package epub

import (
	"testing"
)

func TestRead(t *testing.T) {
	type args struct {
		src string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "case-normal",
			args: args{
				src: "test_data/世界文明启示录.epub",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Read(tt.args.src)
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
