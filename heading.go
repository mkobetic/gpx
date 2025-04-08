package main

import (
	"fmt"
	"strings"
)

// HeadingRange represents the heading range of a Segment
type HeadingRange struct {
	Min       int // counter-clockwise extreme of the heading range
	Mid       int // mid-point of the heading range
	Max       int // clockwise extreme of the heading range
	Variation int // the magnitude of the range in degrees (max-min)
}

func NewHeadingRange(min, max int) *HeadingRange {
	dist := headingDist(min, max)
	return &HeadingRange{
		Min:       min,
		Mid:       headingAdd(min, dist/2),
		Max:       max,
		Variation: dist,
	}
}

func (hr *HeadingRange) Includes(h int) bool {
	return headingBetween(h, hr.Min, hr.Max)
}

func (hr *HeadingRange) Overlaps(hr2 *HeadingRange) bool {
	return hr.Includes(hr2.Min) || hr.Includes(hr2.Max) || hr2.Includes(hr.Min) || hr2.Includes(hr.Max)
}

func (hr *HeadingRange) Merge(hr2 *HeadingRange) *HeadingRange {
	if !hr.Overlaps(hr2) {
		return nil
	}
	min, max := hr.Min, hr.Max
	if headingBetween(min, hr2.Min, hr2.Max) {
		min = hr2.Min
	}
	if headingBetween(max, hr2.Min, hr2.Max) {
		max = hr2.Max
	}
	return NewHeadingRange(min, max)
}

func (hr *HeadingRange) IsBetween(min, max int) bool {
	return headingBetween(hr.Min, min, max) && headingBetween(hr.Max, min, max)
}

func (hr *HeadingRange) String() string {
	return fmt.Sprintf("<%d,%d>", hr.Min, hr.Max)
}

// Represents discontinuous set of heading ranges sorted in clockwise direction.
// Adding a heading to a set merges it with ranges that overlap
// creating new heading that contains all the overlapping ranges.
type HeadingSet []*HeadingRange

func (hs HeadingSet) Add(h *HeadingRange) HeadingSet {
	if len(hs) == 0 {
		return HeadingSet{h}
	}
	result := HeadingSet{}
	prev := hs[0]
	for _, next := range hs {
		if h == nil {
			result = append(result, next)
			continue
		}
		m := next.Merge(h)
		if m != nil {
			h = m
			continue
		}
		if h.IsBetween(prev.Min, next.Min) {
			result = append(result, h)
			h = nil
		}
		result = append(result, next)
	}
	if h != nil {
		result = append(result, h)
	}
	for len(result) > 1 {
		if m := result[len(result)-1].Merge(result[0]); m != nil {
			result[0] = m
			result = result[:len(result)-1]
		} else {
			break
		}
	}
	return result
}

// Return a set with heading ranges that complement hs.
func (hs HeadingSet) Inverse() HeadingSet {
	result := HeadingSet{NewHeadingRange(hs[len(hs)-1].Max, hs[0].Min)}
	prev := hs[0]
	for _, next := range hs[1:] {
		result = append(result, NewHeadingRange(prev.Max, next.Min))
		prev = next
	}
	return result
}

func (hs HeadingSet) Variation() (width int) {
	for _, hr := range hs {
		width += hr.Variation
	}
	return width
}

func (hs HeadingSet) String() string {
	var out []string
	for _, h := range hs {
		out = append(out, h.String())
	}
	return fmt.Sprintf("[%s](%d)", strings.Join(out, ","), hs.Variation())
}
