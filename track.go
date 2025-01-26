package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/bradfitz/latlong"
	"github.com/tkrajina/gpxgo/gpx"
)

const fnFormat = "060102"
const strFormat = "06-01-02 15:04:05"
const mapWidth = 1000

type Track struct {
	gpx      *gpx.GPXTrack
	tz       *time.Location
	filename string // file from which the track was collected
	// Analysis results
	params   *AnalysisParameters
	Segments Segments
	Distance float64
	Start    time.Time
	End      time.Time
	Duration time.Duration
}

// WriteMapFile generates an SVG map of the track into the specified directory.
func (t *Track) WriteMapFile(dir string) error {
	f, err := os.Create(filepath.Join(dir, t.FileName()+".svg"))
	if err != nil {
		return err
	}
	defer f.Close()
	m := NewMap(t.gpx.Bounds(), mapWidth)
	m.renderLines(f, t)
	return nil
}

// WriteSubtitleFile generates a video subtitles file with the stats.
func (t *Track) WriteSubtitleFile(dir string, offset time.Duration) error {
	f, err := os.Create(filepath.Join(dir, t.FileName()+".vtt"))
	if err != nil {
		return err
	}
	defer f.Close()
	m := NewMap(t.gpx.Bounds(), mapWidth)
	m.renderSubtitles(f, t, offset)
	return nil
}

// WriteGpxFile generates track's GPX file into the specified directory.
func (t *Track) WriteGpxFile(dir string) error {
	f, err := os.Create(filepath.Join(dir, t.FileName()+".gpx"))
	if err != nil {
		return err
	}
	defer f.Close()
	g := &gpx.GPX{}
	g.AppendTrack(t.gpx)
	b, err := g.ToXml(gpx.ToXmlParams{Version: "1.1", Indent: true})
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	return err
}

// Timezone returns the timezone for track's location.
func (t *Track) Timezone() *time.Location {
	if t.tz != nil {
		return t.tz
	}
	b := t.gpx.Bounds()
	lat := (b.MaxLatitude + b.MinLatitude) / 2
	lon := (b.MaxLongitude + b.MinLongitude) / 2
	var err error
	t.tz, err = time.LoadLocation(latlong.LookupZoneName(lat, lon))
	if err != nil {
		t.tz = time.UTC
	}
	return t.tz
}

// FileName generates a file name based on track's time bounds and length.
func (t *Track) FileName() string {
	return fmt.Sprintf("%s-%dh%02d-%04.1f%s",
		t.Start.In(t.Timezone()).Format(fnFormat),
		int(t.Duration.Hours()),
		int(t.Duration.Minutes())%60,
		t.params.asLongDistance(t.Distance),
		t.params.longDistance())
}

// Extent returns box dimensions of the track in specified units.
func (t *Track) Extent(unit unit) (width, height float64) {
	b := t.gpx.Bounds()
	coef := math.Cos((b.MaxLatitude + b.MinLatitude) * math.Pi / 360)
	height = (b.MaxLongitude - b.MinLongitude) * coef * float64(unit)
	width = (b.MaxLatitude - b.MinLatitude) * float64(unit)
	return width, height
}

// String returns track description.
func (t *Track) String() string {
	tb := t.gpx.TimeBounds()
	w, h := t.Extent(t.params.longDistanceUnit)
	unit := t.params.longDistance()
	return fmt.Sprintf("%s %05.2f%s %05.2f%s x %05.2f%s (%s) [%d segments]",
		tb.StartTime.In(t.Timezone()).Format(strFormat),
		t.params.asLongDistance(t.Distance),
		unit,
		w,
		unit,
		h,
		unit,
		tb.EndTime.Sub(tb.StartTime),
		len(t.Segments),
	)
}

func (t *Track) gpxAnalyze(params *AnalysisParameters) {
	t.params = params
	var segments Segments
	tMap := NewMap(t.gpx.Bounds(), 1000)
	for i := range t.gpx.Segments {
		segment := &t.gpx.Segments[i]
		segments = append(segments, gpxAnalyzeSegment(segment, t.filename, tMap, params)...)
	}
	var distance float64
	for _, s := range segments {
		distance += s.Distance
	}
	t.Segments = segments
	t.gpx = &gpx.GPXTrack{}
	for _, s := range t.Segments {
		t.gpx.AppendSegment(s.gpx)
	}
	t.Start = segments[0].Start
	t.End = segments[len(segments)-1].End
	t.Duration = t.End.Sub(t.Start)
	t.Distance = distance
}

func (t *Track) posClassify(windDirection direction) {
	windDirection = t.windDirection(windDirection)
	for _, s := range t.Segments {
		if s.Mode == Moving {
			s.Type = windDirection.pointOfSail(s.Heading.Mid)
		} else if s.Mode == Turning {
			from := windDirection.pointOfSail(s.Points[0].Heading)
			to := windDirection.pointOfSail(s.Points[len(s.Points)-1].Heading)
			s.Type = windDirection.turnType(from, to)
		} else {
			s.Type = drifting
		}
	}
}

// Try to determine the prevailing wind direction by
// combining heading ranges of all moving segments,
// finding the gaps that are more than 80 degrees wide,
// and picking the one that includes the defaultDirection,
// or picking the largest one.
func (t *Track) windDirection(defaultDirection direction) direction {
	headings := HeadingSet{}
	for _, s := range t.Segments {
		if s.Speed.Min < 2*Sailing.movingSpeed {
			continue
		}
		if _, turning, static := s.ModeCounts(); turning+static > 0 {
			continue
		}
		headings = headings.Add(s.Heading)
		// fmt.Println(s.String())
		// fmt.Println(headings.String())
	}

	var candidates HeadingSet
	for _, hr := range headings.Inverse() {
		if hr.Variation > 60 {
			candidates = append(candidates, hr)
		}
	}
	if len(candidates) == 0 { // no candidates
		return defaultDirection
	}
	// sort by range width
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Variation > candidates[j].Variation
	})
	// fmt.Println("candidates: ", candidates.String())
	// pick the one that includes defaultDirection
	for _, hr := range candidates {
		if hr.Includes(int(defaultDirection)) {
			return Direction(hr.Mid)
		}
	}
	return Direction(candidates[0].Mid)
}

type Tracks []Track

func (ts Tracks) String() string {
	var ss []string
	for _, t := range ts {
		ss = append(ss, t.String())
	}
	return strings.Join(ss, "\n")
}
