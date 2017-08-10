#!/bin/sh

set -x


sudo stdbuf -o0 hexdump -C /dev/hidraw0 | awk -W interactive '
/02 01 00 00 08 d0 02 1a  03 52 c1 1a 01 00 00 00/ { print "file1 "; play=1; file="file1.mp3" ; }
/02 01 00 00 08 d0 02 1a  03 52 c1 4b ad 00 00 00/ { print "file2 "; play=1; file="file2.mp3" ; }
/02 01 00 00 04 3f d7 5f  35 00 00 00 00 00 00 00/ { print "dir1 "; play=1; file="dir1/*.mp3" ; }
/02 01 00 00 04 ea 4d 9d  11 00 00 00 00 00 00 00/ { print "dir2 "; play=1; file="dir2/*.mp3" ; }
/02 01 00 00 04 8f 74 05  5e 00 00 00 00 00 00 00/ { print "dir3 "; play=1; file="dir3/*.mp3" ; }
/02 01 00 00 04 6f af ea  c2 00 00 00 00 00 00 00/ { print "dir4 "; play=1; file="dir4/*.mp3" ; }
/02 01 00 00 04 3f 08 36  be 00 00 00 00 00 00 00/ { print "dir5 "; play=1; file="dir5/*.mp3" ; }
/02 01 00 00 04 8f 85 e5  e3 00 00 00 00 00 00 00/ { print "dir6 "; play=1; file="dir6/*.mp3" ; }
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
