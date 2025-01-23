package main

type AnalysisParameters struct {
	distanceUnit  float64 // what is point distance measured in
	speedUnit     float64 // what is speed measured in (in distance units, time is implied m/s, km/h, nm/h=kts)
	lookAround    float64 // how far back and ahead to look when analyzing a point (in distanceUnits)
	movingSpeed   float64 // what's the minimum speed to be considered as moving as opposed to stationary (in speedUnits)
	turningChange int     // what's the minimum heading change to consider the point to be part of a turn (in degrees)
}

var Sailing = &AnalysisParameters{
	distanceUnit:  meter,
	speedUnit:     nm,
	lookAround:    50, // m
	movingSpeed:   1,  // kts
	turningChange: 60, // degrees
}
