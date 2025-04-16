package main

import (
	_ "embed"
	"fmt"
	"math"
	"time"

	"github.com/tkrajina/gpxgo/gpx"
)

const border = 20      // map padding area width in points of SVG coordinates
const tlUnitHeight = 3 // timeline height of a vertical unit (e.g. a knot) in points of SVG coordinates
const tlHeight = 25 * tlUnitHeight

//go:embed map.js
var script string

// Map is used to translate track coordinates into SVG coordinate system for rendering.
type Map struct {
	w, h   float64 // width, height in svg coordinates derived from lw/lh
	lw, lh float64 // width, height in lat/lon degrees
	lx, ly float64 // bottom/left offset in lat/lon degrees
	coef   float64 // longitudinal adjustment coeficient (degrees of longitude are shorter in higher latitudes)
}

func NewMap(b gpx.GpxBounds, unit unit) *Map {
	m := &Map{lx: b.MinLongitude, ly: b.MinLatitude}
	// calculate the coeficient for longitudinal adjustment
	m.coef = math.Cos((b.MaxLatitude + b.MinLatitude) * math.Pi / 360)
	// calculate lat/long dimensions
	m.lw = (b.MaxLongitude - b.MinLongitude) * m.coef
	m.lh = b.MaxLatitude - b.MinLatitude
	// compute SVG dimensions as lat/long dimensions * unit
	// i.e. one point in SVG coordinates is 1 unit
	m.w = m.lw*m.coef*float64(unit) + 2*border
	m.h = m.lh/m.lw*float64(m.w) + 2*border
	return m
}

// Point translates a GPS point into SVG coordinates.
func (m *Map) Point(p *gpx.GPXPoint) (x, y int) {
	y = int(m.h-(p.Latitude-m.ly)/m.lh*float64(m.h)) + border
	x = int((p.Longitude-m.lx)*m.coef/m.lw*float64(m.w)) + border
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
