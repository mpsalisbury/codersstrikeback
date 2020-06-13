package main

import (
	"strings"
	"testing"
)

func TestPlus(t *testing.T) {
	tests := []struct {
		p    Point
		v    Vector
		want Point
	}{
		{
			p:    Point{1, 2},
			v:    Vector{2, 3},
			want: Point{3, 5},
		},
		{
			p:    Point{-1, 4},
			v:    Vector{4, -1},
			want: Point{3, 3},
		},
	}

	for _, test := range tests {
		t.Run("Simple", func(t *testing.T) {
			got, want := *test.p.plus(&test.v), test.want
			if got != want {
				t.Errorf("%v.plus(%v) = %v; want %v", test.p, test.v, got, want)
			}
		})
	}
}

func TestMinus(t *testing.T) {
	tests := []struct {
		p1   Point
		p2   Point
		want Vector
	}{
		{
			p1:   Point{2, 3},
			p2:   Point{1, 2},
			want: Vector{1, 1},
		},
		{
			p1:   Point{-1, 4},
			p2:   Point{4, -1},
			want: Vector{-5, 5},
		},
	}

	for _, test := range tests {
		t.Run("Simple", func(t *testing.T) {
			got, want := *test.p1.minus(&test.p2), test.want
			if got != want {
				t.Errorf("%v.minus(%v) = %v; want %v", test.p1, test.p2, got, want)
			}
		})
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		p    Point
		want string
	}{
		{
			p:    Point{1, 2},
			want: "1 2",
		},
		{
			p:    Point{92, 74},
			want: "92 74",
		},
	}

	for _, test := range tests {
		t.Run("Simple", func(t *testing.T) {
			got, want := test.p.String(), test.want
			if got != want {
				t.Errorf("%v.String() = %s; want %s", test.p, got, want)
			}
		})
	}
}

func TestReadInput(t *testing.T) {
	tests := []struct {
		input string
		want  Point
	}{
		{
			input: "1 2",
			want:  Point{1, 2},
		},
		{
			input: "92 74",
			want:  Point{92, 74},
		},
	}

	for _, test := range tests {
		t.Run("Simple", func(t *testing.T) {
			var got Point
			got.ReadInput(strings.NewReader(test.input))
			want := test.want
			if got != want {
				t.Errorf("ReadInput(%s) = %v; want %v", test.input, got, want)
			}
		})
	}
}
