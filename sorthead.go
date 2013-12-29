/*
This command is the replacement for `sort ... |head` which works
orders of magnitude faster on large amount of data.

Usage of sorthead:
	-I=false: interactive mode
	-N=10: print the first N lines instead of the first 10
	-k=0: sort by field number N, not the whole string
	-n=false: compare according to string numerical value
	-r=false: reverse the result of comparisons
*/
package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	flag "github.com/ogier/pflag"
	"io"
	"log"
	"os"
	"runtime/pprof"
	"time"
)

type numkeyType float64

// topval[0] and numkey[0] are for current string,
// the rest are current top values
var topval [][]byte
var keyStart, keyEnd []int
var numkey []numkeyType

// maximum len(topval) is toplen+1
// (element 0 for current value and elements 1..toplen for the top)
// TODO: panic if toplen is not positive
var toplen int

func init() {
	topval = [][]byte{{}}
	keyStart = []int{0}
	keyEnd = []int{0}
	numkey = make([]numkeyType, 1)
}

func curToNum() (out numkeyType) {
	for _, char := range topval[0][keyStart[0]:keyEnd[0]] {
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
	a := topval[0][keyStart[0]:keyEnd[0]]
	b := topval[n][keyStart[n]:keyEnd[n]]
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
		keyStart = append(keyStart, 0)
		keyEnd = append(keyEnd, 0)
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
	keyStart[to] = keyStart[from]
	keyEnd[to] = keyEnd[from]
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
	curlen := 0
	inField := false
	curFieldNum := 0
	keyStart[0] = 0
	keyEnd[0] = 0
	for {
		if bufEnd == bufStart {
			if bufEnd > 0 {
				bufStart, bufEnd = 0, 0
			}
			if 0 == len(input) {
				if 0 == flagField {
					keyStart[0] = 0
					keyEnd[0] = curlen
				} else if keyEnd[0] < keyStart[0] {
					keyEnd[0] = curlen
				}
				return len(topval[0]) > 0
			}
			n, err := input[0].Read(buffer[bufEnd:])
			if n > 0 {
				bufEnd += n
				doneBytes += int64(n)
			} else if io.EOF == err {
				input = input[1:]
				continue
			} else {
				log.Fatalln("read error:", err)
			}
		}
		curByte := buffer[bufStart]
		bufStart++
		if flagField != 0 {
			gotWhitespace := ' ' == curByte || '\t' == curByte || '\n' == curByte
			if inField && gotWhitespace { // end of field
				inField = false
				if curFieldNum == flagField {
					keyEnd[0] = curlen
				}
			} else if (!inField) && (!gotWhitespace) { // start of field
				inField = true
				curFieldNum++
				if curFieldNum == flagField {
					keyStart[0] = curlen
				}
			}
		}
		if '\n' == curByte {
			if 0 == flagField {
				keyStart[0] = 0
				keyEnd[0] = curlen
			}
			return true
		}
		// this append() takes more than half of program run time
		topval[0] = append(topval[0], curByte)
		curlen++
	}
	panic("")
}

var flagNum, flagRev, flagInteractive bool
var flagField int
var doneBytes, doneStrings, doneSeconds int64

func main() {
	var cpuprofile *string
	//cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
	flag.BoolVarP(&flagNum, "numeric-sort", "n", false, "compare according to string numerical value")
	flag.BoolVarP(&flagRev, "reverse", "r", false, "reverse the result of comparisons")
	flag.BoolVarP(&flagInteractive, "interactive", "I", false, "interactive mode (it is the default when no -N given)")
	flag.IntVarP(&toplen, "lines", "N", 10, "print the first N lines instead of the first 10")
	flag.IntVarP(&flagField, "key", "k", 0, "sort by field number N, not the whole string")
	flag.Parse()
	flagGiven := map[string]bool{}
	flag.Visit(func(pf *flag.Flag) {
		flagGiven[pf.Name] = true
	})
	if !flagGiven["interactive"] && !flagGiven["lines"] {
		flagInteractive = true
	}
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
	chStop := make(chan bool)
	chDone := make(chan struct{})
	if flagInteractive {
		go draw(chStop, chDone)
	}
	for readString() {
		add()
		doneStrings++
	}
	if flagInteractive {
		chStop <- true
		<-chDone
	}
	finalOutput(0)
}

func finalOutput(code int) {
	for _, str := range topval[1:] {
		fmt.Println(string(str))
	}
	if 0 != code {
		log.Println("Warning: command interrupted, not all data processed")
	}
	os.Exit(code)
}

func draw(chStop chan bool, chDone chan struct{}) {
	if err := termbox.Init(); err != nil {
		log.Fatalln("Cannot initialize termbox", err)
	}
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	chTimer := make(chan struct{})
	go timer(chTimer)
	go poller(chStop)
	drawOnce()
	for {
		select {
		case allDone := <-chStop:
			termbox.Close()
			if allDone {
				chDone <- struct{}{}
				return
			}
			finalOutput(1)
		case <-chTimer:
			drawOnce()
		}
	}
}

func drawOnce() {
	buffer := termbox.CellBuffer()
	for _, cell := range buffer {
		cell.Ch = ' '
		cell.Fg = termbox.ColorDefault
		cell.Bg = termbox.ColorDefault
	}
	xsize, ysize := termbox.Size()
	for i := 0; i < ysize && i < len(topval); i++ {
		var str []rune
		if 0 == i {
			str = []rune(fmt.Sprintf("Processed %d strings in %d seconds. Hit any key to interrupt.", doneStrings, doneSeconds))
		} else {
			str = []rune(string(topval[i]))
		}
		for j := 0; j < xsize && j < len(str); j++ {
			termbox.SetCell(j, i, str[j], termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	if err := termbox.Flush(); err != nil {
		log.Fatalln("Cannot flush termbox:", err)
	}
}

func timer(chTimer chan struct{}) {
	for {
		time.Sleep(time.Second)
		chTimer <- struct{}{}
		doneSeconds++
	}
}

func poller(chStop chan bool) {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			chStop <- false
		}
	}
}

/*
TODO:
	-H 1	keep one header line
	GNU sort options
	-10	top 10 lines
	file names in cmdline, not just stdin
*/
