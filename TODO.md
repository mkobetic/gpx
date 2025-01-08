
# SVG map
* split track to sections of some reasonable predefined length, hover-highlight only those, not whole gpx segments
* add total stats (max speed, etc)
* render directional arrows onto the track
* add speed timeline chart
* hover on timeline should show where on the track it is and vice versa
* allow selecting a timerange on the timeline and show only or highlight that part of the track
* split track into tack, gybe, leg sections
  * can we figure out prevailing wind direction?
  * maybe allow providing prevailing wind direction as input (e.g SW)
* highlight tacks, gybes on the timeline
* fetch and add satellite or chart tile as background
* maybe add playback, little boat running along the track with the stats subtitles or something, different playback speeds

# GPX files
* add stats to the track description

# OTHER
* add GH actions to build binary releases linux, mac, windows
* emit chapter metadata for video files 
  (e.g. each tack starts a new chapter)
  https://ffmpeg.org/ffmpeg-all.html#Metadata-2