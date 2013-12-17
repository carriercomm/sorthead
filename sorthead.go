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

type numkeyType float64

// topval[0] and numkey[0] are for current string,
// the rest are current top values
var topval [][]byte
var numkey []numkeyType

// maximum len(topval) is toplen+1
// (element 0 for current value and elements 1..toplen for the top)
// TODO: panic if toplen is not positive
var toplen int

func init() {
	topval = [][]byte{{}}
	numkey = make([]numkeyType, 1)
}

func curToNum() (out numkeyType) {
	for _, char := range topval[0] {
		if char >= '0' && char <= '9' {
			out = 10*out + numkeyType(char-'0')
		} else {
			break
		}
	}
	return
}
func strMore(n int) bool {
	//if n >= len(topval) {
	//	log.Fatalf("topval: %d >= %d", n, len(topval))
	//}
	a := topval[0][:]
	b := topval[n][:]
	for i := 0; i < len(a); i++ {
		if i >= len(b) {
			return true
		}
		if a[i] == b[i] {
			continue
		}
		return a[i] > b[i]
	}
	return true
}
func xor(a, b bool) bool {
	if a {
		return !b
	}
	return b
}
func add() {
	var curnum numkeyType
	if flagNum {
		curnum = curToNum()
	}
	pos := len(topval)
	for i := len(topval) - 1; i > 0; i-- {
		numi := numkey[i]
		if flagNum {
			if xor(flagRev, curnum > numi) {
				break
			}
		} else {
			if xor(flagRev, strMore(i)) {
				break
			}
		}
		pos = i
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

var flagNum, flagRev bool

func main() {
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.BoolVar(&flagNum, "n", false, "compare according to string numerical value")
	flag.BoolVar(&flagRev, "r", false, "reverse the result of comparisons")
	flag.IntVar(&toplen, "N", 10, "print the first N lines instead of the first 10")
	flag.Parse()
	if toplen < 1 {
		log.Fatalf("-N must have positive argument")
	}
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
