![tests](https://github.com/mkobetic/gpx/actions/workflows/test.yaml/badge.svg)
![downloads](https://img.shields.io/github/downloads/mkobetic/gpx/total.svg)


Simple GPS track processor for sailing race tracks.

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

![sample track](https://github.com/user-attachments/assets/fbaca31c-839c-43a1-a1fe-6c5ff18f4e89)

Hovering over the track highlights the segment under the cursor and shows metrics at that point.
1. First line shows stats for the point under cursor
2. Second line shows stats for the segment under cursor
3. Third line shows the segment type determined by analysis

```
time: distance @ speed ↑ heading = track distance
length/duration @ min/avg/max speed ↑ min/max heading static|moving|turning
segment analysis given the determined wind direction
```

The map can be zoomed with mouse wheel scroll or touchpad pinch. The map can be panned with mouse or touchpad drag.

## usage

The repository is set up to automatically compile and release binaries for common desktop platforms (linux/mac/windows).
Download .tgz archive suitable for your platform from the latest release here https://github.com/mkobetic/gpx/releases.
The archive contains single file that is the compiled binary. It doesn't need anything, there's no installation process,
just put it somewhere on your $PATH or into a directory where you intend to use it. If you don't want it anymore, just delete the file.

To get the latest usage information, run the binary with the -h option, e.g.

```
> gpx -h
Usage: gpx [flags] files...

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

Flags:
  -a value
        analyze tracks using specified activity type
        supported types: sail
  -o string
        directory for generated files (default ".")
  -ss int
        discard segments that are shorter than this number of points (default 20)
  -v    verbose, print more processing details
  -version
        print version information
  -vo value
        video time offset for subtitles or chapters file, e.g -3.5m or 5m22s
        positive offset means video starts ahead of the track
        requires -a
  -wd value
        wind direction to use for analyzing the track, e.g. NE, or SSW
        if UNK then deduce direction from the track
        implies -a sail

```

## samples

`samples/out` directory contains files generated from the input files in `samples/in` directory.

The output files can be regenerated as follows:
```
$ rm samples/out/*; make && ./gpx -o samples/out -a sail -vo 0s samples/in/*
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

## track analysis

If the -a option is used the chosen activity type is used to analyse the tracks and split them into relatively "straight" moving, turning and static segments. The analysis is performed using parameters associated with the selected activity type. Currently the only supported activity is `sail` which is suitable for sail racing GPS tracks. Parameters for other activity types can be added (create an issue describing what you would like to see).

If the -wd (wind direction, e.g -wd NW) option is used, the track segments are further classified based on the provided wind direction. Moving segments are assigned their corresponding point of sail and tack, turning segments are assigned their turn type (tack, gybe, round up, bear away) and tack.

If the wind direction is specified as UNK (unknown), it will be determined by analyzing the moving segments of the track. If the determination fails a warning will be printed and the point of sail analysis will be skipped.


## gps video subtitles

It is nice to be able to overlay GPS information over the video that you may have recorded on your boat. There are many guides out there showing how to use video editors to render cute measurement gauges into your video recording. It can look pretty good but is very manual and time consuming.

An alternative approach that isn't as artistic but is simple and fast, because you don't have to re-encode the whole video, is to use subtitles to show the GPS information. It's text only, no graphics but can convey most of the same information. Adding subtitles to a video is nearly instant compared to rerendering the full video.

This tool can spit out a subtitles file for your video with following GPS stats:

```
time: distance @ speed ↑ heading = total distance
```

It will look something like this

```
17:51:43: 26.8 m @ 5.8 kts ↑ 33° NNE = 0.07 nm
```

The subtitle granularity matches the granularity of the gps track, i.e. new subtitle for each track point.

The subtitle file generation is gated by the `-vo` flag that requires a "video offset" as its argument. This is because the start of the gps track more than likely doesn't align with the start of the video. The offset specifies how much they are off. Positive offset means the video starts before the gps track, negative offset means the video starts after the gps track. If miraculously they align perfectly set offset to 0. The format of the offset argument is documented here https://pkg.go.dev/time#ParseDuration, e.g. `-3.5m` or `10m44s`.

Note that you may need to enable subtitles in your video player to have them show up. Here's a screenshot of the subtitles shown in a video

![Screenshot 2025-01-06 at 22 59 10](https://github.com/user-attachments/assets/da055ae4-da4f-4a53-9424-6002aa509368)


## segment video chapters

The `-vo` option also triggers generation of a video metadata file that defines a chapter for each segment of the track. Chapters can be used to navigate to the corresponding section of the video in players that support this facility. Same as with subtitles the precision of the chapter definitions depends on correctly determined video offset.

Here's a screenshot of MacOS video player with the chapter list open

![Screenshot 2025-04-10 at 12 50 52](https://github.com/user-attachments/assets/f812a961-dfa6-4eab-89c7-957834a1e22a)

### figuring out video offset

#### when camera and gps watch time is the same

If you can make sure that the clocks on the camera and on the gps device are in sync, you can use `ffmpeg` to lift the `creation_time` off the video as follows:

```
ffprobe -show_entries format_tags=creation_time -of csv=print_section=0 <the-video-file>
```

That should give you something like `2024-08-24T15:08:29.000000Z`. Note that if the timezone is set to Z it's quite likely local time, not UTC.

Then check the timestamp in the beginning of your gpx file, it could look something like this

```
  ...
  <metadata>
    <link href="connect.garmin.com">
      <text>Garmin Connect</text>
    </link>
    <time>2024-08-24T19:01:05.000Z</time>
  </metadata>
``` 

This one likely is in UTC, so you'll need to convert to local time. Then simply subtract the video timestamp from the gpx timestamp to get the video offset.

#### when camera and gps watch time is not the same

Unless you prepared well and made sure the camera and gps device clocks are in sync, then they most likely are not. In that case you need to figure out how much they are off and add that difference to the video offset as well.

You may be able to check the time on both after the fact and figure out the difference that way.

Another option, if you think about this ahead of time is to just record your GPS watch on the video when you start recording, that way you'll have a fairly good idea of what the watch time was at a given point in the video, it should be easy to calculate the offset from that. You just have to make sure you have a reasonably clear and sharp picture of your watch. I find that I can bring my watch about 6 inches from my old GoPro session and hold it there for a second.


## adding subtitles and chapters to a video file (using ffmpeg)

Let's assume you have video.mp4 and the subtitles.vtt and video.chapters files that you successfully produced with this tool. Here's how you can add the subtitles and chapters to the video.

```
ffmpeg -i video.mp4 -i subtitles.vtt -i video.chapters -map 0:v -map 0:a -map 1:s -map_metadata 2 -c copy -c:s mov_text -metadata:s:s:0 language=eng -y video-with-subtitles.mp4
```
Here's what the bits of the command mean:
```
-i video.mp4    = use this video file (input number 0)
-i subtitles.vtt = use this subtitle file (input number 1)
-i subtitles.vtt = use this subtitle file (input number 2)
-map 0:v        = take the video stream from the video file (input 0)
-map 0:a        = take the audio stream from the video file (input 0)
-map 1:s        = take the subtitle stream from the subtitle file (input 1)
-map_metadata 2 = take the metadata from the chapter file (input 2)
-c copy         = don't process the video/audio just copy it over as is
-c:s mov_text   = encode the subtitle stream in a way that works with mp4 container
-metadata:s:s:0 language=eng = set the subtitles language to english (that's how it will show in subtitle options in the player)
-y              = just (over)-write the final file don't ask for confirmation
video-with-subtitles.mp4 = name of the video file with subtitles
```
This should produce a new video file with the same content as the original file with the subtitle stream and chapter metadata added into it. It should be fast because the video content is just copied over as is without re-processing.


## References
* https://gist.github.com/spirillen/af307651c4261383a6d651038a82565d
* https://grep.be/blog/en/computer/play/Adding_subtitles_with_FFmpeg/
* https://www.bannerbear.com/blog/how-to-add-subtitles-to-a-video-file-using-ffmpeg/
* https://developer.mozilla.org/en-US/docs/Web/API/WebVTT_API/Web_Video_Text_Tracks_Format
* https://ikyle.me/blog/2020/add-mp4-chapters-ffmpeg
* https://ffmpeg.org/ffmpeg-formats.html#Metadata-2
