package main

import (
	"math"
	"testing"
)

func TestAbs(t *testing.T) {
	tests := []struct {
    name string
		arg, want float64
	}{
		{
      name: "-1",
			arg:  -1.0,
			want: 1.0,
		},
		{
      name: "-5",
			arg:  -5.0,
			want: 5.0,
		},
		{
      name: "+5",
			arg:  5.0,
			want: 5.0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, want := math.Abs(test.arg), test.want
			if got != want {
				t.Errorf("Abs(-1.0) = %0.1f; want %f", got, want)
			}
		})
	}
}
