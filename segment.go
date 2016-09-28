package main

import (
	"time"

	"github.com/tkrajina/gpxgo/gpx"
)

type Segment struct {
	*gpx.GPXTrackSegment
}

// Point returns i-th point of the segment.
func (s Segment) Point(i int) *gpx.GPXPoint {
	return &s.Points[i]
}


// EachPair iterates over a segment with pairs of subsequent points.
func (s Segment) EachPair(f func(prev, next *gpx.GPXPoint)) {
	prev := s.Point(0)
	for i := 1; i < len(s.Points); i++ {
		next := s.Point(i)
		f(prev, next)
		prev = next
	}
}

type Segments []Segment

// Sort segments by start time
func (s Segments) Len() int      { return len(s) }
func (s Segments) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Segments) Less(i, j int) bool {
	return s[i].TimeBounds().StartTime.Before(s[j].TimeBounds().StartTime)
}

// Dedupe removes subsequent segments with the same time bounds
// and segments that have less than 20 points.
func (s Segments) Dedupe() (t Segments) {
	if len(s) == 0 {
		return
	}
	p := s[0]
	t = append(t, p)
	for _, s := range s[1:] {
		if s.TimeBounds().Equals(p.TimeBounds()) {
			continue
		}
		if s.GetTrackPointsNo() > 20 {
			t = append(t, s)
		}
		p = s
	}
	return
}

// Tracks creates tracks from subsequent segments with time bounds that
// are less than limit time apart.
func (s Segments) Tracks(limit time.Duration) (tracks Tracks) {
	if len(s) == 0 {
		return
	}
	p := s[0]
	t := new(gpx.GPXTrack)
	t.AppendSegment(p.GPXTrackSegment)
	for _, s := range s[1:] {
		if s.TimeBounds().StartTime.Sub(p.TimeBounds().EndTime) > limit {
			tracks = append(tracks, Track{t, nil})
			t = new(gpx.GPXTrack)
		}
		t.AppendSegment(s.GPXTrackSegment)
		p = s
	}
	tracks = append(tracks, Track{t, nil})
	return
}
