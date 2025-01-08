package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bradfitz/latlong"
	"github.com/tkrajina/gpxgo/gpx"
)

type Segment struct {
	*gpx.GPXTrackSegment
	filename string
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

func (s Segment) Timezone() *time.Location {
	b := s.Bounds()
	var err error
	tz, err := time.LoadLocation(latlong.LookupZoneName(b.MinLatitude, b.MinLongitude))
	if err != nil {
		tz = time.UTC
	}
	return tz
}

func (s Segment) String() string {
	tb := s.TimeBounds()
	return fmt.Sprintf("%s @ %s = %05.2fnm (%d)",
		tb.EndTime.Sub(tb.StartTime),
		tb.StartTime.In(s.Timezone()).Format(strFormat),
		s.Length2D()/1852,
		s.GetTrackPointsNo(),
	)
}

func (s Segment) Split(limit time.Duration) Segments {
	limitSeconds := limit.Seconds()
	prev := s.Point(len(s.Points) - 1)
	for i := len(s.Points) - 2; i >= 0; i-- {
		next := s.Point(i)
		if next.TimeDiff(prev) > limitSeconds {
			s1, s2 := s.GPXTrackSegment.Split(i)
			return append(Segment{s1, s.filename}.Split(limit), Segment{s2, s.filename})
		}
		prev = next
	}
	return Segments{s}
}

type Segments []Segment

func GetSegments(g *gpx.GPX, filename string) (s Segments) {
	for _, t := range g.Tracks {
		for i := range t.Segments {
			s = append(s, Segment{&t.Segments[i], filename})
		}
	}
	return s
}

func (ss Segments) String() string {
	var all []string
	for _, s := range ss {
		all = append(all, s.String())
	}
	return strings.Join(all, "\n")
}

// Sort segments by start time
func (s Segments) Len() int      { return len(s) }
func (s Segments) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Segments) Less(i, j int) bool {
	return s[i].TimeBounds().StartTime.Before(s[j].TimeBounds().StartTime)
}

// Dedupe removes subsequent segments with the same time bounds
// and segments that have less than @min points.
func (s Segments) Dedupe(min int) (t Segments) {
	if len(s) == 0 {
		return
	}
	p := s[0]
	t = append(t, p)
	for _, s := range s[1:] {
		if s.TimeBounds().Equals(p.TimeBounds()) {
			continue
		}
		if s.GetTrackPointsNo() > min {
			t = append(t, s)
		}
		p = s
	}
	return
}

func (s Segments) Split(limit time.Duration) (t Segments) {
	for _, seg := range s {
		t = append(t, seg.Split(limit)...)
	}
	return t
}

// Tracks creates tracks from subsequent segments with time bounds that
// are less than limit time apart.
func (ss Segments) Tracks(limit time.Duration) (tracks Tracks) {
	if len(ss) == 0 {
		return
	}
	p := ss[0]
	t := new(gpx.GPXTrack)
	t.AppendSegment(p.GPXTrackSegment)
	for _, s := range ss[1:] {
		if s.TimeBounds().StartTime.Sub(p.TimeBounds().EndTime) > limit {
			tracks = append(tracks, Track{t, nil, s.filename})
			t = new(gpx.GPXTrack)
		}
		t.AppendSegment(s.GPXTrackSegment)
		p = s
	}
	tracks = append(tracks, Track{t, nil, p.filename})
	return
}
