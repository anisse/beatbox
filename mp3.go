package main

import (
	"context"
	"fmt"
	"os"

	"github.com/anisse/alsa"
	mp3 "github.com/hajimehoshi/go-mp3"
)

func playMp3(ctx context.Context, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open %s: %s", filename, err.Error())
	}
	defer f.Close()

	dec, err := mp3.NewDecoder(f)
	if err != nil {
		return fmt.Errorf("could not decode %s: %s", filename, err.Error())
	}
	defer dec.Close() // closes the file, too

	sampleRate := dec.SampleRate()
	p, err := alsa.NewPlayer(sampleRate, 2, 2, 4096)
	if err != nil {
		return fmt.Errorf("could not init alsa: %s", err.Error())
	}
	defer p.Close()

	_, err = CopyCtx(ctx, p, dec)

	return err
}
