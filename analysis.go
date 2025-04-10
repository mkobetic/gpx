package main

import (
	"fmt"
	"math"
)

// units for Distance and Speed functions,
// expressed as the length of one degree of longitude at the equator
type unit float64

const EquatorialRadius = 6378
const km unit = 2 * math.Pi * EquatorialRadius / 360
const meter unit = 1000 * km
const nm unit = 60

func (u unit) distance() string {
	switch u {
	case km:
		return "km"
	case meter:
		return "m"
	case nm:
		return "nm"
	default:
		panic(fmt.Errorf("unknown unit %f", u))
	}
}

func (u unit) convertDistance(d float64, from unit) float64 {
	return d * float64(u) / float64(from)
}

func (u unit) speed() string {
	switch u {
	case km:
		return "km/h"
	case meter:
		return "m/s"
	case nm:
		return "kts" // nm/h
	default:
		panic(fmt.Errorf("unknown unit %f", u))
	}
}

// NOTE: all distance values are in distanceUnits!
// longDistanceUnit can be used to convert the distance values to larger units if desired, e.g. full track length
type AnalysisParameters struct {
	distanceUnit     unit    // what is point distance measured in
	longDistanceUnit unit    // unit to use for longer distance (e.g. track distance)
	speedUnit        unit    // what is speed measured in (in distance units, time is implied m/s, km/h, nm/h=kts)
	lookAround       float64 // how far back and ahead to look when analyzing a point (in distanceUnits)
	movingSpeed      float64 // what's the minimum speed to be considered as moving as opposed to stationary (in speedUnits)
	turningChange    int     // what's the minimum heading change to consider the point to be part of a turn (in degrees)
}

func (params *AnalysisParameters) asLongDistance(dist float64) float64 {
	return params.longDistanceUnit.convertDistance(dist, params.distanceUnit)
}

func (params *AnalysisParameters) distance() string {
	return params.distanceUnit.distance()
}

func (params *AnalysisParameters) longDistance() string {
	return params.longDistanceUnit.distance()
}

func (params *AnalysisParameters) speed() string {
	return params.speedUnit.speed()
}

type Activity *AnalysisParameters

var Sailing Activity = &AnalysisParameters{
	distanceUnit:     meter,
	longDistanceUnit: nm,
	speedUnit:        nm,
	lookAround:       50, // m
	movingSpeed:      1,  // kts
	turningChange:    60, // degrees
}

var Activities = map[string]Activity{
	"sail": Sailing,
}

var KnownActivities []string

func init() {
	for a := range Activities {
		KnownActivities = append(KnownActivities, a)
	}
}
