// Simple GPS track processor for sailing race tracks (.gpx)
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/tkrajina/gpxgo/gpx"
)

var (
	// Build Parameters
	Built     string // time built in UTC
	Commit    string // source commit SHA
	Branch    string // source branch
	GoVersion string // Go version used to build
)

const usage = `Usage: gpx [flags] files...

Simple GPS track processor for sailing race tracks:
* reads all gpx files specified on the command line
* pulls out all track segments
* discards any duplicate or superfluous (very short) segments
* combines segments that are no more than 1h apart into a track
* renders each track into a map rendered as an svg file
* saves each track into a new gpx file
* (optional) analyses each track and splits it into straight moving, turning and static segments
* (optional) generates subtitles file with gps metrics that can be embedded in a video file
(see https://github.com/mkobetic/gpx/blob/master/README.md for more details)

Flags:`

func init() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintln(w, usage)
		flag.PrintDefaults()
	}
}

func main() {
	// flags
	out := flag.String("o", ".", "directory for generated files")
	var videoOffset *time.Duration
	flag.Func("vo", "video time offset for subtitles file, positive offset means video starts ahead of the track", func(ts string) error {
		d, err := time.ParseDuration(ts)
		if err != nil {
			return err
		}
		videoOffset = &d
		return nil
	})
	var windDirection *direction
	flag.Func("wd", "wind direction to use for analyzing the track, e.g. NE, or SSW", func(wd string) error {
		for d, s := range directionToString {
			if s == wd {
				windDirection = &d
				return nil
			}
		}
		return fmt.Errorf("invalid wind direction value %s", wd)
	})
	fMinSegmentLength := flag.Int("ss", 20, "discard segments that are shorter than this number of points")
	fVersion := flag.Bool("version", false, "print version information")
	fVerbose := flag.Bool("v", false, "verbose, print more processing details")
	flag.Parse()

	if *fVersion {
		fmt.Println("built on " + Built)
		fmt.Printf("built from %s@%s\n", Branch, Commit)
		fmt.Println("built with " + GoVersion)
		os.Exit(0)
	}

	// args
	if len(flag.Args()) == 0 {
		fmt.Println("Transforms specified gpx files into a gpx, svg and video subtitles files for individual race tracks")
		fmt.Println("Usage: gpx [-o <dir>] [-vo <duration>] <files>")
		flag.Usage()
		return
	}

	// Collect all the original segments from the parsed files.
	// Using Segment instead of gpx.GPXTrackSegment so that we can attach the filenames that they came from.
	var protoSegments Segments
	for _, fn := range flag.Args() {
		g, err := gpx.ParseFile(fn)
		if err != nil {
			fmt.Printf("Error opening %s: %s\n", fn, err)
			return
		}
		protoSegments = append(protoSegments, gpxGetSegments(g, filepath.Base(fn))...)
	}
	sort.Sort(protoSegments)
	sn := len(protoSegments)
	protoSegments = gpxDedupeSegments(protoSegments, *fMinSegmentLength)
	protoSegments = gpxSplitSegments(protoSegments, time.Hour)
	protoSegments = gpxDedupeSegments(protoSegments, *fMinSegmentLength)
	fmt.Printf("Dropped %d duplicate and short segments\n", sn-len(protoSegments))
	for _, t := range gpxBuildTracks(protoSegments, time.Hour) {
		if windDirection != nil {
			t.gpxAnalyze(Sailing)
			t.posClassify(*windDirection)
		}
		fmt.Println(t.String())
		if *fVerbose {
			for i, s := range t.Segments {
				fmt.Printf("%d: %s\n", i, s.String())
			}
		}
		if err := t.WriteMapFile(*out); err != nil {
			fmt.Println(err)
		}
		if videoOffset != nil {
			if err := t.WriteSubtitleFile(*out, *videoOffset); err != nil {
				fmt.Println(err)
			}
		}
		if err := t.WriteGpxFile(*out); err != nil {
			fmt.Println(err)
		}
	}
}
