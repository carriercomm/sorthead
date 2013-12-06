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
	"os"
	"runtime/pprof"
)

type SortHead struct {
	str    []string
	num    []int
	length int // TODO: panic if it's not positive
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
func (top SortHead) String() (out string) {
	for _, str := range top.str {
		out = out + str + "\n"
	}
	return
}
func (top *SortHead) Add(str string) {
	num := toNum(str)
	pos := top.length
	for i := len(top.str) - 1; i >= 0; i-- {
		numi := top.num[i]
		if num > numi {
			pos = i
		} else if num < numi {
			break
		}
	}
	if len(top.str) < top.length {
		top.str = append(top.str, "")
		top.num = append(top.num, 0)
	}
	for i := len(top.str) - 1; i > pos; i-- {
		top.str[i] = top.str[i-1]
		top.num[i] = top.num[i-1]
	}
	if pos < len(top.str) {
		top.str[pos] = str
		top.num[pos] = num
	}
}

func main() {
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")
	flag.Parse()
	if *cpuprofile != "" {
		pf, err := os.Create(*cpuprofile)
		if err != nil {
			dief("cannot create %s: %s", *cpuprofile, err)
		}
		pprof.StartCPUProfile(pf)
		defer pprof.StopCPUProfile()
	}

	top := SortHead{str: make([]string, 0), length: 10}
	reader := bufio.NewReader(os.Stdin)
	for {
		cur, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			dief("read error: %s", err)
		}
		if cur[len(cur)-1] == '\n' {
			cur = cur[0 : len(cur)-1]
		}
		top.Add(cur)
	}
	fmt.Printf("%s", top)
}

func warnf(f string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Warning: "+f+"\n", args...)
}

func dief(f string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+f+"\n", args...)
	os.Exit(1)
}
