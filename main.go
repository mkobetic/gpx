// Simple GPS track processor for sailing race tracks (.gpx)
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
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
* combines segments that are no more than 1h apart into tracks
* renders each track into a map saved as an SVG file
* saves each track into a new GPX file
* (optional) analyses each track and splits it into straight moving, turning and static segments (-a)
* (optional) point of sail analysis and classification of the segments (-wd)
* (optional) generates subtitles file with gps metrics that can be embedded in a video file (-vo)
* (optional) generates chapter metadata file with a chapter for each segment that can be embedded in a video file (-vo)
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
	var usage string
	out := flag.String("o", ".", "directory for generated files")
	fMinSegmentLength := flag.Int("ss", 20, "discard segments that are shorter than this number of points")
	fVersion := flag.Bool("version", false, "print version information")
	fVerbose := flag.Bool("v", false, "verbose, print more processing details")

	var fActivity Activity
	usage = "analyze tracks using specified activity type\nsupported types: " + strings.Join(KnownActivities, ", ")
	flag.Func("a", usage, func(at string) error {
		for a, params := range Activities {
			if a == at {
				fActivity = params
				return nil
			}
		}
		return fmt.Errorf("%s is not a recognized activity type\nknown activities are "+strings.Join(KnownActivities, ", "), at)
	})

	var fWindDirection *direction
	usage = "wind direction to use for analyzing the track, e.g. NE, or SSW\nif UNK then deduce direction from the track\nimplies -a sail"
	flag.Func("wd", usage, func(wd string) error {
		for d, s := range directionToString {
			if s == wd {
				fActivity = Sailing
				fWindDirection = &d
				return nil
			}
		}
		return fmt.Errorf("%s is not a recognized wind direction\nvalid values are "+strings.Join(WindDirections, ", "), wd)
	})

	var fVideoOffset *time.Duration
	usage = "video time offset for subtitles or chapters file, e.g -3.5m or 5m22s\npositive offset means video starts ahead of the track\nrequires -a"
	flag.Func("vo", usage, func(ts string) error {
		if fActivity == nil {
			return fmt.Errorf("option -vo requires analysis (option -a)")
		}
		d, err := time.ParseDuration(ts)
		if err != nil {
			return fmt.Errorf(err.Error() + "\nformat of the offset value is documented at https://pkg.go.dev/time#ParseDuration")
		}
		fVideoOffset = &d
		return nil
	})

	flag.Parse()

	if *fVersion {
		fmt.Println("built on " + Built)
		fmt.Printf("built from %s@%s\n", Branch, Commit)
		fmt.Println("built with " + GoVersion)
		os.Exit(0)
	}

	// args
	if len(flag.Args()) == 0 {
		fmt.Println("Transforms specified gpx files into a gpx, svg and video subtitle and chapter files for individual race tracks.")
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

	// Reassemble tracks from gathered proto-segments and process them.
	for _, t := range gpxBuildTracks(protoSegments, time.Hour) {
		if fActivity != nil {
			t.gpxAnalyze(Sailing)
			if fWindDirection != nil {
				windDirection := *fWindDirection
				if windDirection == UNK {
					windDirection = t.windDirection()
				}
				if windDirection == UNK {
					fmt.Printf("%s\n  WARNING: Could not determine wind direction, skipping point of sail analysis\n", t.String())
				} else {
					t.posClassify(windDirection)
				}
			}
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
		if fActivity != nil && fVideoOffset != nil {
			if err := t.WriteSubtitleFile(*out, *fVideoOffset); err != nil {
				fmt.Println(err)
			}
			if err := t.WriteChapterFile(*out, *fVideoOffset); err != nil {
				fmt.Println(err)
			}
		}
		if err := t.WriteGpxFile(*out); err != nil {
			fmt.Println(err)
		}
	}
}
