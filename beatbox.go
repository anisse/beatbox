package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/librespot-org/librespot-golang/src/librespot/core"
)

const (
	TAGLEN   = 16
	EMPTY    = "00000000000000000000000000000000"
	SHUTDOWN = "01050000000000000000000000000000"
	STOPTAG  = "020200000"

	NEXTDELAY = 5 * time.Second
)

type player struct {
	tag     string
	list    []string
	index   int
	stopped time.Time
	spotify bool
	session *core.Session
	stop    context.CancelFunc
}

func (p *player) play(w chan<- error) context.CancelFunc {
	ctx, stop := context.WithCancel(context.Background())
	fmt.Println("Now playing:", p.list[p.index])
	go func() {
		if p.spotify {
			playTrack(ctx, p.session, p.list[p.index])
		} else {
			w <- playMp3(ctx, p.list[p.index])
		}
	}()
	return stop
}
func newPlayer() (*player, error) {
	session, err := openSession()
	if err != nil {
		return nil, err
	}
	return &player{
		session: session,
		stop:    func() {},
	}, nil
}

func (p *player) player(c <-chan string) {
	w := make(chan error)
	for {
		select {
		case s := <-c:
			switch {
			case s == SHUTDOWN:
				p.stop()
			case strings.HasPrefix(s, STOPTAG):
				p.stopped = time.Now()
				p.stop()
			case s == p.tag && time.Since(p.stopped) < NEXTDELAY &&
				len(p.list) > p.index+1:
				//next
				p.index++

				p.stop()
				p.stop = p.play(w)
			default:
				list, err := filepath.Glob("/perm/beatbox-data/" + s + "/*.mp3")
				if err != nil || len(list) == 0 {
					// try with spotify
					if playlist, err := ioutil.ReadFile("/perm/beatbox-data/" + s); err == nil && len(playlist) > 0 {
						list, err = playlistTracks(p.session, string(playlist))
						if err != nil {
							fmt.Println("Unknown playlist/track:", err.Error(), string(playlist), s)
							continue
						}
						p.spotify = true
					} else {
						fmt.Println("Unknown tag/command:", s)
						continue
					}
				} else {
					p.spotify = false
				}
				p.index = 0 //first
				p.tag = s
				p.list = list

				p.stop()
				p.stop = p.play(w)
			}
		case err := <-w:
			if err != nil {
				if err.Error() != "context canceled" { //this is us doing the cancellations
					fmt.Println("Error playing mp3", err)
				}
				continue
			}
			//next
			if len(p.list) > p.index+1 {
				p.index++
				p.stop = p.play(w)
			}
		}
	}
}

func main() {
	f, err := os.Open("/dev/hidraw0")
	if err != nil {
		fmt.Println("Error opening HID device", err)
		return
	}
	p, err := newPlayer()
	if err != nil {
		fmt.Println("Error creating player", err)
		return
	}
	c := make(chan string)
	go p.player(c)

	for {
		t := make([]byte, TAGLEN)
		_, err := io.ReadAtLeast(f, t, TAGLEN)
		if err != nil {
			fmt.Println("Error reading", err)
			return
		}
		s := hex.EncodeToString(t)
		if s != EMPTY {
			c <- s
		}
	}
}
