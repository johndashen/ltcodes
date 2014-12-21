package main

import "flag"
import "fmt"
import "ltcodes/lt"
import "os"
import "time"
func Usage() {
		fmt.Println("decode")
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

	var decoder lt.BlockDecoder
	wrapUp := func(d []byte) {
		stopT := time.Now()
		Eprintf := func(format string, a ...interface{}) {
			fmt.Fprintf(os.Stderr, format, a...)
		}
		Eprintln := func(a ...interface{}) {
			fmt.Fprintln(os.Stderr, a...)
		}

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
			fmt.Println(err.Error())
			ctrs.drop++
			// return
		} else {
			ctrs.proc ++
		}
		decoder.Include(b)
	}
}
