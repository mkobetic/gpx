package main

import (
	"math"
	"sort"
)

func init() {
	for _, wd := range directionToString {
		WindDirections = append(WindDirections, wd)
	}
	sort.Strings(WindDirections)
}

type direction int
type pointOfSail int
type turn int
type tack int
type windAttitude int

const (
	UNK direction = -1 // UNKNOWN
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
	UNK: "UNK",
}

var WindDirections []string

var posToString = map[pointOfSail]string{
	irons:   "irons",
	closePT: "close reach port",
	beamPT:  "beam reach port",
	broadPT: "broad reach port",
	run:     "downwind run",
	closeSB: "close reach starboard",
	beamSB:  "beam reach starboard",
	broadSB: "broad reach starboard",
}

const (
	starboard tack = 1
	unknown   tack = 0
	port      tack = -1

	upwind   windAttitude = 1
	beam     windAttitude = 0
	downwind windAttitude = -1

	// turn bits:
	// 0 - FROM DOWNWIND=0 / UPWIND=1
	// 1 - FROM PT=0 / SB=1
	// 2 - TO DOWNWIND=0 / UPWIND=1
	// 3 - TO PT=0 / SB=1
	bearawayPT turn = 0b0001
	bearawaySB turn = 0b1011
	roundupPT  turn = 0b0100
	roundupSB  turn = 0b1110

	tackSBPT turn = 0b0111
	tackPTSB turn = 0b1101
	gybeSBPT turn = 0b0010
	gybePTSB turn = 0b1000

	drifting turn = 0b0000
)

var turnToString = map[turn]string{
	bearawayPT: "bear away port",
	bearawaySB: "bear away starboard",
	roundupPT:  "round up port",
	roundupSB:  "round up starboard",
	tackPTSB:   "tack port to starboard",
	tackSBPT:   "tack starboard to port",
	gybePTSB:   "gybe port to starboard",
	gybeSBPT:   "gybe starboard to port",
	drifting:   "drifting",
}

var posToWindAttitude = map[pointOfSail]windAttitude{
	closeSB: upwind,
	irons:   upwind,
	closePT: upwind,
	beamPT:  beam,
	broadPT: downwind,
	run:     downwind,
	broadSB: downwind,
	beamSB:  beam,
}

var posToTack = map[pointOfSail]tack{
	irons:   unknown,
	closePT: port,
	beamPT:  port,
	broadPT: port,
	run:     unknown,
	broadSB: starboard,
	beamSB:  starboard,
	closeSB: starboard,
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

func (windDirection direction) pointOfSail(heading int) *SegmentType {
	// adjust heading for wind direction
	adjusted := heading - int(windDirection)
	if adjusted < 0 {
		adjusted = 360 + adjusted
	}
	idx := int(math.Floor((float64(adjusted) + 22.5) / 45))
	if idx > 7 {
		idx = 0
	}
	pos := []pointOfSail{irons, closePT, beamPT, broadPT, run, broadSB, beamSB, closeSB}[idx]
	return &SegmentType{windDirection: windDirection, pointOfSail: pos}
}

func (windDirection direction) turnType(from, to *SegmentType) *SegmentType {
	return &SegmentType{windDirection: windDirection, turn: turnType(from.pointOfSail, to.pointOfSail)}
}

func (pos pointOfSail) String() string {
	return posToString[pos]
}

func (pos pointOfSail) windAttitude() windAttitude {
	return posToWindAttitude[pos]
}

func (pos pointOfSail) tack() tack {
	return posToTack[pos]
}

type turns []turn

func (ts turns) set(bitValue int, ambiguous bool) turns {
	var ts2 turns
	if ambiguous {
		ts2 = append(ts2, ts...)
	}
	for _, t := range ts {
		ts2 = append(ts2, t+turn(bitValue))
	}
	return ts2
}

func turnType(from, to pointOfSail) turn {
	toWA := to.windAttitude()
	toTack := to.tack()
	fromWA := from.windAttitude()
	fromTack := from.tack()
	// When tack or attitude is uncertain, we add both cases as candidates
	var candidates turns = turns{0}
	if toTack == starboard {
		candidates = candidates.set(8, false)
	} else if toTack == unknown {
		candidates = candidates.set(8, true)
	}
	if toWA == upwind {
		candidates = candidates.set(4, false)
	} else if toWA == beam {
		candidates = candidates.set(4, true)
	}
	if fromTack == starboard {
		candidates = candidates.set(2, false)
	} else if fromTack == unknown {
		candidates = candidates.set(2, true)
	}
	if fromWA == upwind {
		candidates = candidates.set(1, false)
	} else if fromWA == beam {
		candidates = candidates.set(1, true)
	}
	// See if any of the candidates are valid,
	// return first one that is
	for _, t := range candidates {
		if _, found := turnToString[t]; found {
			return t
		}
	}
	// Otherwise return unknown
	return 15
}

func (t turn) windAttitude() windAttitude {
	if t&5 == 0 {
		return downwind
	} else if t&5 == 5 {
		return upwind
	} else {
		return beam
	}
}

func (t turn) String() string {
	if t, found := turnToString[t]; !found {
		return "unknown"
	} else {
		return t
	}
}

type SegmentType struct {
	turn
	pointOfSail
	windDirection direction
}

func (st *SegmentType) String() string {
	if st.turn == 0 {
		return st.pointOfSail.String() + " " + st.windDirection.String()
	} else {
		return st.turn.String() + " " + st.windDirection.String()
	}
}

func (st *SegmentType) windAttitude() windAttitude {
	if st.turn == 0 {
		return st.pointOfSail.windAttitude()
	} else {
		return st.turn.windAttitude()
	}
}
