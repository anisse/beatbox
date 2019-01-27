package main

import (
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"math"

	"github.com/anisse/alsa"
	"github.com/jfreymuth/oggvorbis"
)

// Reader that does float32le -> s16le conversion
type resampleReader struct {
	dec *oggvorbis.Reader
}

func (r *resampleReader) Read(p []byte) (n int, err error) {
	fBuf := make([]float32, len(p)/2)
	n, err = r.dec.Read(fBuf)
	for i := 0; i < n; i += 1 {
		val := int16(fBuf[i] * math.MaxInt16)
		binary.LittleEndian.PutUint16(p[i*2:], uint16(val))
	}
	return n * 2, err
}

func playOgg(ctx context.Context, r io.Reader) error {
	dec, err := oggvorbis.NewReader(r)
	if err != nil {
		return fmt.Errorf("could not initialize ogg reader: %s", err.Error())
	}
	p, err := alsa.NewPlayer(dec.SampleRate(), dec.Channels(), 2, 4096)
	if err != nil {
		return fmt.Errorf("could not init alsa: %s", err.Error())
	}
	defer p.Close()

	_, err = CopyCtx(ctx, p, &resampleReader{dec})
	return err
}
