package main

import (
	"bytes"
	"fmt"
	"math"
	"time"

	"github.com/tkrajina/gpxgo/gpx"
)

const border = 0.05

// Map is used to translate track coordinates into SVG coordinate system for rendering.
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

// Point translates a GPS point into SVG coordinates.
func (m *Map) Point(p *gpx.GPXPoint) (x, y int) {
	y = m.h - int((p.Latitude-m.ly)/m.lh*float64(m.h))
	x = int((p.Longitude - m.lx) * m.coef / m.lw * float64(m.w))
	return
}

// Distance computes the distance between two GPS points in specified units.
func (m *Map) Distance(p1, p2 *gpx.GPXPoint, unit unit) float64 {
	x := p2.Latitude - p1.Latitude
	y := (p2.Longitude - p1.Longitude) * m.coef
	return float64(unit) * math.Sqrt(x*x+y*y)
}

// Speed computes the average speed between two GPS points in specified units of distance.
// The time aspect is derived from the distance unit, i.e. meter => m/s, km => km/h, nm => kts.
func (m *Map) Speed(p1, p2 *gpx.GPXPoint, unit unit) float64 {
	t := float64(p2.Timestamp.Sub(p1.Timestamp))
	if unit == meter {
		t /= float64(time.Second)
	} else {
		t /= float64(time.Hour)
	}
	return m.Distance(p1, p2, unit) / t
}

// Heading computes the direction from p1 to p2 in degrees.
func (m *Map) Heading(p1, p2 *gpx.GPXPoint) int {
	lat := p2.Latitude - p1.Latitude
	lon := (p2.Longitude - p1.Longitude) * m.coef
	deg := int(math.Round(math.Atan2(lon, lat) / math.Pi * 180))
	if deg < 0 {
		return 360 + deg
	} else {
		return deg
	}
}

var palette = func() (palette []int) {
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
}()

// SpeedColor return the RGB color code matching the speed between two GPS points.
func (m *Map) SpeedColor(speed float64) string {
	s := int(speed)
	if s >= len(palette) {
		s = len(palette) - 1
	}
	return fmt.Sprintf("#%03x", palette[s])
}

func (m *Map) polylinePoints(s *Segment) string {
	b := bytes.NewBuffer(nil)
	for i := range s.gpx.Points {
		x, y := m.Point(s.gpxPoint(i))
		fmt.Fprintf(b, "%d,%d ", x, y)
	}
	return b.String()
}
