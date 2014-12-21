#!/bin/bash

# no arguments
function usage {
	echo "Usage: $0 <filename> <blocksize> " 
	exit 1
}

# one argument: executable go project
function or_go_build {
	bin=$1
	if [ ! -e $bin ]; then
		go build src/ltcodes/$bin
		if [ $? -ne 0 ] ; then
			echo "Can't build go package [$1], exiting..."
			exit 1
		fi
	fi
}

########## main:
# test whether encode or decode exist
test "$#" -eq 2 || usage
FILE=$1
BSIZE=$2

# build encode and decode
# following assumes we run from the gocode directory
or_go_build encode
or_go_build decode

### actual test
#SEED=1

./encode $FILE $BSIZE | ./decode | cmp $1 - -l
if [ $? -eq 0 ]; then
	echo "test.sh: File correctly transmitted"
fi



