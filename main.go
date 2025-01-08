// Simple GPS track processor for sailing race tracks (.gpx)
package main

import (
	"flag"
	"fmt"
	"path/filepath"
	"sort"
	"time"

	"github.com/tkrajina/gpxgo/gpx"
)

func main() {
	out := flag.String("o", ".", "directory for generated files")
	var offset *time.Duration
	flag.Func("vo", "video time offset for subtitles file, positive offset means video starts ahead of the track", func(ts string) error {
		d, err := time.ParseDuration(ts)
		if err != nil {
			return err
		}
		offset = &d
		return nil
	})
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("Transforms specified gpx files into a gpx, svg and video subtitles files for individual race tracks")
		fmt.Println("Usage: gpx [-o <dir>] [-vo <duration>] <files>")
		flag.Usage()
		return
	}
	var segments Segments
	for _, fn := range flag.Args() {
		g, err := gpx.ParseFile(fn)
		if err != nil {
			fmt.Printf("Error opening %s: %s\n", fn, err)
			return
		}
		segments = append(segments, GetSegments(g, filepath.Base(fn))...)
	}
	sort.Sort(segments)
	sn := len(segments)
	segments = segments.Dedupe(20).Split(time.Hour).Dedupe(20)
	fmt.Printf("Dropped %d duplicate and bogus segments\n", sn-len(segments))
	for _, t := range segments.Tracks(time.Hour) {
		fmt.Println(t.String())
		if err := t.WriteMapFile(*out); err != nil {
			fmt.Println(err)
		}
		if offset != nil {
			if err := t.WriteSubtitleFile(*out, *offset); err != nil {
				fmt.Println(err)
			}
		}
		if err := t.WriteGpxFile(*out); err != nil {
			fmt.Println(err)
		}
	}
}
