package main

import (
	"testing"
)

func Test_AnalyzeTurn1(t *testing.T) {
	trk := readTrackSample(t, turn1)
	trk.gpxAnalyze(Sailing)
	assertEqual(t, len(trk.Segments), 2)
	assertEqual(t, trk.Segments[0].Mode, Turning)
	assertEqual(t, len(trk.Segments[0].Points), 15)
	assertEqual(t, trk.Segments[1].Mode, Moving)
	assertEqual(t, len(trk.Segments[1].Points), 8)
	assertEqual(t, trk.FileName(), "240824-0h01-00.2nm")
}

func Test_AnalyzeTurn2(t *testing.T) {
	trk := readTrackSample(t, turn2)
	trk.gpxAnalyze(Sailing)
	assertEqual(t, len(trk.Segments), 2)
	assertEqual(t, trk.Segments[0].Mode, Moving)
	assertEqual(t, len(trk.Segments[0].Points), 9)
	assertEqual(t, trk.Segments[1].Mode, Moving)
	assertEqual(t, len(trk.Segments[1].Points), 11)
	assertEqual(t, trk.FileName(), "240824-0h01-00.1nm")
	// logTrackPoints(t, trk)
}

func Test_AnalyzeTurn3(t *testing.T) {
	trk := readTrackSample(t, turn3)
	trk.gpxAnalyze(Sailing)
	assertEqual(t, len(trk.Segments), 3)
	assertEqual(t, trk.Segments[0].Mode, Moving)
	assertEqual(t, len(trk.Segments[0].Points), 5)
	assertEqual(t, trk.Segments[1].Mode, Turning)
	assertEqual(t, len(trk.Segments[1].Points), 7)
	assertEqual(t, trk.Segments[2].Mode, Moving)
	assertEqual(t, len(trk.Segments[2].Points), 11)
	assertEqual(t, trk.FileName(), "240824-0h01-00.1nm")
	// logTrackPoints(t, trk)
}

/*
* 1 static
* 2 moving
* 12 turning
* 8 moving
 */
const turn1 = `
			<trkpt lat="44.08929976634681" lon="-76.90646200440824">
				<ele>79</ele>
				<time>2024-08-24T19:09:56Z</time>
			</trkpt>
			<trkpt lat="44.08893154934049" lon="-76.90613636747003">
				<ele>80.19999694824219</ele>
				<time>2024-08-24T19:10:04Z</time>
			</trkpt>
			<trkpt lat="44.088761480525136" lon="-76.90594056621194">
				<ele>80.5999984741211</ele>
				<time>2024-08-24T19:10:08Z</time>
			</trkpt>
			<trkpt lat="44.0886424575001" lon="-76.9058041088283">
				<ele>81</ele>
				<time>2024-08-24T19:10:11Z</time>
			</trkpt>
			<trkpt lat="44.08860817551613" lon="-76.90575976856053">
				<ele>81.19999694824219</ele>
				<time>2024-08-24T19:10:12Z</time>
			</trkpt>
			<trkpt lat="44.08850843086839" lon="-76.90558014437556">
				<ele>81.19999694824219</ele>
				<time>2024-08-24T19:10:16Z</time>
			</trkpt>
			<trkpt lat="44.08847716636956" lon="-76.90544377081096">
				<ele>81</ele>
				<time>2024-08-24T19:10:20Z</time>
			</trkpt>
			<trkpt lat="44.08847423270345" lon="-76.90541669726372">
				<ele>80.5999984741211</ele>
				<time>2024-08-24T19:10:21Z</time>
			</trkpt>
			<trkpt lat="44.088479932397604" lon="-76.90537487156689">
				<ele>80.80000305175781</ele>
				<time>2024-08-24T19:10:23Z</time>
			</trkpt>
			<trkpt lat="44.08848856575787" lon="-76.90535953268409">
				<ele>81</ele>
				<time>2024-08-24T19:10:24Z</time>
			</trkpt>
			<trkpt lat="44.08851086162031" lon="-76.90534109249711">
				<ele>81.5999984741211</ele>
				<time>2024-08-24T19:10:26Z</time>
			</trkpt>
			<trkpt lat="44.088552352041006" lon="-76.90535022877157">
				<ele>81.5999984741211</ele>
				<time>2024-08-24T19:10:29Z</time>
			</trkpt>
			<trkpt lat="44.0885856281966" lon="-76.90537738613784">
				<ele>82</ele>
				<time>2024-08-24T19:10:31Z</time>
			</trkpt>
			<trkpt lat="44.088723342865705" lon="-76.90550747327507">
				<ele>80.19999694824219</ele>
				<time>2024-08-24T19:10:38Z</time>
			</trkpt>
			<trkpt lat="44.088884107768536" lon="-76.90567938610911">
				<ele>80.19999694824219</ele>
				<time>2024-08-24T19:10:45Z</time>
			</trkpt>
			<trkpt lat="44.08902509137988" lon="-76.90584467723966">
				<ele>79.80000305175781</ele>
				<time>2024-08-24T19:10:51Z</time>
			</trkpt>
			<trkpt lat="44.08923916518688" lon="-76.90609311684966">
				<ele>80.5999984741211</ele>
				<time>2024-08-24T19:10:59Z</time>
			</trkpt>
			<trkpt lat="44.089392805472016" lon="-76.90633988007903">
				<ele>81.19999694824219</ele>
				<time>2024-08-24T19:11:06Z</time>
			</trkpt>
			<trkpt lat="44.08942373469472" lon="-76.90641338936985">
				<ele>81.19999694824219</ele>
				<time>2024-08-24T19:11:08Z</time>
			</trkpt>
			<trkpt lat="44.0894815698266" lon="-76.9065678678453">
				<ele>81.5999984741211</ele>
				<time>2024-08-24T19:11:12Z</time>
			</trkpt>
			<trkpt lat="44.08961760811508" lon="-76.90686106681824">
				<ele>82.19999694824219</ele>
				<time>2024-08-24T19:11:19Z</time>
			</trkpt>
			<trkpt lat="44.08976278267801" lon="-76.90712484531105">
				<ele>82.19999694824219</ele>
				<time>2024-08-24T19:11:25Z</time>
			</trkpt>
			<trkpt lat="44.089788515120745" lon="-76.9071706943214">
				<ele>82.19999694824219</ele>
				<time>2024-08-24T19:11:26Z</time>
			</trkpt>
`

const turn2 = `
<trkpt lat="44.08867908641696" lon="-76.90315953455865">
<ele>95.5999984741211</ele>
<time>2024-08-24T20:48:49Z</time>
</trkpt>
<trkpt lat="44.088592836633325" lon="-76.90307093784213">
<ele>95.5999984741211</ele>
<time>2024-08-24T20:48:51Z</time>
</trkpt>
<trkpt lat="44.0884611569345" lon="-76.90293800085783">
<ele>95.5999984741211</ele>
<time>2024-08-24T20:48:54Z</time>
</trkpt>
<trkpt lat="44.08829024992883" lon="-76.90276122651994">
<ele>95.4000015258789</ele>
<time>2024-08-24T20:48:58Z</time>
</trkpt>
<trkpt lat="44.088209699839354" lon="-76.90267573110759">
<ele>95.4000015258789</ele>
<time>2024-08-24T20:49:00Z</time>
</trkpt>
<trkpt lat="44.08812252804637" lon="-76.90262384712696">
<ele>95</ele>
<time>2024-08-24T20:49:02Z</time>
</trkpt>
<trkpt lat="44.08807843923569" lon="-76.90260934643447">
<ele>95</ele>
<time>2024-08-24T20:49:03Z</time>
</trkpt>
<trkpt lat="44.088001661002636" lon="-76.90259702503681">
<ele>95.4000015258789</ele>
<time>2024-08-24T20:49:05Z</time>
</trkpt>
<trkpt lat="44.08792303875089" lon="-76.90260364674032">
<ele>95.5999984741211</ele>
<time>2024-08-24T20:49:08Z</time>
</trkpt>
<trkpt lat="44.08790795132518" lon="-76.90260842442513">
<ele>95.5999984741211</ele>
<time>2024-08-24T20:49:13Z</time>
</trkpt>
<trkpt lat="44.087941478937864" lon="-76.90258235670626">
<ele>96</ele>
<time>2024-08-24T20:49:23Z</time>
</trkpt>
<trkpt lat="44.087942065671086" lon="-76.90257925540209">
<ele>96</ele>
<time>2024-08-24T20:49:24Z</time>
</trkpt>
<trkpt lat="44.08796620555222" lon="-76.90254262648523">
<ele>95</ele>
<time>2024-08-24T20:49:36Z</time>
</trkpt>
<trkpt lat="44.08800878562033" lon="-76.90262971445918">
<ele>94.80000305175781</ele>
<time>2024-08-24T20:49:46Z</time>
</trkpt>
<trkpt lat="44.08803979866207" lon="-76.90269626677036">
<ele>95</ele>
<time>2024-08-24T20:49:50Z</time>
</trkpt>
<trkpt lat="44.08807139843702" lon="-76.90277966670692">
<ele>95.19999694824219</ele>
<time>2024-08-24T20:49:55Z</time>
</trkpt>
<trkpt lat="44.08813610672951" lon="-76.90294009633362">
<ele>94.80000305175781</ele>
<time>2024-08-24T20:50:04Z</time>
</trkpt>
<trkpt lat="44.088179022073746" lon="-76.90301301889122">
<ele>95</ele>
<time>2024-08-24T20:50:07Z</time>
</trkpt>
<trkpt lat="44.0882177464664" lon="-76.903075883165">
<ele>95</ele>
<time>2024-08-24T20:50:09Z</time>
</trkpt>
<trkpt lat="44.08831346780062" lon="-76.90323731862009">
<ele>95.19999694824219</ele>
<time>2024-08-24T20:50:14Z</time>
</trkpt>
`

const turn3 = `			<trkpt lat="44.10134145990014" lon="-76.922099115327">
				<ele>84</ele>
				<time>2024-08-24T19:17:00Z</time>
			</trkpt>
			<trkpt lat="44.101427625864744" lon="-76.92223599180579">
				<ele>83.5999984741211</ele>
				<time>2024-08-24T19:17:04Z</time>
			</trkpt>
			<trkpt lat="44.101488729938865" lon="-76.92231561988592">
				<ele>83.5999984741211</ele>
				<time>2024-08-24T19:17:07Z</time>
			</trkpt>
			<trkpt lat="44.101545894518495" lon="-76.92238460294902">
				<ele>83.4000015258789</ele>
				<time>2024-08-24T19:17:10Z</time>
			</trkpt>
			<trkpt lat="44.10162141546607" lon="-76.92249759100378">
				<ele>83.4000015258789</ele>
				<time>2024-08-24T19:17:16Z</time>
			</trkpt>
			<trkpt lat="44.10172292031348" lon="-76.92271183244884">
				<ele>83.80000305175781</ele>
				<time>2024-08-24T19:17:23Z</time>
			</trkpt>
			<trkpt lat="44.101816797629" lon="-76.9229123275727">
				<ele>84</ele>
				<time>2024-08-24T19:17:29Z</time>
			</trkpt>
			<trkpt lat="44.1018548514694" lon="-76.9229691568762">
				<ele>84</ele>
				<time>2024-08-24T19:17:31Z</time>
			</trkpt>
			<trkpt lat="44.10196700133383" lon="-76.92300729453564">
				<ele>83.80000305175781</ele>
				<time>2024-08-24T19:17:36Z</time>
			</trkpt>
			<trkpt lat="44.10201637074351" lon="-76.92293914966285">
				<ele>83.80000305175781</ele>
				<time>2024-08-24T19:17:42Z</time>
			</trkpt>
			<trkpt lat="44.10201192833483" lon="-76.92292121239007">
				<ele>83.80000305175781</ele>
				<time>2024-08-24T19:17:43Z</time>
			</trkpt>
			<trkpt lat="44.101986195892096" lon="-76.92287787795067">
				<ele>83.5999984741211</ele>
				<time>2024-08-24T19:17:46Z</time>
			</trkpt>
			<trkpt lat="44.10197270102799" lon="-76.92286362871528">
				<ele>83.5999984741211</ele>
				<time>2024-08-24T19:17:50Z</time>
			</trkpt>
			<trkpt lat="44.10198242403567" lon="-76.92284476943314">
				<ele>83.80000305175781</ele>
				<time>2024-08-24T19:18:00Z</time>
			</trkpt>
			<trkpt lat="44.10196012817323" lon="-76.92277729511261">
				<ele>83.80000305175781</ele>
				<time>2024-08-24T19:18:09Z</time>
			</trkpt>
			<trkpt lat="44.10192936658859" lon="-76.92272557877004">
				<ele>84</ele>
				<time>2024-08-24T19:18:13Z</time>
			</trkpt>
			<trkpt lat="44.10186759196222" lon="-76.92259901203215">
				<ele>84.19999694824219</ele>
				<time>2024-08-24T19:18:20Z</time>
			</trkpt>
			<trkpt lat="44.10180011764169" lon="-76.92247093655169">
				<ele>84.19999694824219</ele>
				<time>2024-08-24T19:18:29Z</time>
			</trkpt>
			<trkpt lat="44.10178469493985" lon="-76.92243631929159">
				<ele>84.19999694824219</ele>
				<time>2024-08-24T19:18:32Z</time>
			</trkpt>
			<trkpt lat="44.10173314623535" lon="-76.92234571091831">
				<ele>84.4000015258789</ele>
				<time>2024-08-24T19:18:39Z</time>
			</trkpt>
			<trkpt lat="44.1017219144851" lon="-76.92232886329293">
				<ele>84.5999984741211</ele>
				<time>2024-08-24T19:18:40Z</time>
			</trkpt>
			<trkpt lat="44.10161931999028" lon="-76.92221051082015">
				<ele>84.80000305175781</ele>
				<time>2024-08-24T19:18:46Z</time>
			</trkpt>
			<trkpt lat="44.10143801942468" lon="-76.92203155718744">
				<ele>84.19999694824219</ele>
				<time>2024-08-24T19:18:55Z</time>
			</trkpt>`
