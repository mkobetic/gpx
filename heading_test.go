package main

import (
	"fmt"
	"testing"
)

func Test_NewHeadingRange(t *testing.T) {
	for i, tt := range []struct {
		min, max, mid, width int
	}{
		{42, 121, 81, 79},
	} {
		t.Run(fmt.Sprintf("%d: <%d, %d>", i, tt.min, tt.max), func(t *testing.T) {
			hr := NewHeadingRange(tt.min, tt.max)
			if tt.mid != hr.Mid {
				t.Errorf("exp %d got %d", tt.mid, hr.Mid)
			}
			if tt.width != hr.Variation {
				t.Errorf("exp %d got %d", tt.width, hr.Variation)
			}
		})
	}
}

func Test_HeadingRangeOverlaps(t *testing.T) {
	for i, tt := range []struct {
		min1, max1 int
		min2, max2 int
		overlaps   bool
	}{
		{20, 60, 80, 200, false},
		{80, 200, 20, 60, false},
		{20, 80, 60, 200, true},
		{60, 200, 20, 80, true},
		{340, 80, 300, 30, true},
		{340, 80, 30, 300, true},
		{340, 30, 80, 300, false},
		{20, 200, 60, 80, true},
		{200, 20, 60, 80, false},
		{20, 200, 80, 60, true},
		{60, 80, 20, 200, true},
		{20, 40, 40, 60, true},
	} {
		t.Run(fmt.Sprintf("%d: <%d, %d> overlaps <%d, %d>", i, tt.min1, tt.max1, tt.min2, tt.max2), func(t *testing.T) {
			hr1 := NewHeadingRange(tt.min1, tt.max1)
			hr2 := NewHeadingRange(tt.min2, tt.max2)
			got := hr1.Overlaps(hr2)
			if tt.overlaps != got {
				t.Errorf("exp %t got %t", tt.overlaps, got)
			}
		})
	}
}

func Test_HeadingRangeMerge(t *testing.T) {
	for i, tt := range []struct {
		min1, max1 int
		min2, max2 int
		min3, max3 int
	}{
		{20, 80, 60, 200, 20, 200},
		{60, 200, 20, 80, 20, 200},
		{340, 80, 300, 30, 300, 80},
		{340, 80, 30, 300, 340, 300},
		{20, 200, 60, 80, 20, 200},
		{20, 200, 80, 60, 80, 60},
		{60, 80, 20, 200, 20, 200},
		{20, 40, 40, 60, 20, 60},
	} {
		t.Run(fmt.Sprintf("%d: <%d,%d> overlaps <%d,%d>", i, tt.min1, tt.max1, tt.min2, tt.max2), func(t *testing.T) {
			hr1 := NewHeadingRange(tt.min1, tt.max1)
			hr2 := NewHeadingRange(tt.min2, tt.max2)
			got := hr1.Merge(hr2)
			if tt.min3 != got.Min {
				t.Errorf("min: exp %d got %d", tt.min3, got.Min)
			}
			if tt.max3 != got.Max {
				t.Errorf("max: exp %d got %d", tt.max3, got.Max)
			}
		})
	}
}

func Test_HeadingSetAdd(t *testing.T) {
	for i, tt := range []struct {
		set1, set2, set3 []int
	}{
		{[]int{}, []int{20, 80, 160, 200}, []int{20, 80, 160, 200}},
		{[]int{}, []int{20, 80, 60, 200}, []int{20, 200}},
		{[]int{50, 150, 200, 300}, []int{250, 100}, []int{200, 150}},
		{[]int{50, 150, 200, 300, 340, 30}, []int{20, 100}, []int{200, 300, 340, 150}},
		{[]int{200, 300, 340, 30}, []int{20, 100}, []int{200, 300, 340, 100}},
		{[]int{279, 32, 121, 214, 219, 250}, []int{208, 235}, []int{279, 32, 121, 250}},
		{[]int{208, 250, 279, 351, 16, 32, 134, 207}, []int{343, 47}, []int{208, 250, 279, 47, 134, 207}},
	} {
		set1 := setFrom(tt.set1)
		set2 := setFrom(tt.set2)
		set3 := setFrom(tt.set3)
		t.Run(fmt.Sprintf("%d: %s add %s", i, set1, set2), func(t *testing.T) {
			for _, hr := range set2 {
				set1 = set1.Add(hr)
			}
			if set1.String() != set3.String() {
				t.Errorf("exp %s got %s", set3, set1)
			}
		})
	}
}

func setFrom(nrs []int) (hs HeadingSet) {
	for i := 0; i < len(nrs); i += 2 {
		hs = append(hs, NewHeadingRange(nrs[i], nrs[i+1]))
	}
	return hs
}
