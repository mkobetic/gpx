package main

import (
	"fmt"
	"testing"
)

func Test_Direction(t *testing.T) {
	for i, tt := range []struct {
		heading   int
		direction direction
	}{
		{0, N},
		{11, N},
		{12, NNE},
		{33, NNE},
		{34, NE},
		{56, NE},
		{57, ENE},
		{78, ENE},
		{79, E},
		{101, E},
		{102, ESE},
		{123, ESE},
		{124, SE},
		{146, SE},
		{147, SSE},
		{168, SSE},
		{169, S},
		{191, S},
		{192, SSW},
		{213, SSW},
		{214, SW},
		{236, SW},
		{237, WSW},
		{258, WSW},
		{259, W},
		{281, W},
		{282, WNW},
		{303, WNW},
		{304, NW},
		{326, NW},
		{327, NNW},
		{348, NNW},
		{349, N},
		{359, N},
	} {
		t.Run(fmt.Sprintf("%d: %d %s", i, tt.heading, tt.direction.String()), func(t *testing.T) {
			dir := Direction(tt.heading)
			if dir != tt.direction {
				t.Errorf("exp: %s got: %s", tt.direction.String(), dir.String())
			}
		})
	}
}

func Test_PointOfSail(t *testing.T) {
	for i, tt := range []struct {
		wind    direction
		heading int
		pos     pointOfSail
	}{
		{N, 45, closePT},
		{N, 90, beamPT},
		{N, 135, broadPT},
		{N, 180, run},
		{N, 215, broadSB},
		{N, 270, beamSB},
		{N, 315, closeSB},

		{NW, 90, broadPT},
		{NW, 180, broadSB},
		{S, 90, beamSB},
		{S, 270, beamPT},
	} {
		t.Run(fmt.Sprintf("%d: %s %d", i, tt.wind.String(), tt.heading), func(t *testing.T) {
			pos := tt.wind.pointOfSail(tt.heading)
			if pos.pointOfSail != tt.pos {
				t.Errorf("exp %s got %s", tt.pos.String(), pos.String())
			}
		})
	}
}

func Test_TurnType(t *testing.T) {
	for i, tt := range []struct {
		from, to pointOfSail
		turn     turn
	}{
		{closePT, closeSB, tackPTSB},
		{closeSB, closePT, tackSBPT},
		{closeSB, broadSB, bearawaySB},
		{closePT, broadPT, bearawayPT},
		{broadPT, broadSB, gybePTSB},
		{broadSB, broadPT, gybeSBPT},
		{broadSB, closeSB, roundupSB},
		{broadPT, closePT, roundupPT},
	} {
		t.Run(fmt.Sprintf("%d: %s-%s", i, tt.from, tt.to), func(t *testing.T) {
			got := turnType(tt.from, tt.to)
			if got != tt.turn {
				t.Errorf("exp %s got %s", tt.turn, got)
			}
		})
	}
}

func Test_WindDirectionTurnType(t *testing.T) {
	for i, tt := range []struct {
		wind     direction
		from, to int
		turn     string
	}{
		{N, 310, 40, "tack starboard to port N"},
	} {
		t.Run(fmt.Sprintf("%d: %d->%d@%s", i, tt.from, tt.to, tt.wind), func(t *testing.T) {
			got := tt.wind.turnType(tt.wind.pointOfSail(tt.from), tt.wind.pointOfSail(tt.to))
			if got.String() != tt.turn {
				t.Errorf("exp %s got %s", tt.turn, got)
			}
		})
	}
}
