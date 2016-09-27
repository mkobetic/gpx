package main

import (
	"bytes"
	"fmt"
	"math"
	"time"

	"github.com/tkrajina/gpxgo/gpx"
)

const border = 0.05

type Map struct {
	w, h   int     // widht, height in svg coordinates
	lw, lh float64 // width, height in lat/lon degrees
	lx, ly float64 // bottom/left offset in lat/lon degrees
	coef   float64 // longitudinal adjustment coeficient (degrees of longitude are shorter in higher latitudes)
}

func NewMap(b gpx.GpxBounds, width int) *Map {
	m := &Map{w: width, lx: b.MinLongitude, ly: b.MinLatitude}
	m.coef = math.Cos((b.MaxLatitude + b.MinLatitude) * math.Pi / 360)
	m.lw = (b.MaxLongitude - b.MinLongitude) * m.coef
	m.lh = b.MaxLatitude - b.MinLatitude
	bx := border * m.lh
	m.lh += 2 * bx
	m.lx -= bx
	m.lw += 2 * bx
	m.ly -= bx
	m.h = int(m.lh / m.lw * float64(m.w))
	return m
}

func (m *Map) Point(p *gpx.GPXPoint) (x, y int) {
	y = m.h - int((p.Latitude-m.ly)/m.lh*float64(m.h))
	x = int((p.Longitude - m.lx) * m.coef / m.lw * float64(m.w))
	return
}

// units for Distance and Speed functions.
const km = 2 * math.Pi * 6371 / 360
const meter = 1000 * km
const nm = 2 * math.Pi * 3440 / 360

func (m *Map) Distance(p1, p2 *gpx.GPXPoint, unit float64) float64 {
	x := p2.Latitude - p1.Latitude
	y := (p2.Longitude - p1.Longitude) * m.coef
	return unit * math.Sqrt(x*x+y*y)
}

func (m *Map) Speed(p1, p2 *gpx.GPXPoint, unit float64) float64 {
	t := float64(p2.Timestamp.Sub(p1.Timestamp))
	if unit == meter {
		t /= float64(time.Second)
	} else {
		t /= float64(time.Hour)
	}
	return m.Distance(p1, p2, unit) / t
}

func compute_palette() (palette []int) {
	for i := 0; i < 16; i += 2 {
		palette = append(palette, i*16+15)
	}
	for i := 0; i < 16; i += 4 {
		palette = append(palette, 15*16+15-i)
	}
	for i := 0; i < 16; i += 4 {
		palette = append(palette, (i*16+15)*16)
	}
	for i := 0; i < 16; i += 2 {
		palette = append(palette, (17*15-i)*16)
	}
	return palette
}

var palette = compute_palette()

func (m *Map) SpeedColor(p1, p2 *gpx.GPXPoint) string {
	s := int(m.Speed(p1, p2, nm))
	if s >= len(palette) {
		s = len(palette) - 1
	}
	return fmt.Sprintf("#%03x", palette[s])
}

func (m *Map) PolylinePoints(s Segment) string {
	b := bytes.NewBuffer(nil)
	for i := range s.Points {
		x, y := m.Point(s.Point(i))
		fmt.Fprintf(b, "%d,%d ", x, y)
	}
	return b.String()
}
