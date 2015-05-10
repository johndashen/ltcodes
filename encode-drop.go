// +build encode-drop
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
	fmt.Println("encode <filename> <blockSize> [drop = 0.0]")
}

func main() {
	// argument parsing
	flag.Parse()
	if flag.NArg() < 3 || flag.NArg() > 4 {
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
	drop := float64(0.0)
	if flag.NArg() == 3 {
		drop, err = strconv.ParseFloat(flag.Arg(2), 64)
		if err != nil {
			Usage()
			return
		}
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
	
	badBlock := lt.EmptyCodedBlock(/*fileSize*/ 137, /*	blockSize*/ 37)

	firstBlock := true
	for err == nil {
		nextBlock := encoder.NextCodedBlock()
		if firstBlock || rand.Float64() > drop {
			_, err = os.Stdout.Write(nextBlock.Pack())
			firstBlock = false
		} else {
			os.Stdout.Write(badBlock.Pack())
		}
	}
}
