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

	b, err := lt.ReadBlockFrom(os.Stdin)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	wrapUp := func(d []byte) {
		stopT := time.Now()
		os.Stderr.Write([]byte(fmt.Sprintf("file received in %fms\n", stopT.Sub(startT).Seconds() * 1000)))
		os.Stderr.Write([]byte(fmt.Sprintf("total size: %d bytes\n", len(d))))
		os.Stdout.Write(d)
	}

	decoder := lt.NewDecoder(b)
	if done, data := decoder.AttemptDone(); done {
		wrapUp(data)
		return
	}
	for err == nil {
		b, err := lt.ReadBlockFrom(os.Stdin)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		decoder.Include(b)

		if done, data := decoder.AttemptDone(); done {
			wrapUp(data)
			return
		}
	}
}
