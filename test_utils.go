package main

import (
	"testing"
	"time"

	"github.com/tkrajina/gpxgo/gpx"
)

func readTrackSample(t *testing.T, data string) *Track {
	g, err := gpx.ParseString(`<?xml version="1.0" encoding="UTF-8"?>
	<gpx xmlns="http://www.topografix.com/GPX/1/1" version="1.1">
	<trk><trkseg>` + data + `</trkseg></trk></gpx>`)
	if err != nil {
		t.Error(err)
	}
	if g == nil {
		t.Error("failed to parse track sample")
	}
	ss := gpxGetSegments(g, "")
	ts := gpxBuildTracks(ss, time.Hour)
	if len(ts) != 1 {
		t.Errorf("found %d tracks", len(ts))
	}
	return &ts[0]
}

func logTrackPoints(t *testing.T, trk *Track) {
	for _, s := range trk.Segments {
		t.Log(s.String())
		for _, p := range s.Points {
			t.Log(p.String())
		}
	}
}

func assertEqual[T comparable](t *testing.T, got T, exp T) {
	t.Helper()
	if got == exp {
		return
	}
	t.Errorf("got: %v\nexp: %v", got, exp)
}
