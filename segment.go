package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bradfitz/latlong"
	"github.com/tkrajina/gpxgo/gpx"
)

type Segment struct {
	gpx      *gpx.GPXTrackSegment
	filename string
	Points   Points
	// Analysis results
	AvgHeading       int
	HeadingVariation int
	Distance         float64
	AvgSpeed         float64
	Duration         time.Duration
	Start            time.Time
	End              time.Time
	Mode             Mode
	// POS              pointOfSail
}

func SegmentFromPoints(ps Points, mode Mode, filename string) *Segment {
	gpxPoints := make([]gpx.GPXPoint, len(ps))
	for i, p := range ps {
		gpxPoints[i] = *p.gpx
	}
	return &Segment{
		gpx:              &gpx.GPXTrackSegment{Points: gpxPoints},
		filename:         filename,
		Points:           ps,
		AvgHeading:       ps.averageHeading(),
		HeadingVariation: ps.headingVariation(),
		Distance:         ps.distance(),
		AvgSpeed:         ps.averageSpeed(),
		Start:            ps.start(),
		End:              ps.end(),
		Duration:         ps.duration(),
		Mode:             mode,
	}
}

// gpxPoint returns i-th point of the segment.
func (s *Segment) gpxPoint(i int) *gpx.GPXPoint {
	return &s.gpx.Points[i]
}

// gpxEachPair iterates over a segment with pairs of subsequent points.
func (s *Segment) gpxEachPair(f func(prev, next *gpx.GPXPoint)) {
	prev := s.gpxPoint(0)
	for i := 1; i < len(s.gpx.Points); i++ {
		next := s.gpxPoint(i)
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
	return fmt.Sprintf("%s @ %s = %05.2fnm (%d)",
		tb.EndTime.Sub(tb.StartTime),
		tb.StartTime.In(s.Timezone()).Format(strFormat),
		s.gpx.Length2D()/1852,
		s.gpx.GetTrackPointsNo(),
	)
}

func (s *Segment) String() string {
	return fmt.Sprintf("%.0fm/%.0fs @ %.1f kts \u2191 %d\u00b0 %s < %d\u00b0 %s (%d)",
		s.Distance, s.Duration.Seconds(), s.AvgSpeed, s.AvgHeading, Direction(s.AvgHeading).String(), s.HeadingVariation, s.Mode, len(s.Points))

}

// gpxSplit segment where the time difference between points is more than @limit.
func (s *Segment) gpxSplit(limit time.Duration) Segments {
	limitSeconds := limit.Seconds()
	prev := s.gpxPoint(len(s.gpx.Points) - 1)
	for i := len(s.gpx.Points) - 2; i >= 0; i-- {
		next := s.gpxPoint(i)
		if next.TimeDiff(prev) > limitSeconds {
			s1, s2 := s.gpx.Split(i)
			return append((&Segment{gpx: s1, filename: s.filename}).gpxSplit(limit), &Segment{gpx: s2, filename: s.filename})
		}
		prev = next
	}
	return Segments{s}
}

// gpxAnalyze the segment and split it up into runs of points of the same Mode of movement (static, moving, turning).
// The Map is the context to use for the analysis derived from the Track, it is the same for all segments of the track.
func (s *Segment) gpxAnalyze(m *Map, params *AnalysisParameters) Segments {
	previous := &Point{gpx: s.gpxPoint(0)}
	points := Points{previous}
	s.gpxEachPair(func(prev, next *gpx.GPXPoint) {
		nextPoint := &Point{
			gpx:      next,
			previous: previous,
			Heading:  m.Heading(prev, next),
			Distance: m.Distance(prev, next, params.distanceUnit),
			Speed:    m.Speed(prev, next, params.speedUnit),
		}
		previous.next = nextPoint
		previous = nextPoint
		points = append(points, nextPoint)
	})
	for _, p := range points {
		p.Analyze(params)
	}
	segments := Segments{}
	points.eachRun(func(run Points, mode Mode) {
		segments = append(segments, SegmentFromPoints(run, mode, s.filename))
	})
	return segments
}

type Segments []*Segment

func gpxGetSegments(g *gpx.GPX, filename string) (s Segments) {
	for _, t := range g.Tracks {
		for i := range t.Segments {
			s = append(s, &Segment{gpx: &t.Segments[i], filename: filename})
		}
	}
	return s
}

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

// gpxDedupe removes subsequent segments with the same time bounds
// and segments that have less than @min points.
func (s Segments) gpxDedupe(min int) (t Segments) {
	if len(s) == 0 {
		return
	}
	p := s[0]
	t = append(t, p)
	for _, s := range s[1:] {
		if s.gpx.TimeBounds().Equals(p.gpx.TimeBounds()) {
			continue
		}
		if s.gpx.GetTrackPointsNo() > min {
			t = append(t, s)
		}
		p = s
	}
	return
}

func (s Segments) gpxSplit(limit time.Duration) (t Segments) {
	for _, seg := range s {
		t = append(t, seg.gpxSplit(limit)...)
	}
	return t
}

// gpxTracks creates tracks from subsequent segments with time bounds that
// are less than limit time apart.
func (ss Segments) gpxTracks(limit time.Duration) (tracks Tracks) {
	if len(ss) == 0 {
		return
	}
	p := ss[0]
	t := new(gpx.GPXTrack)
	t.AppendSegment(p.gpx)
	for _, s := range ss[1:] {
		if s.gpx.TimeBounds().StartTime.Sub(p.gpx.TimeBounds().EndTime) > limit {
			tracks = append(tracks, Track{gpx: t, filename: s.filename})
			t = new(gpx.GPXTrack)
		}
		t.AppendSegment(s.gpx)
		p = s
	}
	tracks = append(tracks, Track{gpx: t, filename: p.filename})
	return
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
