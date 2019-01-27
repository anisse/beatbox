package main

import (
	"context"
	"io"
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
