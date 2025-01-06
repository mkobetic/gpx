package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/bradfitz/latlong"
	"github.com/tkrajina/gpxgo/gpx"
)

const fnFormat = "060102"
const strFormat = "06-01-02 15:04:05"
const mapWidth = 1000

type Track struct {
	*gpx.GPXTrack
	tz *time.Location
}

// Segment returns i-th segment of the track.
func (t *Track) Segment(i int) Segment {
	return Segment{&t.Segments[i]}
}

type Tracks []Track

func (ts Tracks) String() string {
	var ss []string
	for _, t := range ts {
		ss = append(ss, t.String())
	}
	return strings.Join(ss, "\n")
}

// WriteMapFile generates an SVG map of the track into the specified directory.
func (t *Track) WriteMapFile(dir string) error {
	f, err := os.Create(filepath.Join(dir, t.FileName()+".svg"))
	if err != nil {
		return err
	}
	defer f.Close()
	m := NewMap(t.Bounds(), mapWidth)
	m.renderLines(f, t)
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
	g.AppendTrack(t.GPXTrack)
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
	b := t.Bounds()
	var err error
	t.tz, err = time.LoadLocation(latlong.LookupZoneName(b.MinLatitude, b.MinLongitude))
	if err != nil {
		t.tz = time.UTC
	}
	return t.tz
}

// FileName generates a file name based on track's time bounds and length.
func (t *Track) FileName() string {
	tb := t.TimeBounds()
	start := tb.StartTime.In(t.Timezone()).Format(fnFormat)
	d := tb.EndTime.Sub(tb.StartTime)
	return fmt.Sprintf("%s-%dh%02d-%04.1fnm",
		start,
		int(d.Hours()),
		int(d.Minutes())%60,
		t.Length2D()/1852)
}

// Extent returns box dimensions of the track in specified units.
func (t *Track) Extent(unit float64) (width, height float64) {
	b := t.Bounds()
	coef := math.Cos((b.MaxLatitude + b.MinLatitude) * math.Pi / 360)
	height = (b.MaxLongitude - b.MinLongitude) * coef * unit
	width = (b.MaxLatitude - b.MinLatitude) * unit
	return width, height
}

// String returns track description.
func (t *Track) String() string {
	tb := t.TimeBounds()
	w, h := t.Extent(nm)
	return fmt.Sprintf("%s %05.2fnm %05.2fnm x %05.2fnm (%s)",
		tb.StartTime.In(t.Timezone()).Format(strFormat),
		t.Length2D()/1852,
		w,
		h,
		tb.EndTime.Sub(tb.StartTime),
	)
}
