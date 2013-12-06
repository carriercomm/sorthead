/*
TODO:
	-H 1	keep one header line
	GNU sort options
	-10	top 10 lines
	file names in cmdline, not just stdin
	make it actually faster than GNU sort :)
*/
package main
import (
	"fmt"
	"bufio"
	"os"
	"io"
	"regexp"
	"sort"
)

type NumSort []string
func (slice NumSort) Len() int {
	return len(slice)
}
func (slice NumSort) Swap(i, j int) {
	tmp := slice[i]
	slice[i] = slice[j]
	slice[j] = tmp
}
func toNum(str string) (out int) {
	_, err := fmt.Sscanf(str, "%d", &out)
	if err != nil {
		out = 0
	}
	return
}
func (slice NumSort) Less(i, j int) bool {
	return toNum(slice[i]) > toNum(slice[j])
}

func main() {
	maxlen := 10
	top := make([]string, 0)
	re := regexp.MustCompile("\n$")
	reader := bufio.NewReader(os.Stdin)
	for {
		cur, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		} else if err != nil {
			dief("read error: %s", err)
		}
		cur = re.ReplaceAllString(cur, "")
		top = append(top, cur)
		//sort.Sort(sort.StringSlice(top))
		sort.Sort(NumSort(top))
		if len(top) > maxlen {
			top = top[0:maxlen-1]
		}
		//warnf("top: %v", top) //////
	}
	for _, str := range top {
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
