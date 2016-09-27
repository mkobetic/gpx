Simple GPS track processor for sailing race tracks.

* reads all gpx files specified on the command line
* pulls out all segments
* discards any duplicate or degenerate (very short) segments
* combines segments that are no more than 1h apart into a track
* renders each track into an svg file
* saves each track into a new gpx file

The `samples/out` directory contains files generated from the `samples/in` directory.
The output files can be reproduced with `./gpx -o samples/out samples/in/*`
