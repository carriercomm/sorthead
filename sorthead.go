/*
TODO:
	-H 1	keep one header line
	GNU sort options
	-10	top 10 lines
	file names in cmdline, not just stdin
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"runtime/pprof"
)

// topval[0] and numkey[0] are for current string,
// the rest are current top values
var topval [][]rune
var numkey []int

// maximum len(topval) is toplen+1
// (element 0 for current value and elements 1..toplen for the top)
// TODO: panic if toplen is not positive
var toplen int

func init() {
	topval = [][]rune{{}}
	numkey = make([]int, 1)
	toplen = 10
}

func toNum(str []rune) (out int) {
	for _, char := range str {
		if char >= '0' && char <= '9' {
			out = 10*out + int(char-'0')
		} else {
			break
		}
	}
	return
}
func add() {
	curnum := toNum(topval[0])
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
		topval = append(topval, []rune{})
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

var reader *bufio.Reader

func init() {
	reader = bufio.NewReader(os.Stdin)
}
func readString() bool {
	topval[0] = topval[0][0:0]
	for {
		r, _, err := reader.ReadRune()
		if r == '\n' {
			return true
		} else if err == io.EOF {
			return false
		} else if err != nil {
			log.Fatalf("read error: %s", err)
		} else {
			topval[0] = append(topval[0], r)
		}
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
	for readString() {
		add()
	}
	for _, str := range topval[1:] {
		fmt.Println(string(str))
	}
}
