package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	TAGLEN   = 16
	EMPTY    = "00000000000000000000000000000000"
	SHUTDOWN = "01050000000000000000000000000000"
	STOPTAG  = "020200000"

	NEXTDELAY = 5 * time.Second
)

type playerState struct {
	tag     string
	list    []string
	index   int
	stopped time.Time
}

func (p *playerState) play(w chan<- error) context.CancelFunc {
	ctx, stop := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "mpg321", p.list[p.index])
	fmt.Println("Now playing:", p.list[p.index])
	err := cmd.Start()
	if err != nil {
		w <- err
	}
	go func() {
		w <- cmd.Wait()
	}()
	return stop
}

func player(c <-chan string) {
	var stop context.CancelFunc = func() {}
	var p playerState
	w := make(chan error)
	for {
		select {
		case s := <-c:
			switch {
			case s == SHUTDOWN:
				stop()
			case strings.HasPrefix(s, STOPTAG):
				p.stopped = time.Now()
				stop()
			default:
				list, err := filepath.Glob(s + "/*.mp3")
				if err != nil || len(list) == 0 {
					fmt.Println("Unknown tag/command:", s)
					continue
				}
				stop()
				if s == p.tag && time.Since(p.stopped) < NEXTDELAY &&
					len(p.list) > p.index+1 {
					p.index++ //next
				} else {
					p.index = 0 //first
				}
				p.tag = s
				p.list = list
				stop = p.play(w)
			}
		case err := <-w:
			if err != nil {
				if err.Error() != "signal: killed" { //this is us doing the killing
					fmt.Println("Error running mpg321", err)
				}
				continue
			}
			//next
			if len(p.list) > p.index+1 {
				p.index++
				stop = p.play(w)
			}
		}
	}
}

func main() {
	f, err := os.Open("/dev/hidraw0")
	if err != nil {
		fmt.Println("Error opening hidraw0", err)
		return
	}
	c := make(chan string)
	go player(c)

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
