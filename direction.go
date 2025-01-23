package main

import "math"

type direction int
type pointOfSail int

const (
	N   direction = 0
	NNE direction = 23
	NE  direction = 45
	ENE direction = 68
	E   direction = 90
	ESE direction = 113
	SE  direction = 135
	SSE direction = 158
	S   direction = 180
	SSW direction = 203
	SW  direction = 225
	WSW direction = 248
	W   direction = 270
	WNW direction = 293
	NW  direction = 315
	NNW direction = 338

	// value is heading assuming wind from N
	irons   pointOfSail = 0
	closePT pointOfSail = 45
	beamPT  pointOfSail = 90
	broadPT pointOfSail = 135
	run     pointOfSail = 180
	broadSB pointOfSail = 225
	beamSB  pointOfSail = 270
	closeSB pointOfSail = 315
)

var directionToString = map[direction]string{
	N:   "N",
	NNE: "NNE",
	NE:  "NE",
	ENE: "ENE",
	E:   "E",
	ESE: "ESE",
	SE:  "SE",
	SSE: "SSE",
	S:   "S",
	SSW: "SSW",
	SW:  "SW",
	WSW: "WSW",
	W:   "W",
	WNW: "WNW",
	NW:  "NW",
	NNW: "NNW",
}

var posToString = map[pointOfSail]string{
	irons:   "parked",
	closePT: "close reach port",
	beamPT:  "beam reach port",
	broadPT: "broad reach port",
	run:     "downwind run",
	closeSB: "close reach starboard",
	beamSB:  "beam reach starboard",
	broadSB: "broad reach starboard",
}

func Direction(heading int) direction {
	idx := int(math.Floor((float64(heading) + 11.25) / 22.5))
	if idx > 15 {
		idx = 0
	}
	return []direction{N, NNE, NE, ENE, E, ESE, SE, SSE, S, SSW, SW, WSW, W, WNW, NW, NNW}[idx]
}

func (d direction) String() string {
	return directionToString[d]
}

func (windDirection direction) pointOfSail(heading int) pointOfSail {
	// adjust heading for wind direction
	adjusted := heading - int(windDirection)
	if adjusted < 0 {
		adjusted = 360 + adjusted
	}
	idx := int(math.Floor((float64(adjusted) + 22.5) / 45))
	if idx > 7 {
		idx = 0
	}
	return []pointOfSail{irons, closePT, beamPT, broadPT, run, broadSB, beamSB, closeSB}[idx]
}

func (pos pointOfSail) String() string {
	return posToString[pos]
}
