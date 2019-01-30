# beatbox

A gokrazy appliance to play music with an NFC reader

[![Build Status](https://travis-ci.org/anisse/beatbox.svg?branch=master)](https://travis-ci.org/anisse/beatbox)
[![Go Report Card](https://goreportcard.com/badge/github.com/anisse/beatbox)](https://goreportcard.com/report/github.com/anisse/beatbox)
[![GoDoc](https://godoc.org/github.com/anisse/beatbox?status.svg)](http://godoc.org/github.com/anisse/beatbox)

# Description

Beatbox is small application that reads data from a Violet Mir:ror using the hidraw interface. Depending on which NFC tag is used, it will play local or remote playlists.

See the [initial concept article](https://anisse.astier.eu/awk-driven-iot.html), and [the FOSDEM talk](https://fosdem.org/2019/schedule/event/embeddedwithgo/) for understanding of the motivations.

[![Demo Video](https://img.youtube.com/vi/I1Vc38DaTQY/0.jpg)](https://www.youtube.com/watch?v=I1Vc38DaTQY)


# Configuration

Since this is a gokrazy appliance, the configuration and data path is hardcoded to /perm/beatbox-data, there you'll need these files:
 - a directory for each tag, named with the tag ID, for example 0201000004042867d200000000000000. This is the tag ID that beatbox shows when an unkown tag is put on the NFC reader.
 - inside this directory, put the mp3 files you want to play

And for Spotify support:
 - "blob.bin", the file generated when using librespot, that serves as an authentication token.
 - "username", a file containing you Spotify username
 - files named like the directories mentionned earlier, that contain a Spotify playlist id, for example user/spotify/playlist/4Y7Ch1e3A3qWP9Gt3g9Hiv

# Building

To build the binary for your machine:
```sh
   go get -u github.com/anisse/beatbox
```

And to build the SDCARD image for your Raspberry Pi 3, if the sdcard device is `/dev/mmcblk0` (be **VERY careful** before running, or you might overwrite the wrong disk/device):

```sh
   go get -u github.com/gokrazy/tools/cmd/gokr-packer
   sudo setfacl -m u:${USER}:rw /dev/mmcblk0
   gokr-packer -kernel_package github.com/anisse/beatbox-kernel -overwrite=/dev/mmcblk0 github.com/anisse/beatbox
```
