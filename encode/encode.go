package main

import "flag"
import "fmt"
import "ltcodes/lt"
import "math/rand"
import "strconv"
import "os"
import "time"
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
	encoder := lt.NewEncoder(filename, uint32(blockSize), seed)
	
	for err == nil {
		nextBlock := encoder.NextCodedBlock()
		_, err = os.Stdout.Write(nextBlock.Pack())
	}
}
