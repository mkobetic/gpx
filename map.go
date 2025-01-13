package main

import (
	"bytes"
	"fmt"
	"io"
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

// units for Distance and Speed functions,
// expressed as the length of one degree of longitude at the equator
const km = 2 * math.Pi * 6371 / 360
const meter = 1000 * km
const nm = 60

// Distance computes the distance between two GPS points in specified units.
func (m *Map) Distance(p1, p2 *gpx.GPXPoint, unit float64) float64 {
	x := p2.Latitude - p1.Latitude
	y := (p2.Longitude - p1.Longitude) * m.coef
	return unit * math.Sqrt(x*x+y*y)
}

// Speed computes the average speed between two GPS points in specified units of distance.
// The time aspect is derived from the distance unit, i.e. meter => m/s, km => km/h, nm => kts.
func (m *Map) Speed(p1, p2 *gpx.GPXPoint, unit float64) float64 {
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
func (m *Map) SpeedColor(p1, p2 *gpx.GPXPoint) string {
	s := int(m.Speed(p1, p2, nm))
	if s >= len(palette) {
		s = len(palette) - 1
	}
	return fmt.Sprintf("#%03x", palette[s])
}

func direction(heading int) string {
	idx := int(math.Floor((float64(heading) + 11.25) / 22.5))
	if idx > 15 {
		idx = 0
	}
	return []string{"N", "NNE", "NE", "ENE", "E", "ESE", "SE", "SSE", "S", "SSW", "SW", "WSW", "W", "WNW", "NW", "NNW"}[idx]
}

func (m *Map) polylinePoints(s *Segment) string {
	b := bytes.NewBuffer(nil)
	for i := range s.gpx.Points {
		x, y := m.Point(s.Point(i))
		fmt.Fprintf(b, "%d,%d ", x, y)
	}
	return b.String()
}

// Renders a VTT subtitle file based on the track.
// Positive @videoOffset means the video starts ahead of the track, the timestamps will be adjusted accordingly.
// Negative @videoOffset means the video starts later and therefore the corresponding initial part of the track will be skipped.
// See https://developer.mozilla.org/en-US/docs/Web/API/WebVTT_API/Web_Video_Text_Tracks_Format
func (m *Map) renderSubtitles(w io.Writer, t *Track, videoOffset time.Duration) {
	fmt.Fprintln(w, "WEBVTT")
	fmt.Fprintln(w)
	fmt.Fprintf(w, "NOTE generated by gpx at %s\n", time.Now().Format(time.RFC3339))
	fmt.Fprintf(w, "source = %s\n", t.filename)
	fmt.Fprintf(w, "video offset from source: %s\n", videoOffset)
	fmt.Fprintln(w)

	currentOffset := videoOffset
	totalDistance := float64(0)
	cueCounter := 0
	for i := range t.gpx.Segments {
		t.Segment(i).EachPair(func(prev, next *gpx.GPXPoint) {
			duration := next.Timestamp.Sub(prev.Timestamp)
			newOffset := currentOffset + duration
			if newOffset < 0 {
				currentOffset = newOffset
				return
			}
			cueCounter++
			totalDistance += m.Distance(prev, next, nm)
			heading := m.Heading(prev, next)
			direction := direction(heading)
			fmt.Fprintf(w, "%d\n", cueCounter)
			fmt.Fprintf(w, "%s --> %s\n", vttTimestamp(currentOffset), vttTimestamp(newOffset))
			fmt.Fprintf(w, "%s: %0.1f m @ %0.1f kts \u2191 %d\u00b0 %s = %0.2f nm\n",
				next.Timestamp.In(t.Timezone()).Format(time.TimeOnly),
				m.Distance(prev, next, meter),
				m.Speed(prev, next, nm),
				heading,
				direction,
				totalDistance)
			fmt.Fprintln(w)
			currentOffset = newOffset
		})
	}
}

func vttTimestamp(ts time.Duration) string {
	if ts < 0 {
		return "00:00:00.000"
	}
	total := ts.Milliseconds()
	ms := total % 1000
	total /= 1000
	s := total % 60
	total /= 60
	m := total % 60
	total /= 60
	return fmt.Sprintf("%02d:%02d:%02d.%03d", total, m, s, ms)
}
