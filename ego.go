package main
import (
"fmt"
"io"
"github.com/tkrajina/gpxgo/gpx"
)
//line map2.ego:1
 func (m *Map) renderLines(w io.Writer, t *Track) error  {
//line map2.ego:2
_, _ = fmt.Fprintf(w, "\n")
//line map2.ego:3
_, _ = fmt.Fprintf(w, "\n")
//line map2.ego:4
_, _ = fmt.Fprintf(w, "\n<svg width=\"")
//line map2.ego:4
_, _ = fmt.Fprintf(w, "%v",  m.w )
//line map2.ego:4
_, _ = fmt.Fprintf(w, "\" height=\"")
//line map2.ego:4
_, _ = fmt.Fprintf(w, "%v",  m.h )
//line map2.ego:4
_, _ = fmt.Fprintf(w, "\" version=\"1.1\" xmlns=\"http://www.w3.org/2000/svg\">\n    <style type=\"text/css\" >\n        <![CDATA[\n            .segment { fill: none; stroke-width: 4 }\n            .segment:hover { stroke-width: 8 }\n        ]]>\n    </style>\n    <g id=\"legend\">\n        ")
//line map2.ego:12
 for i := range palette { 
//line map2.ego:13
_, _ = fmt.Fprintf(w, "\n        <rect x=\"")
//line map2.ego:13
_, _ = fmt.Fprintf(w, "%v",  30*i )
//line map2.ego:13
_, _ = fmt.Fprintf(w, "\" y=\"0\" width=\"30\" height=\"20\" fill=\"")
//line map2.ego:13
_, _ = fmt.Fprintf(w, "%v",  fmt.Sprintf("#%03x",palette[i]) )
//line map2.ego:13
_, _ = fmt.Fprintf(w, "\"/>\n        <rect x=\"")
//line map2.ego:14
_, _ = fmt.Fprintf(w, "%v",  30*i )
//line map2.ego:14
_, _ = fmt.Fprintf(w, "\" y=\"")
//line map2.ego:14
_, _ = fmt.Fprintf(w, "%v",  m.h-20 )
//line map2.ego:14
_, _ = fmt.Fprintf(w, "\" width=\"30\" height=\"20\" fill=\"")
//line map2.ego:14
_, _ = fmt.Fprintf(w, "%v",  fmt.Sprintf("#%03x",palette[i]) )
//line map2.ego:14
_, _ = fmt.Fprintf(w, "\"/>\n        ")
//line map2.ego:15
 } 
//line map2.ego:16
_, _ = fmt.Fprintf(w, "\n        ")
//line map2.ego:16
 for i := 0; i < len(palette); i += 5 { 
//line map2.ego:17
_, _ = fmt.Fprintf(w, "\n        <text x=\"")
//line map2.ego:17
_, _ = fmt.Fprintf(w, "%v",  30*i+5 )
//line map2.ego:17
_, _ = fmt.Fprintf(w, "\" y=\"16\" fill=\"white\">")
//line map2.ego:17
_, _ = fmt.Fprintf(w, "%v",  fmt.Sprint(i) )
//line map2.ego:17
_, _ = fmt.Fprintf(w, "kts</text>\n        <text x=\"")
//line map2.ego:18
_, _ = fmt.Fprintf(w, "%v",  30*i+5 )
//line map2.ego:18
_, _ = fmt.Fprintf(w, "\" y=\"")
//line map2.ego:18
_, _ = fmt.Fprintf(w, "%v",  m.h-4 )
//line map2.ego:18
_, _ = fmt.Fprintf(w, "\" fill=\"white\">")
//line map2.ego:18
_, _ = fmt.Fprintf(w, "%v",  fmt.Sprint(i) )
//line map2.ego:18
_, _ = fmt.Fprintf(w, "kts</text>\n        ")
//line map2.ego:19
 } 
//line map2.ego:20
_, _ = fmt.Fprintf(w, "\n    </g>\n    ")
//line map2.ego:21
 for i := range t.Segments { 
//line map2.ego:22
_, _ = fmt.Fprintf(w, "\n    <g class=\"segment\">\n        ")
//line map2.ego:23
 t.Segment(i).EachPair(func(prev, next *gpx.GPXPoint) { 
//line map2.ego:24
_, _ = fmt.Fprintf(w, "\n        ")
//line map2.ego:24
 x1, y1 := m.Point(prev); x2, y2 := m.Point(next); c := m.SpeedColor(prev, next) 
//line map2.ego:25
_, _ = fmt.Fprintf(w, "\n        <!-- ")
//line map2.ego:25
_, _ = fmt.Fprintf(w, "%v",  m.Distance(prev,next,nm) )
//line map2.ego:25
_, _ = fmt.Fprintf(w, " nm; ")
//line map2.ego:25
_, _ = fmt.Fprintf(w, "%v",  m.Speed(prev,next,nm) )
//line map2.ego:25
_, _ = fmt.Fprintf(w, " kts -->\n        <line class=\"step\" x1=\"")
//line map2.ego:26
_, _ = fmt.Fprintf(w, "%v",  x1 )
//line map2.ego:26
_, _ = fmt.Fprintf(w, "\" y1=\"")
//line map2.ego:26
_, _ = fmt.Fprintf(w, "%v",  y1 )
//line map2.ego:26
_, _ = fmt.Fprintf(w, "\" x2=\"")
//line map2.ego:26
_, _ = fmt.Fprintf(w, "%v",  x2 )
//line map2.ego:26
_, _ = fmt.Fprintf(w, "\" y2=\"")
//line map2.ego:26
_, _ = fmt.Fprintf(w, "%v",  y2 )
//line map2.ego:26
_, _ = fmt.Fprintf(w, "\" stroke=\"")
//line map2.ego:26
_, _ = fmt.Fprintf(w, "%v",  c )
//line map2.ego:26
_, _ = fmt.Fprintf(w, "\"/>\n\t    ")
//line map2.ego:27
 }) 
//line map2.ego:28
_, _ = fmt.Fprintf(w, "\n    </g>\n\t")
//line map2.ego:29
 } 
//line map2.ego:30
_, _ = fmt.Fprintf(w, "\n</svg>\n")
return nil
}
//line map.ego:1
 func (m *Map) renderPolylines(w io.Writer, t *Track) error  {
//line map.ego:2
_, _ = fmt.Fprintf(w, "\n<svg width=\"")
//line map.ego:2
_, _ = fmt.Fprintf(w, "%v",  m.w )
//line map.ego:2
_, _ = fmt.Fprintf(w, "\" height=\"")
//line map.ego:2
_, _ = fmt.Fprintf(w, "%v",  m.h )
//line map.ego:2
_, _ = fmt.Fprintf(w, "\" version=\"1.1\" xmlns=\"http://www.w3.org/2000/svg\">\n    <style type=\"text/css\" >\n        <![CDATA[\n            .segment { fill: none; stroke: blue; stroke-width: 3 }\n            .segment:hover { stroke: red }\n        ]]>\n    </style>\n    ")
//line map.ego:9
 for i := range t.Segments { 
//line map.ego:10
_, _ = fmt.Fprintf(w, "\n    <polyline class=\"segment\" points=\"")
//line map.ego:10
_, _ = fmt.Fprintf(w, "%v",  m.polylinePoints(t.Segment(i)) )
//line map.ego:10
_, _ = fmt.Fprintf(w, "\"/>\n\t")
//line map.ego:11
 } 
//line map.ego:12
_, _ = fmt.Fprintf(w, "\n</svg>\n")
return nil
}
