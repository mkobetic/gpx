Simple GPS track processor for sailing race tracks.

* reads all gpx files specified on the command line
* pulls out all segments
* discards any duplicate or superfluous (very short) segments
* combines segments that are no more than 1h apart into a track
* renders each track into an svg file
* saves each track into a new gpx file

`samples/out` directory contains files generated from the input files in `samples/in` directory.

The output files can be regenerated as follows:
```
$ rm samples/out/*; make && ./gpx -o samples/out samples/in/*
Dropped 32 duplicate and bogus segments
16-05-25 17:56:07 16.09nm 01.32nm x 00.86nm (2h21m36s)
16-06-01 18:57:23 06.76nm 01.08nm x 00.57nm (1h11m2s)
16-06-04 11:03:46 04.89nm 01.25nm x 00.52nm (1h23m41s)
16-06-05 09:56:20 23.06nm 01.71nm x 01.98nm (3h1m4s)
16-06-22 17:51:25 18.56nm 01.42nm x 00.71nm (2h29m44s)
16-07-26 18:48:29 04.04nm 01.41nm x 00.76nm (42m51s)
16-08-08 18:29:00 01.61nm 00.73nm x 00.17nm (13m45s)
16-08-17 18:20:50 04.37nm 01.28nm x 00.58nm (59m19s)
16-08-20 17:34:32 10.07nm 01.70nm x 01.98nm (1h12m42s)
16-08-24 17:50:55 14.82nm 01.31nm x 00.71nm (2h4m15s)
```
