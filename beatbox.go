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

	"github.com/anisse/librespot-golang/src/librespot/core"
)

const (
	TAGLEN   = 16
	EMPTY    = "00000000000000000000000000000000"
	SHUTDOWN = "01050000000000000000000000000000"
	STOPTAG  = "020200000"

	NEXTDELAY = 5 * time.Second

	DATADIR = "/perm/beatbox-data/"
)

type player struct {
	tag     string
	list    []string
	index   int
	stopped time.Time
	spotify bool
	session *core.Session
	stop    context.CancelFunc
	w       chan error
}

func (p *player) play() context.CancelFunc {
	ctx, stop := context.WithCancel(context.Background())
	fmt.Println("Now playing:", p.list[p.index])
	go func() {
		if p.spotify {
			p.w <- playTrack(ctx, p.session, p.list[p.index])
		} else {
			p.w <- playMp3(ctx, p.list[p.index])
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
		w:       make(chan error),
	}, nil
}

func (p *player) process(s string) {
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
		p.stop = p.play()
	default:
		list, err := filepath.Glob(DATADIR + s + "/*.mp3")
		if err != nil || len(list) == 0 {
			// try with spotify
			playlist, err := ioutil.ReadFile(DATADIR + s)
			if err == nil && len(playlist) > 0 {
				list, err = playlistTracks(p.session, string(playlist))
				if err != nil {
					fmt.Println("Unknown playlist/track:",
						err.Error(), string(playlist), s)
					return
				}
				p.spotify = true
			} else {
				fmt.Println("Unknown tag/command:", s)
				return
			}
		} else {
			p.spotify = false
		}
		p.index = 0 //first
		p.tag = s
		p.list = list

		p.stop()
		p.stop = p.play()
	}
}
func (p *player) trackFinished() {
	//next
	if len(p.list) > p.index+1 {
		p.index++
		p.stop = p.play()
	}
}

func (p *player) run(c <-chan string) {
	for {
		select {
		case s := <-c:
			p.process(s)
		case err := <-p.w:
			if err != nil {
				if err.Error() != "context canceled" { //this is us doing the cancellations
					fmt.Println("Error playing ", err)
				}
				continue
			}
			p.trackFinished()
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
	go p.run(c)

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
