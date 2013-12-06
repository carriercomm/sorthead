/*
TODO:
	-H 1	keep one header line
	GNU sort options
	-10	top 10 lines
	file names in cmdline, not just stdin
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
	_, err := fmt.Sscanf(str, "%d", &out)
	if err != nil {
		out = 0
	}
	return
}
func (top NumSort) Less(i, j int) bool {
	return top.num[i] > top.num[j]
}

func main() {
	maxlen := 10
	top := NumSort{ str: make([]string, 0) }
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
		top.str = append(top.str, cur)
		curnum := toNum(cur)
		top.num = append(top.num, curnum)
		//sort.Sort(sort.StringSlice(top))
		sort.Sort(top)
		if len(top.str) > maxlen {
			top.str = top.str[0:maxlen-1]
			top.num = top.num[0:maxlen-1]
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
