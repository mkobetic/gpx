package main

import (
	"fmt"
	"time"

	"github.com/tkrajina/gpxgo/gpx"
)

type Mode string

const (
	Static  Mode = "static"
	Moving  Mode = "moving"
	Turning Mode = "turning"
)

// Point is a track point with computed motion parameters and analysis results.
// Note that the starting point has zero Speed and Heading!
type Point struct {
	gpx *gpx.GPXPoint
	// Analysis results
	params        *AnalysisParameters
	previous      *Point  // previous point on the track
	next          *Point  // next point on the track
	Speed         float64 // speed from previous point
	Heading       int     // heading from previous point
	Distance      float64 // distance from previous point
	HeadingChange int     // how much does the heading change in the lookAround range
	Mode          Mode
}

func (p *Point) String() string {
	return fmt.Sprintf("%0.1fm @ %0.1f %s \u2191 %d\u00b0 %s < %d (%s)", p.Distance, p.Speed, p.params.speedUnit.speed(), p.Heading, Direction(p.Heading).String(), p.HeadingChange, p.Mode)
}

func (p *Point) ShortString() string {
	return fmt.Sprintf("%0.1fm @ %0.1f %s \u2191 %d\u00b0 %s", p.Distance, p.Speed, p.params.speedUnit.speed(), p.Heading, Direction(p.Heading).String())
}

func (p *Point) Analyze(params *AnalysisParameters) {
	var speed float64
	if p.previous == nil {
		p.HeadingChange = p.next.headingChange(0, params.lookAround-p.next.Distance, true, params.movingSpeed)
		speed = p.next.Speed
	} else {
		p.HeadingChange = p.headingChange(0, params.lookAround, true, params.movingSpeed) - p.headingChange(0, params.lookAround, false, params.movingSpeed)
		speed = p.Speed
	}
	if speed < params.movingSpeed {
		p.Mode = Static
		return
	}
	if abs(p.HeadingChange) < params.turningChange {
		p.Mode = Moving
	} else {
		p.Mode = Turning
	}
}

func (p *Point) headingChange(change int, distance float64, forward bool, movingSpeed float64) int {
	var next *Point
	if forward {
		next = p.next
	} else {
		next = p.previous
	}
	if next == nil || next.Speed < movingSpeed {
		return change
	}
	change += headingDiff(p.Heading, next.Heading)
	distance -= next.Distance
	if distance <= 0 {
		return change
	}
	return next.headingChange(change, distance, forward, movingSpeed)
}

type Points []*Point

func (ps Points) mode() Mode {
	modes := make(map[Mode]int)
	for _, p := range ps {
		modes[p.Mode] += 1
	}
	var maxn int
	var maxm Mode
	for m, n := range modes {
		if n > maxn {
			maxm = m
		}
	}
	return maxm
}

func (ps Points) speed(mode Mode) SpeedRange {
	var sum float64
	var count int
	start := 0
	// Ignore stray points with different mode
	// to avoid skewing the statistics.
	for ; ps[start].Mode != mode; start++ {
	}
	min := ps[start].Speed
	max := min
	for _, p := range ps[start:] {
		if p.Mode != mode {
			continue
		}
		count++
		sum += p.Speed
		if p.Speed < min {
			min = p.Speed
		}
		if max < p.Speed {
			max = p.Speed
		}
	}
	return SpeedRange{Min: min, Avg: sum / float64(len(ps)), Max: max}
}

func (ps Points) distance() float64 {
	var sum float64
	for _, p := range ps {
		sum += p.Distance
	}
	return sum
}

func (ps Points) duration() time.Duration {
	return ps.end().Sub(ps.start())
}

func (ps Points) end() time.Time {
	return ps[len(ps)-1].gpx.Timestamp
}

func (ps Points) start() time.Time {
	return ps[0].gpx.Timestamp
}

func (ps Points) headingRange(mode Mode) *HeadingRange {
	start := 0
	// Ignore stray points with different mode
	// to avoid skewing the statistics.
	for ; ps[start].Mode != mode; start++ {
	}
	prev := ps[start]
	min := prev.Heading
	max := min
	start++
	for _, p := range ps[start:] {
		if p.Mode != mode {
			continue
		}
		diff := headingDiff(prev.Heading, p.Heading)
		if diff < 0 && headingDiff(p.Heading, min) > 0 {
			min = p.Heading
		} else if diff > 0 && headingDiff(p.Heading, max) < 0 {
			max = p.Heading
		}
		prev = p
	}
	return NewHeadingRange(min, max)
}

// Split points into longest runs by mode.
func (ps Points) eachRun(f func(run Points, mode Mode)) {
	for {
		mode := ps[0].Mode
		var i int
		for i = 1; i < len(ps) && ps[i].Mode == mode; i++ {
		}
		// Don't want segments with single point
		if i == 1 && len(ps) > 1 {
			mode := ps[1].Mode
			for i = 2; i < len(ps) && ps[i].Mode == mode; i++ {
			}
		}
		// Check that we're not left with a single point run at the end
		if i+1 == len(ps) {
			i += 1
		}
		f(ps[0:i], mode)
		if i == len(ps) {
			break
		}
		ps = ps[i:]
	}
}

// Difference in degrees from heading a to heading b (-179 ... 180),
// positive if turning clockwise, negative counter-clockwise.
func headingDiff(a, b int) int {
	diff := b - a
	if diff > 180 {
		diff -= 360
	} else if diff <= -180 {
		diff += 360
	}
	return diff
}

// Is heading h between heading min and max inclusive (clockwise)
func headingBetween(h, min, max int) bool {
	minDiff := headingDiff(h, min)
	maxDiff := headingDiff(h, max)
	if headingDist(min, max) < 180 {
		return minDiff <= 0 && maxDiff >= 0
	} else {
		return !(minDiff > 0 && maxDiff < 0)
	}
}

// The distance in degrees between min and max inclusive (0-359) (clockwise)
func headingDist(min, max int) int {
	diff := headingDiff(min, max)
	if diff >= 0 {
		return diff
	}
	return 360 + diff
}

// Add diff d (-179...180) to heading h
func headingAdd(h, d int) int {
	h2 := h + d
	if h2 > 359 {
		return h2 - 360
	} else if h2 < 0 {
		return 360 + h2
	} else {
		return h2
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
