package main

import (
	"math"
	"strings"
	"testing"
)

const TOLERANCE = 0.000000001

func almostEqual(f1, f2 float64) bool {
	return math.Abs(f1-f2) < TOLERANCE
}

func almostEqualV(v1, v2 Vector) bool {
	return almostEqual(v1.vx, v2.vx) && almostEqual(v1.vy, v2.vy)
}

func TestVectorLen(t *testing.T) {
	tests := []struct {
		v    Vector
		want float64
	}{
		{
			v:    Vector{2, 0},
			want: 2,
		},
		{
			v:    Vector{0, -4},
			want: 4,
		},
		{
			v:    Vector{1, 1},
			want: math.Sqrt2,
		},
	}

	for _, test := range tests {
		t.Run("Simple", func(t *testing.T) {
			got, want := test.v.len(), test.want
			if !almostEqual(got, want) {
				t.Errorf("%v.norm() = %0.4f; want %0.4f", test.v, got, want)
			}
		})
	}
}

func TestVectorLen2(t *testing.T) {
	tests := []struct {
		v    Vector
		want float64
	}{
		{
			v:    Vector{2, 0},
			want: 4,
		},
		{
			v:    Vector{0, -4},
			want: 16,
		},
		{
			v:    Vector{1, 1},
			want: 2,
		},
	}

	for _, test := range tests {
		t.Run("Simple", func(t *testing.T) {
			got, want := test.v.len2(), test.want
			if !almostEqual(got, want) {
				t.Errorf("%v.norm() = %0.4f; want %0.4f", test.v, got, want)
			}
		})
	}
}

func TestVectorNorm(t *testing.T) {
	tests := []struct {
		v    Vector
		want Vector
	}{
		{
			v:    Vector{2, 0},
			want: Vector{1, 0},
		},
		{
			v:    Vector{0, -4},
			want: Vector{0, -1},
		},
		{
			v:    Vector{1, 1},
			want: Vector{1.0 / math.Sqrt2, 1.0 / math.Sqrt2},
		},
	}

	for _, test := range tests {
		t.Run("Simple", func(t *testing.T) {
			got, want := *test.v.norm(), test.want
			if !almostEqualV(got, want) {
				t.Errorf("%v.norm() = %v; want %v", test.v, got, want)
			}
		})
	}
}

func TestVectorMinus(t *testing.T) {
	tests := []struct {
		v1   Vector
		v2   Vector
		want Vector
	}{
		{
			v1:   Vector{2, 3},
			v2:   Vector{1, 2},
			want: Vector{1, 1},
		},
		{
			v1:   Vector{-1, 4},
			v2:   Vector{4, -1},
			want: Vector{-5, 5},
		},
	}

	for _, test := range tests {
		t.Run("Simple", func(t *testing.T) {
			got, want := *test.v1.minus(&test.v2), test.want
			if got != want {
				t.Errorf("%v.minus(%v) = %v; want %v", test.v1, test.v2, got, want)
			}
		})
	}
}

func TestVectorDot(t *testing.T) {
	tests := []struct {
		v1   Vector
		v2   Vector
		want float64
	}{
		{
			v1:   Vector{1, 2},
			v2:   Vector{2, 3},
			want: 8,
		},
		{
			v1:   Vector{-1, 4},
			v2:   Vector{4, -1},
			want: -8,
		},
		{
			v1:   Vector{1, 2},
			v2:   Vector{-2, 1},
			want: 0,
		},
	}

	for _, test := range tests {
		t.Run("Simple", func(t *testing.T) {
			got, want := test.v1.dot(&test.v2), test.want
			if got != want {
				t.Errorf("%v.dot(%v) = %0.1f; want %0.1f", test.v1, test.v2, got, want)
			}
		})
	}
}

func TestVectorPerpendicular(t *testing.T) {
	tests := []struct {
		v Vector
	}{
		{
			v: Vector{2, 3},
		},
		{
			v: Vector{-1, 4},
		},
		{
			v: Vector{1, -4},
		},
	}

	for _, test := range tests {
		t.Run("Simple", func(t *testing.T) {
			perp := test.v.perpendicular()
			dot := test.v.dot(perp)
			if dot != 0.0 {
				t.Errorf("%v not perp to %v (dot = %0.1f)", perp, test.v, dot)
			}
		})
	}
}

func TestVectorTimes(t *testing.T) {
	tests := []struct {
		v    Vector
		f    float64
		want Vector
	}{
		{
			v:    Vector{2, 3},
			f:    2.0,
			want: Vector{4, 6},
		},
		{
			v:    Vector{-1, 4},
			f:    4.0,
			want: Vector{-4, 16},
		},
	}

	for _, test := range tests {
		t.Run("Simple", func(t *testing.T) {
			got, want := *test.v.times(test.f), test.want
			if got != want {
				t.Errorf("%v.times(%0.1f) = %v; want %v", test.v, test.f, got, want)
			}
		})
	}
}

func TestVectorString(t *testing.T) {
	tests := []struct {
		v    Vector
		want string
	}{
		{
			v:    Vector{1, 2},
			want: "1 2",
		},
		{
			v:    Vector{92, 74},
			want: "92 74",
		},
	}

	for _, test := range tests {
		t.Run("Simple", func(t *testing.T) {
			got, want := test.v.String(), test.want
			if got != want {
				t.Errorf("%v.String() = %s; want %s", test.v, got, want)
			}
		})
	}
}

func TestVectorReadInput(t *testing.T) {
	tests := []struct {
		input string
		want  Vector
	}{
		{
			input: "1 2",
			want:  Vector{1, 2},
		},
		{
			input: "92 74",
			want:  Vector{92, 74},
		},
	}

	for _, test := range tests {
		t.Run("Simple", func(t *testing.T) {
			var got Vector
			got.ReadInput(strings.NewReader(test.input))
			want := test.want
			if got != want {
				t.Errorf("ReadInput(%s) = %v; want %v", test.input, got, want)
			}
		})
	}
}
