package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bradfitz/latlong"
	"github.com/tkrajina/gpxgo/gpx"
)

// SpeedRange represents the speed range of a Segment
type SpeedRange struct {
	Min float64
	Avg float64
	Max float64
}

func (s *SpeedRange) String(unit unit) string {
	return fmt.Sprintf("%.1f/%.1f/%.1f %s", s.Avg, s.Min, s.Max, unit.speed())
}

// Segment represents a track section of continuous straight movement, turn or stagnation.
// This segment type is represented by its Mode.
// As such a neighboring segments in a track should always have different Mode.
type Segment struct {
	gpx      *gpx.GPXTrackSegment
	filename string
	// Analysis results
	params   *AnalysisParameters
	Points   Points
	previous *Segment
	next     *Segment
	Heading  *HeadingRange
	Distance float64 // length of the segment
	Speed    SpeedRange
	Duration time.Duration
	Start    time.Time
	End      time.Time
	Mode     Mode
	// Activity specific segment type
	Type fmt.Stringer
}

func SegmentFromPoints(ps Points, mode Mode, filename string, params *AnalysisParameters) *Segment {
	gpxPoints := make([]gpx.GPXPoint, len(ps))
	for i, p := range ps {
		gpxPoints[i] = *p.gpx
	}

	return &Segment{
		gpx:      &gpx.GPXTrackSegment{Points: gpxPoints},
		filename: filename,
		params:   params,
		Points:   ps,
		Heading:  ps.headingRange(mode),
		Distance: ps.distance(),
		Speed:    ps.speed(mode),
		Start:    ps.start(),
		End:      ps.end(),
		Duration: ps.duration(),
		Mode:     mode,
	}
}

// gpxPoint returns i-th point of the segment.
func (s *Segment) gpxPoint(i int) *gpx.GPXPoint {
	return &s.gpx.Points[i]
}

// EachPair iterates over a segment with pairs of subsequent points.
func (s *Segment) EachPair(f func(prev, next *Point)) {
	prev := s.Points[0]
	for _, next := range s.Points[1:] {
		f(prev, next)
		prev = next
	}
}

func (s *Segment) Timezone() *time.Location {
	b := s.gpx.Bounds()
	var err error
	tz, err := time.LoadLocation(latlong.LookupZoneName(b.MinLatitude, b.MinLongitude))
	if err != nil {
		tz = time.UTC
	}
	return tz
}

func (s *Segment) gpxString() string {
	tb := s.gpx.TimeBounds()
	return fmt.Sprintf("%s @ %s = %05.2f%s (%d)",
		tb.EndTime.Sub(tb.StartTime),
		tb.StartTime.In(s.Timezone()).Format(strFormat),
		s.params.longDistanceUnit.convertDistance(s.gpx.Length2D(), meter),
		s.params.distance(),
		s.gpx.GetTrackPointsNo(),
	)
}

func (s *Segment) ModeCounts() (moving, turning, static int) {
	for _, p := range s.Points {
		if p.Mode == Moving {
			moving++
		} else if p.Mode == Turning {
			turning++
		} else {
			static++
		}
	}
	return moving, turning, static
}

func (s *Segment) String() string {
	moving, turning, static := s.ModeCounts()
	return fmt.Sprintf("%.0fm/%.0fs @ %.1f/%.1f/%.1f %s \u2191 %d\u00b0/%d\u00b0 < %d\u00b0 %s (M:%d/T:%d/S:%d)",
		s.Distance, s.Duration.Seconds(),
		s.Speed.Min, s.Speed.Avg, s.Speed.Max, s.params.speedUnit.speed(),
		s.Heading.Min, s.Heading.Max, s.Heading.Variation,
		s.Mode, moving, turning, static)
}

func (s *Segment) ShortString() string {
	return fmt.Sprintf("%.0fm/%.0fs @ %.1f/%.1f/%.1f %s \u2191 %d\u00b0/%d\u00b0 %s",
		s.Distance, s.Duration.Seconds(),
		s.Speed.Min, s.Speed.Avg, s.Speed.Max, s.params.speedUnit.speed(),
		s.Heading.Min, s.Heading.Max,
		s.Mode)
}

type Segments []*Segment

func (ss Segments) gpxString() string {
	var all []string
	for _, s := range ss {
		all = append(all, s.gpxString())
	}
	return strings.Join(all, "\n")
}

// Sort segments by start time
func (s Segments) Len() int      { return len(s) }
func (s Segments) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s Segments) Less(i, j int) bool {
	return s[i].gpx.TimeBounds().StartTime.Before(s[j].gpx.TimeBounds().StartTime)
}

func (ss Segments) EachPair(f func(prev *Point, next *Point)) {
	var prev *Point
	for _, s := range ss {
		for _, p := range s.Points {
			if prev == nil {
				prev = p
				continue
			}
			f(prev, p)
			prev = p
		}
	}
}
