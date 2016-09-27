package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
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

func (t *Track) Segment(i int) Segment {
	return Segment{&t.Segments[i]}
}

type Tracks []Track

func (t *Track) WriteMapFile(dir string) error {
	fn, err := t.FileName()
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(dir, fn+".svg"))
	if err != nil {
		return err
	}
	defer f.Close()
	m := NewMap(t.Bounds(), mapWidth)
	return m.RenderLines(f, t)
}

func (t *Track) WriteGpxFile(dir string) error {
	fn, err := t.FileName()
	if err != nil {
		return err
	}
	f, err := os.Create(filepath.Join(dir, fn+".gpx"))
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

func (t *Track) FileName() (string, error) {
	tb := t.TimeBounds()
	start := tb.StartTime.In(t.Timezone()).Format(fnFormat)
	d := tb.EndTime.Sub(tb.StartTime)
	return fmt.Sprintf("%s-%dh%02d-%04.1fnm",
			start,
			int(d.Hours()),
			int(d.Minutes())%60,
			t.Length2D()/1852),
		nil
}

func (t *Track) Extent(unit float64) (width, height float64) {
	b := t.Bounds()
	coef := math.Cos((b.MaxLatitude + b.MinLatitude) * math.Pi / 360)
	height = (b.MaxLongitude - b.MinLongitude) * coef * unit
	width = (b.MaxLatitude - b.MinLatitude) * unit
	return width, height
}

func (t *Track) String() string {
	tb := t.TimeBounds()
	w, h := t.Extent(km)
	return fmt.Sprintf("%s %.03fkm %.03fkmx%.03fkm (%s)",
		tb.StartTime.In(t.Timezone()).Format(strFormat),
		t.Length2D()/1000,
		w,
		h,
		tb.EndTime.Sub(tb.StartTime),
	)
}
