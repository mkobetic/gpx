// Simple GPS track processor for sailing race tracks (.gpx)
package main

import (
	"flag"
	"fmt"
	"sort"
	"time"

	"github.com/tkrajina/gpxgo/gpx"
)

func main() {
	out := flag.String("o", ".", "directory for generated files")
	flag.Parse()
	if len(flag.Args()) == 0 {
		fmt.Println("Transforms specified gpx files into a gpx and svg file for individual race tracks")
		fmt.Println("Usage: gpx [-o <dir>] <files>")
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
		for _, t := range g.Tracks {
			for i := range t.Segments {
				segments = append(segments, Segment{&t.Segments[i]})
			}
		}
	}
	sort.Sort(segments)
	sn := len(segments)
	segments = segments.Dedupe()
	fmt.Printf("Dropped %d duplicate and bogus segments\n", sn-len(segments))
	for _, t := range segments.Tracks(time.Hour) {
		fmt.Println(t.String())
		if err := t.WriteMapFile(*out); err != nil {
			fmt.Println(err)
		}
		if err := t.WriteGpxFile(*out); err != nil {
			fmt.Println(err)
		}
	}
}
