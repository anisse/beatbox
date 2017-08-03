#!/bin/sh

set -x


sudo stdbuf -o0 hexdump -C /dev/hidraw0 | awk -W interactive '
/02 01 00 00 08 d0 02 1a  03 52 c1 1a 01 00 00 00/ { print "file1 "; play=1; file="file1.mp3" ; }
/02 01 00 00 08 d0 02 1a  03 52 c1 4b ad 00 00 00/ { print "file2 "; play=1; file="file2.mp3" ; }
/02 01 00 00 04 3f d7 5f  35 00 00 00 00 00 00 00/ { print "dir1"; play=1; file="dir1/*.mp3" ; }
/  02 02 00 00 0. |01 05 00 00 00 00 00 00  00 00 00 00 00 00 00 00/ { print "stop"; system("killall -q mpg321"); }
{ 
if (play) {
	system("mpg321 -q " file " & ");
	}
play=0 ;
file=x ;
}

END {
}
'
