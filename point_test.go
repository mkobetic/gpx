package main

import (
	"fmt"
	"testing"
)

func Test_HeadingDiff(t *testing.T) {
	for i, tt := range []struct {
		a, b, diff int
	}{
		{20, 60, 40},
		{60, 20, -40},
		{20, 340, -40},
		{20, 201, -179},
		{20, 200, 180},
		{200, 20, 180},
		{20, 199, 179},
		{340, 20, 40},
		{340, 300, -40},
		{340, 161, -179},
		{340, 160, 180},
		{160, 340, 180},
		{340, 159, 179},
	} {
		t.Run(fmt.Sprintf("%d: %d - %d", i, tt.a, tt.b), func(t *testing.T) {
			got := headingDiff(tt.a, tt.b)
			if tt.diff != got {
				t.Errorf("exp %d got %d", tt.diff, got)
			}
		})
	}
}

func Test_HeadingDist(t *testing.T) {
	for i, tt := range []struct {
		a, b, dist int
	}{
		{20, 20, 0},
		{21, 20, 359},
		{20, 60, 40},
		{60, 20, 320},
		{20, 340, 320},
		{20, 201, 181},
		{20, 200, 180},
		{200, 20, 180},
		{20, 199, 179},
		{340, 20, 40},
		{340, 300, 320},
		{340, 161, 181},
		{340, 160, 180},
		{160, 340, 180},
		{340, 159, 179},
	} {
		t.Run(fmt.Sprintf("%d: %d - %d", i, tt.a, tt.b), func(t *testing.T) {
			got := headingDist(tt.a, tt.b)
			if tt.dist != got {
				t.Errorf("exp %d got %d", tt.dist, got)
			}
		})
	}
}

func Test_HeadingAdd(t *testing.T) {
	for i, tt := range []struct {
		h, diff, result int
	}{
		{20, 60, 80},
		{60, -20, 40},
		{20, -60, 320},
		{300, 100, 40},
		{60, 180, 240},
		{240, -179, 61},
	} {
		t.Run(fmt.Sprintf("%d: %d - %d", i, tt.h, tt.diff), func(t *testing.T) {
			got := headingAdd(tt.h, tt.diff)
			if tt.result != got {
				t.Errorf("exp %d got %d", tt.result, got)
			}
		})
	}
}

func Test_HeadingBetween(t *testing.T) {
	for i, tt := range []struct {
		h, min, max int
		isIn        bool
	}{
		{20, 60, 80, false},
		{20, 80, 60, true},
		{81, 60, 80, false},
		{59, 60, 80, false},
		{70, 60, 80, true},
		{60, 60, 80, true},
		{80, 60, 80, true},
	} {
		t.Run(fmt.Sprintf("%d: %d between %d - %d", i, tt.h, tt.min, tt.max), func(t *testing.T) {
			got := headingBetween(tt.h, tt.min, tt.max)
			if tt.isIn != got {
				t.Errorf("exp %t got %t", tt.isIn, got)
			}
		})
	}
}
