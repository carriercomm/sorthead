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
var topval []string
var numkey []int

// maximum len(topval) is toplen+1
// (element 0 for current value and elements 1..toplen for the top)
// TODO: panic if toplen is not positive
var toplen int

func init() {
	topval = make([]string, 1)
	numkey = make([]int, 1)
	toplen = 10
}

func toNum(str string) (out int) {
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
		topval = append(topval, "")
		numkey = append(numkey, 0)
	}
	for i := len(topval) - 1; i > pos; i-- {
		topval[i] = topval[i-1]
		numkey[i] = numkey[i-1]
	}
	if pos < len(topval) {
		topval[pos] = topval[0]
		numkey[pos] = curnum
	}
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

	reader := bufio.NewReader(os.Stdin)
	for {
		cur, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("read error: %s", err)
		}
		if cur[len(cur)-1] == '\n' {
			cur = cur[0 : len(cur)-1]
		}
		topval[0] = cur
		add()
	}
	for _, str := range topval[1:] {
		fmt.Println(str)
	}
}
