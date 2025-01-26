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

type Point struct {
	gpx      *gpx.GPXPoint
	previous *Point  // previous point on the track
	next     *Point  // next point on the track
	Speed    float64 // speed from previous point
	Heading  int     // heading from previous point
	Distance float64 // distance from previous point
	// Analysis results
	HeadingChange int // how much does the heading change in the lookAround range
	Mode          Mode
}

func (p *Point) String() string {
	return fmt.Sprintf("%0.1fm @ %0.1f kts \u2191 %d\u00b0 %s < %d (%s)", p.Distance, p.Speed, p.Heading, Direction(p.Heading).String(), p.HeadingChange, p.Mode)
}

func (p *Point) ShortString() string {
	return fmt.Sprintf("%0.1fm @ %0.1f kts \u2191 %d\u00b0 %s", p.Distance, p.Speed, p.Heading, Direction(p.Heading).String())
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

func (ps Points) speed() Speed {
	var sum float64
	min := ps[0].Speed
	max := min
	for _, p := range ps {
		sum += p.Speed
		if p.Speed < min {
			min = p.Speed
		}
		if max < p.Speed {
			max = p.Speed
		}
	}
	return Speed{Min: min, Avg: sum / float64(len(ps)), Max: max}
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

func (ps Points) heading() Heading {
	var maxDiff, minDiff, currentDiff int
	min := ps[0].Heading
	max := min
	for _, p := range ps[1:] {
		diff := headingDiff(p.previous.Heading, p.Heading)
		if diff < 0 && headingDiff(p.Heading, min) > 0 {
			min = p.Heading
		} else if diff > 0 && headingDiff(p.Heading, max) < 0 {
			max = p.Heading
		}
		currentDiff += diff
		if currentDiff > maxDiff {
			maxDiff = currentDiff
		} else if currentDiff < minDiff {
			minDiff = currentDiff
		}
	}
	diff := maxDiff - minDiff
	mid := (min + diff/2) % 360
	return Heading{
		Min:       min,
		Mid:       mid,
		Max:       max,
		Variation: diff,
	}
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
