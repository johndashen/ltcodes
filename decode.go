// +build decode
package main

import (
	"flag"
     "fmt"
	"./lt"
	"os"
	"time"
)

func Usage() {
		fmt.Println("decode")
}

func Eprintf(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}

func Eprintln(a ...interface{}) {
	fmt.Fprintln(os.Stderr, a...)
}

func main() {
	startT := time.Now()
	// argument parsing - no args
	flag.Parse()
	if flag.NArg() > 0 {
		Usage()
		return
	}

	ctrs := struct{
		in, proc, drop int;
	}{}

	decoder := new(lt.BlockDecoder)

	wrapUp := func(d []byte) {
		stopT := time.Now()
		Eprintln("file received in", stopT.Sub(startT))
		Eprintf("total size: %d bytes\n", len(d))
		Eprintf("packets received: %d, packets processed: %d, packets dropped: %d\n", ctrs.in, ctrs.proc, ctrs.drop)
		Eprintf("bytes/packet: %d, bytes total: %d, rate: %0.3f\n", 
			decoder.BlockSize(), uint32(ctrs.proc) * decoder.BlockSize(), 
			float64(decoder.FileSize()) / float64(uint32(ctrs.proc) * decoder.BlockSize()))
		os.Stdout.Write(d)
	}

	var err error
	for err == nil {
		if done, data := decoder.AttemptDone(); done {
			wrapUp(data)
			return
		}

		b, err := lt.ReadBlockFrom(os.Stdin)

		ctrs.in++
		if err != nil {
			Eprintln(err.Error())
			ctrs.drop++
			// return
		} else if !decoder.Validate(b) {
//			Eprintln("Dropped block found")
			ctrs.drop++
		} else {
			ctrs.proc ++
			decoder.Include(b)
		}
	}
}
