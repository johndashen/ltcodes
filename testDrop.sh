#!/bin/bash

# no arguments
function usage {
	echo "Usage: $0 <filename> <blocksize> <drop>" 
	exit 1
}

# one argument: executable go project
function or_go_build {
	bin=$1
	if [ ! -e $bin ]; then
		go build ${bin}.go
		if [ $? -ne 0 ] ; then
			echo "Can't build go package [$1], exiting..."
			exit 1
		fi
	fi
}

########## main:
# test whether encode or decode exist
test "$#" -eq 3 || usage
FILE=$1
BSIZE=$2
DROP=$3
# build encode and decode
# following assumes we run from the gocode directory
or_go_build encode-drop
or_go_build decode

### actual test
#SEED=1

./encode-drop $FILE $BSIZE $DROP | ./decode | cmp $1 - -l
if [ $? -eq 0 ]; then
	echo "test.sh: File correctly transmitted"
fi



