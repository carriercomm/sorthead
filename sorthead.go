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

var topval []string
var numkey []int
var toplen int // TODO: panic if it's not positive

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
func add(str string) {
	num := toNum(str)
	pos := toplen
	for i := len(topval) - 1; i >= 0; i-- {
		numi := numkey[i]
		if num > numi {
			pos = i
		} else if num < numi {
			break
		}
	}
	if len(topval) < toplen {
		topval = append(topval, "")
		numkey = append(numkey, 0)
	}
	for i := len(topval) - 1; i > pos; i-- {
		topval[i] = topval[i-1]
		numkey[i] = numkey[i-1]
	}
	if pos < len(topval) {
		topval[pos] = str
		numkey[pos] = num
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

	topval = make([]string, 0)
	toplen = 10
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
		add(cur)
	}
	for _, str := range topval {
		fmt.Println(str)
	}
}
