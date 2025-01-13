package main

import (
	"fmt"
	"testing"

	"github.com/tkrajina/gpxgo/gpx"
)

func Test_Heading(t *testing.T) {
	bounds := gpx.GpxBounds{MinLatitude: 44, MaxLatitude: 45, MinLongitude: -78, MaxLongitude: -77}
	m := NewMap(bounds, 1000)
	for i, tt := range []struct {
		p1lon, p1lat, p2lon, p2lat float64
		heading                    int // degrees
		direction                  string
		distance                   int // meters
	}{
		{-77.5, 44.5, -77.4, 44.6, 35, "NE", 13658},
		{-77.5, 44.5, -77.4, 44.4, 145, "SE", 13658},
		{-77.5, 44.5, -77.6, 44.4, 215, "SW", 13658},
		{-77.5, 44.5, -77.6, 44.6, 325, "NW", 13658},
	} {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			h := m.Heading(point(tt.p1lon, tt.p1lat), point(tt.p2lon, tt.p2lat))
			if h != tt.heading {
				t.Errorf("exp: %d, got %d", tt.heading, h)
			}
			d := int(m.Distance(point(tt.p1lon, tt.p1lat), point(tt.p2lon, tt.p2lat), meter))
			if d != tt.distance {
				t.Errorf("exp: %d, got %d", tt.distance, d)
			}
		})
	}
}

func point(lon, lat float64) *gpx.GPXPoint {
	return &gpx.GPXPoint{Point: gpx.Point{
		Longitude: lon,
		Latitude:  lat,
	}}
}
