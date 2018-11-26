package main

import (
	"reflect"
	"testing"
)

func Test_counter_included(t *testing.T) {
	type fields struct {
		current uint
		sums    []uint
		len     uint
		capa    uint
	}
	tests := []struct {
		name   string
		fields fields
		want   *counter
	}{
		{
			name: "increment ok",
			fields: fields{
				current: 1,
				sums:    []uint{0},
				len:     1,
				capa:    2,
			},
			want: &counter{
				current: 1,
				sums:    []uint{0, 1},
				len:     2,
				capa:    2,
			},
		},
		{
			name: "increment capa ok",
			fields: fields{
				current: 3,
				sums:    []uint{0, 1},
				len:     2,
				capa:    2,
			},
			want: &counter{
				current: 3,
				sums:    []uint{1, 3},
				len:     2,
				capa:    2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &counter{
				current: tt.fields.current,
				sums:    tt.fields.sums,
				len:     tt.fields.len,
				capa:    tt.fields.capa,
			}
			c.included()

			if !reflect.DeepEqual(c, tt.want) {
				t.Errorf("indluded() = %v, want %v", c, tt.want)
			}
		})
	}
}
