/*
TODO:
	-H 1	keep one header line
	GNU sort options
	-10	top 10 lines
	file names in cmdline, not just stdin
*/
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
)

// topval[0] and numkey[0] are for current string,
// the rest are current top values
var topval [][]byte
var numkey []int64

// maximum len(topval) is toplen+1
// (element 0 for current value and elements 1..toplen for the top)
// TODO: panic if toplen is not positive
var toplen int

func init() {
	topval = [][]byte{{}}
	numkey = make([]int64, 1)
	toplen = 10
}

func curToNum() (out int64) {
	for _, char := range topval[0] {
		if char >= '0' && char <= '9' {
			out = 10*out + int64(char-'0')
		} else {
			break
		}
	}
	return
}
func add() {
	curnum := curToNum()
	pos := toplen + 1
	for i := len(topval) - 1; i > 0; i-- {
		numi := numkey[i]
		if curnum > numi {
			pos = i
		} else if curnum < numi {
			break
		}
	}
	if len(topval) < toplen+1 {
		topval = append(topval, []byte{})
		numkey = append(numkey, 0)
	}
	for i := len(topval) - 1; i > pos; i-- {
		copyVal(i, i-1)
	}
	if pos < len(topval) {
		numkey[0] = curnum
		copyVal(pos, 0)
	}
}

func copyVal(to, from int) {
	if from >= len(topval) || from < 0 || to >= len(topval) || to < 0 {
		log.Fatalf("copyVal bad index: to=%d to=%d len=%d", to, from, len(topval))
	}
	numkey[to] = numkey[from]
	copy(topval[to], topval[from])
	if len(topval[to]) < len(topval[from]) {
		topval[to] = append(topval[to], topval[from][len(topval[to]):]...)
	} else if len(topval[to]) > len(topval[from]) {
		topval[to] = topval[to][:len(topval[from])]
	}
}

var buffer [1024]byte
var bufStart, bufEnd int // automatically 0
var input []io.Reader

func readString() bool {
	topval[0] = topval[0][0:0]
	for {
		if bufEnd == bufStart {
			if bufEnd > 0 {
				bufStart, bufEnd = 0, 0
			}
			if 0 == len(input) {
				return len(topval[0]) > 0
			}
			n, err := input[0].Read(buffer[bufEnd:])
			if n > 0 {
				bufEnd += n
			} else if io.EOF == err {
				input = input[1:]
				continue
			} else {
				log.Fatalln("read error:", err)
			}
		}
		curByte := buffer[bufStart]
		bufStart++
		if '\n' == curByte {
			return true
		}
		// this append() takes more than half of program run time
		topval[0] = append(topval[0], curByte)
	}
	panic("")
}

func main() {
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	if *cpuprofile != "" {
		pf, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatalf("cannot create %s: %s", *cpuprofile, err)
		}
		pprof.StartCPUProfile(pf)
		defer pprof.StopCPUProfile()
	}
	input = []io.Reader{os.Stdin}
	for readString() {
		add()
	}
	for _, str := range topval[1:] {
		fmt.Println(string(str))
	}
}
