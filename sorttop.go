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
	"sort"
)

type NumSort struct {
	str []string
	num []int
}

func (top NumSort) Len() int {
	return len(top.str)
}
func (top NumSort) Swap(i, j int) {
	tmp := top.str[i]
	top.str[i] = top.str[j]
	top.str[j] = tmp
	tmpnum := top.num[i]
	top.num[i] = top.num[j]
	top.num[j] = tmpnum
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
func (top NumSort) Less(i, j int) bool {
	return top.num[i] > top.num[j]
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

	maxlen := 10
	top := NumSort{str: make([]string, 0)}
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
		top.str = append(top.str, cur)
		curnum := toNum(cur)
		top.num = append(top.num, curnum)
		//sort.Sort(sort.StringSlice(top))
		sort.Sort(top)
		if len(top.str) > maxlen {
			top.str = top.str[0 : maxlen-1]
			top.num = top.num[0 : maxlen-1]
		}
		//warnf("top.num: %v", top.num) //////
	}
	for _, str := range top.str {
		fmt.Println(str)
	}
}

func warnf(f string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Warning: "+f+"\n", args...)
}

func dief(f string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, "Error: "+f+"\n", args...)
	os.Exit(1)
}
