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
