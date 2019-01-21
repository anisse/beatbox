package main

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/anisse/alsa"
	mp3 "github.com/hajimehoshi/go-mp3"
)

// Code below from http://ixday.github.io/post/golang-cancel-copy/
type readerFunc func(p []byte) (n int, err error)

func (rf readerFunc) Read(p []byte) (n int, err error) { return rf(p) }

// Note that we cancel on the reader; ensuring that any consumed byte is
// attempted to be written, but it also means that:
//  - we reduce the granularity of the the cancellation window
//  - anyway, if any of read or write blocks, it blocks cancellation, but this
//    is out of scope
func CopyCtx(ctx context.Context, dst io.Writer, src io.Reader) (int64, error) {
	n, err := io.Copy(dst, readerFunc(func(p []byte) (int, error) {

		select {
		case <-ctx.Done():
			return 0, ctx.Err()
		default:
			return src.Read(p)
		}
	}))
	return n, err
}

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
