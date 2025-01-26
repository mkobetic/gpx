package main

import (
	"time"

	"github.com/tkrajina/gpxgo/gpx"
)

// Collects original segments from a GPX file.
// Attaches the filename to the segments.
func gpxGetSegments(g *gpx.GPX, filename string) (ss Segments) {
	for _, t := range g.Tracks {
		for i := range t.Segments {
			ss = append(ss, &Segment{gpx: &t.Segments[i], filename: filename})
		}
	}
	return ss
}

// gpxDedupe removes subsequent segments with the same time bounds
// and segments that have less than @min points.
func gpxDedupeSegments(ss Segments, min int) (t Segments) {
	if len(ss) == 0 {
		return
	}
	p := ss[0]
	t = append(t, p)
	for _, s := range ss[1:] {
		points := s.gpx.GetTrackPointsNo()
		bounds := s.gpx.TimeBounds()
		pBounds := p.gpx.TimeBounds()
		pPoints := p.gpx.GetTrackPointsNo()
		if bounds.Equals(pBounds) && points == pPoints {
			continue
		}
		if points > min {
			t = append(t, s)
		}
		p = s
	}
	return t
}

func gpxSplitSegments(ss Segments, limit time.Duration) (t Segments) {
	for _, seg := range ss {
		t = append(t, gpxSplitSegment(seg, limit)...)
	}
	return t
}

// gpxTracks reassembles tracks from subsequent original segments
// with time bounds that are less than limit time apart.
func gpxBuildTracks(ss Segments, limit time.Duration) (tracks Tracks) {
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

// gpxSplit segment where the time difference between points is more than @limit.
func gpxSplitSegment(s *Segment, limit time.Duration) Segments {
	limitSeconds := limit.Seconds()
	prev := s.gpxPoint(len(s.gpx.Points) - 1)
	for i := len(s.gpx.Points) - 2; i >= 0; i-- {
		next := s.gpxPoint(i)
		if next.TimeDiff(prev) > limitSeconds {
			s1, s2 := s.gpx.Split(i)
			return append(gpxSplitSegment((&Segment{gpx: s1, filename: s.filename}), limit), &Segment{gpx: s2, filename: s.filename})
		}
		prev = next
	}
	return Segments{s}
}

func gpxEachPair(s *gpx.GPXTrackSegment, f func(prev, next *gpx.GPXPoint)) {
	prev := &s.Points[0]
	for i := 1; i < len(s.Points); i++ {
		next := &s.Points[i]
		f(prev, next)
		prev = next
	}
}

// gpxAnalyze the segment and split it up into runs of points of the same Mode of movement (static, moving, turning).
// The Map is the context to use for the analysis derived from the Track, it is the same for all segments of the track.
func gpxAnalyzeSegment(s *gpx.GPXTrackSegment, filename string, m *Map, params *AnalysisParameters) Segments {
	previousPt := &Point{gpx: &s.Points[0], params: params}
	points := Points{previousPt}
	gpxEachPair(s, func(prev, next *gpx.GPXPoint) {
		nextPoint := &Point{
			gpx:      next,
			params:   params,
			previous: previousPt,
			Heading:  m.Heading(prev, next),
			Distance: m.Distance(prev, next, params.distanceUnit),
			Speed:    m.Speed(prev, next, params.speedUnit),
		}
		previousPt.next = nextPoint
		previousPt = nextPoint
		points = append(points, nextPoint)
	})
	// Fix up movement params of the first point (we don't want the Speed and Heading to be zero)
	first := points[0]
	first.Speed = first.next.Speed
	first.Heading = first.next.Heading
	// Run point analysis (determine Mode and HeadingChange)
	for _, p := range points {
		p.Analyze(params)
	}
	// Create segments of points in the same Mode.
	// Combine segments that are too short (< 5 points) with neighbouring segments.
	segments := Segments{}
	var previousSeg *Segment
	var short Points // holds previous run that was too short
	points.eachRun(func(run Points, mode Mode) {
		// join stashed previous short run
		run = append(short, run...)
		// if run was short then recompute mode
		if len(run)/2 < len(short) {
			mode = run.mode()
		}
		// if we are still short stash it and continue
		if len(run) < 5 {
			short = run
			return
		}
		short = nil // clear the short stash
		segment := SegmentFromPoints(run, mode, filename, params)
		segments = append(segments, segment)
		if previousSeg != nil {
			previousSeg.next = segment
			segment.previous = previousSeg
		}
		previousSeg = segment
	})
	if short != nil {
		// We were left with a shortie at the end, append it to the last segment.
		lastIndex := len(segments) - 1
		last := segments[lastIndex]
		segments[lastIndex] = SegmentFromPoints(append(last.Points, short...), last.Mode, filename, params)
	}
	return segments
}
