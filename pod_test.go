package main

import (
    "math"
    "testing"
)

func TestAbs(t *testing.T) {
    got := math.Abs(-1.0)
    want := 1.0
    if got != want {
        t.Errorf("Abs(-1.0) = %0.1f; want %f", got, want)
    }
}
