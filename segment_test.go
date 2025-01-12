package main

import (
	"testing"
	"time"

	"github.com/tkrajina/gpxgo/gpx"
)

var prostart = []byte(`
<?xml version="1.0" encoding="UTF-8"?>
<gpx xmlns="http://www.topografix.com/GPX/1/1" version="1.1" creator="https://github.com/tkrajina/gpxgo">
	<metadata>
		<author></author>
	</metadata>
	<trk>
		<trkseg>

			<trkpt lat="45.35330581665039" lon="-75.82582092285156">
			<time>2018-07-03T00:36:43Z</time>
			</trkpt>
			<trkpt lat="45.35330581665039" lon="-75.82582092285156">
			<time>2018-07-03T00:36:45Z</time>
			</trkpt>
			<trkpt lat="45.35330581665039" lon="-75.82582092285156">
			<time>2018-07-03T00:36:47Z</time>
			</trkpt>
			<trkpt lat="45.353302001953125" lon="-75.82582092285156">
			<time>2018-07-03T00:36:49Z</time>
			</trkpt>
			<trkpt lat="45.353302001953125" lon="-75.82582092285156">
			<time>2018-07-03T00:36:51Z</time>
			</trkpt>
			<trkpt lat="45.353302001953125" lon="-75.82582092285156">
			<time>2018-07-03T00:36:53Z</time>
			</trkpt>
			<trkpt lat="45.353302001953125" lon="-75.82582092285156">
			<time>2018-07-03T00:36:55Z</time>
			</trkpt>
			<trkpt lat="45.353302001953125" lon="-75.82582092285156">
			<time>2018-07-03T00:36:57Z</time>
			</trkpt>
			<trkpt lat="45.353302001953125" lon="-75.82582092285156">
			<time>2018-07-03T00:36:59Z</time>
			</trkpt>
			<trkpt lat="45.353302001953125" lon="-75.82582092285156">
			<time>2018-07-03T00:37:01Z</time>
			</trkpt>
			<trkpt lat="45.353302001953125" lon="-75.82582092285156">
			<time>2018-07-03T00:37:03Z</time>
			</trkpt>
			<trkpt lat="45.353302001953125" lon="-75.82582092285156">
			<time>2018-07-03T00:37:05Z</time>
			</trkpt>
			<trkpt lat="45.353302001953125" lon="-75.82582092285156">
			<time>2018-07-03T00:37:07Z</time>
			</trkpt>
			<trkpt lat="45.353302001953125" lon="-75.82582092285156">
			<time>2018-07-03T00:37:09Z</time>
			</trkpt>
			<trkpt lat="45.353302001953125" lon="-75.82582092285156">
			<time>2018-07-03T00:37:11Z</time>
			</trkpt>
			<trkpt lat="45.353302001953125" lon="-75.82582092285156">
			<time>2018-07-03T00:37:13Z</time>
			</trkpt>
			<trkpt lat="45.353302001953125" lon="-75.82582092285156">
			<time>2018-07-03T00:37:15Z</time>
			</trkpt>
			<trkpt lat="45.353302001953125" lon="-75.82582092285156">
			<time>2018-07-03T00:37:17Z</time>
			</trkpt>
			<trkpt lat="45.353370666503906" lon="-75.82583618164062">
			<time>2018-07-05T22:48:32Z</time>
			</trkpt>
			<trkpt lat="45.353370666503906" lon="-75.82583618164062">
			<time>2018-07-05T22:48:34Z</time>
			</trkpt>
			<trkpt lat="45.353370666503906" lon="-75.82583618164062">
			<time>2018-07-05T22:48:36Z</time>
			</trkpt>
			<trkpt lat="45.353370666503906" lon="-75.82583618164062">
			<time>2018-07-05T22:48:38Z</time>
			</trkpt>
			<trkpt lat="45.353370666503906" lon="-75.82583618164062">
			<time>2018-07-05T22:48:40Z</time>
			</trkpt>
			<trkpt lat="45.353370666503906" lon="-75.82583618164062">
			<time>2018-07-05T22:48:42Z</time>
			</trkpt>
			<trkpt lat="45.353370666503906" lon="-75.82583618164062">
			<time>2018-07-05T22:48:44Z</time>
			</trkpt>
			<trkpt lat="45.353370666503906" lon="-75.82583618164062">
			<time>2018-07-05T22:48:46Z</time>
			</trkpt>
			<trkpt lat="45.353370666503906" lon="-75.82583618164062">
			<time>2018-07-05T22:48:48Z</time>
			</trkpt>
			<trkpt lat="45.353370666503906" lon="-75.82583618164062">
			<time>2018-07-05T22:48:50Z</time>
			</trkpt>
			<trkpt lat="45.353370666503906" lon="-75.82583618164062">
			<time>2018-07-05T22:48:52Z</time>
			</trkpt>
			<trkpt lat="45.353370666503906" lon="-75.82583618164062">
			<time>2018-07-05T22:48:54Z</time>
			</trkpt>
			<trkpt lat="45.353370666503906" lon="-75.82583618164062">
			<time>2018-07-05T22:48:56Z</time>
			</trkpt>
			<trkpt lat="45.353370666503906" lon="-75.82583618164062">
			<time>2018-07-05T22:48:58Z</time>
			</trkpt>
			<trkpt lat="45.353370666503906" lon="-75.82583618164062">
			<time>2018-07-05T22:49:00Z</time>
			</trkpt>
			<trkpt lat="45.353370666503906" lon="-75.82583618164062">
			<time>2018-07-05T22:49:02Z</time>
			</trkpt>
			<trkpt lat="45.353370666503906" lon="-75.82583618164062">
			<time>2018-07-05T22:49:04Z</time>
			</trkpt>
			<trkpt lat="45.353370666503906" lon="-75.82583618164062">
			<time>2018-07-05T22:49:06Z</time>
			</trkpt>
			<trkpt lat="45.353370666503906" lon="-75.82583618164062">
			<time>2018-07-05T22:49:08Z</time>
			</trkpt>
			<trkpt lat="45.353370666503906" lon="-75.82583618164062">
			<time>2018-07-05T22:49:10Z</time>
			</trkpt>
			<trkpt lat="45.353370666503906" lon="-75.82583618164062">
			<time>2018-07-05T22:49:12Z</time>
			</trkpt>
		</trkseg>
	</trk>
</gpx>`)

func Test_Split(t *testing.T) {
	g, err := gpx.ParseBytes(prostart)
	if err != nil {
		t.Error(err)
	}
	ss := GetSegments(g, "")
	if len(ss) != 1 {
		t.Fail()
	}
	ss = ss.Split(time.Hour)
	if len(ss) != 2 {
		t.Error(ss)
	}
	t.Log("\n", ss)
	ts := ss.Tracks(time.Hour)
	if len(ts) != 2 {
		t.Error(ts)
	}
	t.Log("\n", ts)
}
