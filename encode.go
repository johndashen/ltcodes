// +build encode
package main

import (
	"flag"
	"fmt"
	"./lt"
	"math/rand"
	"strconv"
	"os"
	"time"
)

func Usage() {
		fmt.Println("encode <filename> <blockSize> [seed]")
}

func main() {
	// argument parsing
	flag.Parse()
	if flag.NArg() < 2 || flag.NArg() > 3 {
		Usage()
		return
	}
	filename := flag.Arg(0)
	blockSize, err := strconv.Atoi(flag.Arg(1))
	
	if err != nil {
		Usage()
		return
	}
	rand.Seed(time.Now().UnixNano())
	seed := rand.Uint32()
	if flag.NArg() == 3 {
		tmp, err := strconv.ParseUint(flag.Arg(2), 10, 32)
	
		if err != nil {
			Usage()
			return
		}
		seed = uint32(tmp)
	}

	// get file size
	stats, err := os.Lstat(filename)
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}
	fSize := stats.Size()

	// open file
	f, err := os.Open(filename)
	if err != nil {
		fmt.Errorf(err.Error())
		return
	}

	defer func() {
		if err = f.Close(); err != nil {
			fmt.Errorf(err.Error())
		}
	}()
	
	encoder := lt.NewEncoder(f, uint64(fSize), uint32(blockSize), seed)
	
	for err == nil {
		nextBlock := encoder.NextCodedBlock()
		_, err = os.Stdout.Write(nextBlock.Pack())
	}
}
