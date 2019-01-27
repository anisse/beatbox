module beatbox

require (
	github.com/anisse/alsa v0.0.0-20190117080618-9d5ae27c4011
	github.com/hajimehoshi/go-mp3 v0.1.1
	github.com/jfreymuth/oggvorbis v1.0.0
	github.com/jfreymuth/vorbis v1.0.0 // indirect
	github.com/librespot-org/librespot-golang/src/Spotify v0.0.1
	github.com/librespot-org/librespot-golang/src/librespot v0.0.1
)

replace github.com/librespot-org/librespot-golang/src/Spotify => github.com/anisse/librespot-golang/src/Spotify v0.0.0-20190127000001-d164560c700e79a550a

replace github.com/librespot-org/librespot-golang/src/librespot => github.com/anisse/librespot-golang/src/librespot v0.0.0-20190127000001-d164560c700e79a550a
