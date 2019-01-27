package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/librespot-org/librespot-golang/src/Spotify"
	"github.com/librespot-org/librespot-golang/src/librespot"
	"github.com/librespot-org/librespot-golang/src/librespot/core"
	"github.com/librespot-org/librespot-golang/src/librespot/utils"
)

func getPreferredAudioFile(track *Spotify.Track) (selectedFile *Spotify.AudioFile) {
	for _, file := range track.GetFile() {
		// Take highest quality available ogg vorbis file
		if (file.GetFormat() == Spotify.AudioFile_OGG_VORBIS_96 && selectedFile == nil) ||
			(file.GetFormat() == Spotify.AudioFile_OGG_VORBIS_160 && (selectedFile == nil || selectedFile.GetFormat() == Spotify.AudioFile_OGG_VORBIS_96)) ||
			(file.GetFormat() == Spotify.AudioFile_OGG_VORBIS_320) {
			selectedFile = file
		}
	}

	return
}

func playTrack(ctx context.Context, session *core.Session, id string) error {
	track, err := session.Mercury().GetTrack(utils.Base62ToHex(id))
	if err != nil {
		return err
	}
	selectedFile := getPreferredAudioFile(track)
	if selectedFile == nil {
		for _, altTrack := range track.GetAlternative() {
			selectedFile = getPreferredAudioFile(altTrack)

			if selectedFile != nil {
				break
			}
		}
	}

	if selectedFile == nil {
		return fmt.Errorf("Could not find track for id: %s", id)
	}

	audioFile, err := session.Player().LoadTrack(selectedFile, track.GetGid())
	if err != nil {
		return err
	}
	defer audioFile.Cancel()

	err = playOgg(ctx, audioFile)
	return err
}

func playlistTracks(session *core.Session, id string) ([]string, error) {
	list, err := session.Mercury().GetPlaylist(id)
	if err != nil {
		return nil, err
	}

	contents := list.GetContents()
	if contents == nil {
		return nil, fmt.Errorf("No contents in playlist")
	}
	items := contents.GetItems()
	if items == nil {
		return nil, fmt.Errorf("No items in playlist content")
	}

	if len(items) == 0 {
		return []string{}, nil
	}

	tracks := make([]string, 0, len(items))
	for _, v := range items {
		t := strings.TrimPrefix(*v.Uri, "spotify:track:")
		tracks = append(tracks, t)
	}
	return tracks, nil
}

func openSession() (*core.Session, error) {
	const (
		usernameFile = "/perm/beatbox-data/username"
		deviceName   = "beatbox"
		blobFile     = "/perm/beatbox-data/blob.bin"
	)

	blobBytes, err := ioutil.ReadFile(blobFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read auth blob from %s: %s\n", blobFile, err.Error())
	}

	userBytes, err := ioutil.ReadFile(usernameFile)
	if err != nil {
		return nil, fmt.Errorf("unable to username from %s: %s\n", usernameFile, err.Error())
	}

	session, err := librespot.LoginSaved(string(userBytes), blobBytes, deviceName)
	if err != nil {
		return nil, fmt.Errorf("opening spotify session: %s", err.Error())
	}

	return session, nil
}
