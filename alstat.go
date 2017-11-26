package main

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Access struct {
	// field[0] は primary field として特別扱いされることがある
	fields []string
	index  int
}

type ByFields []Access

func (b ByFields) Len() int      { return len(b) }
func (b ByFields) Swap(i, j int) { b[i], b[j] = b[j], b[i] }
func (b ByFields) Less(i, j int) bool {
	// dictionary order
	for k, _ := range b[i].fields {
		if b[i].fields[k] != b[j].fields[k] {
			return b[i].fields[k] < b[j].fields[k]
		}
	}
	return false
}

type AccessAggregation struct {
	keys   []Access
	counts []int
	sums   [][]int
}

func (a *AccessAggregation) Add(key []string, sums []int) {
	for i, k := range a.keys {
		// key exists
		if EqSlices(k.fields, key) {
			a.counts[i]++
			for j, sum := range sums {
				a.sums[i][j] += sum
			}
			return
		}
	}

	// key does not exist
	a.keys = append(a.keys, Access{fields: key, index: len(a.keys)})
	a.counts = append(a.counts, 1)
	a.sums = append(a.sums, sums)
}

func (a *AccessAggregation) Sort() {
	sort.Sort(ByFields(a.keys))
}

func PrintOnce(opt Options) {
	lines, err := Tail(opt.logName, opt.n)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			fmt.Fprintf(os.Stderr, "No such file: %s\n", opt.logName)
			os.Exit(1)
		}
		panic(err)
	}

	acc := AccessAggregation{}

	for _, line := range lines {
		key := make([]string, len(opt.labels))
		sums := make([]int, len(opt.sumLabels))
		for _, lvalue := range strings.Split(line, "\t") {
			pos := strings.IndexRune(lvalue, ':')
			if pos == -1 {
				// ignore broken values
				continue
			}

			// construct the key
			for i, label := range opt.labels {
				if label == lvalue[:pos] {
					if opt.labelRegexps[i] != nil {
						key[i] = opt.labelRegexps[i].FindString(lvalue[pos+1:])
					} else {
						key[i] = lvalue[pos+1:]
					}
					break
				}
			}

			// collect sums
			for i, label := range opt.sumLabels {
				if label == lvalue[:pos] {
					if s, err := strconv.Atoi(lvalue[pos+1:]); err == nil {
						sums[i] = s
					}
				}
			}
		}

		acc.Add(key, sums)
	}

	// calculate width
	width := make([]int, len(opt.labels))
	for i, label := range opt.labels {
		width[i] = len(label)
	}
	for _, key := range acc.keys {
		for i, f := range key.fields {
			if width[i] < len(f) {
				width[i] = len(f)
			}
		}
	}
	widthSums := make([]int, len(opt.sumLabels))
	for i, label := range opt.sumLabels {
		widthSums[i] = len(label) + 5
	}

	// calculate totals within each primary category
	totals := make(map[string]int)
	for i, key := range acc.keys {
		if _, ok := totals[key.fields[0]]; !ok {
			totals[key.fields[0]] = 0
		}
		totals[key.fields[0]] += acc.counts[i]
	}

	// make separator
	sepLength := 6
	for _, w := range width {
		sepLength += w + 2
	}
	if opt.printRate {
		sepLength += 9
	}
	for _, w := range widthSums {
		sepLength += w + 2
	}
	var separatorBytes = make([]byte, sepLength)
	for i := 0; i < sepLength; i++ {
		separatorBytes[i] = '-'
	}
	separator := string(separatorBytes)

	acc.Sort()

	// clear screen
	if opt.interval >= 1 {
		print("\033[H\033[2J")
	}

	// print labels
	for i, label := range opt.labels {
		fmt.Printf("%-*s  ", width[i], label)
	}
	fmt.Printf("access")
	if opt.printRate {
		fmt.Printf("   (rate)")
	}
	for i, label := range opt.sumLabels {
		fmt.Printf("  %-*s", width[i], "sum("+label+")")
	}
	fmt.Println("")

	fmt.Println(separator)

	firstField := acc.keys[0].fields[0]
	for _, key := range acc.keys {
		// print separator
		if opt.printSeparator && firstField != key.fields[0] {
			fmt.Println(separator)
			firstField = key.fields[0]
		}

		// print key
		for i, f := range key.fields {
			if len(f) == 0 {
				fmt.Printf("%-*s  ", width[i], "*")
			} else {
				fmt.Printf("%-*s  ", width[i], f)
			}
		}

		fmt.Printf("%6d", acc.counts[key.index])
		if opt.printRate {
			rate := float64(acc.counts[key.index]) / float64(totals[key.fields[0]]) * 100.0
			fmt.Printf("  %6.2f%%", rate)
		}

		for i, _ := range opt.sumLabels {
			fmt.Printf("  %*d", widthSums[i], acc.sums[key.index][i])
		}

		fmt.Println("")
	}
}

var opt Options

func main() {
	opt.Load()

	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt)

	if opt.interval <= 0 {
		PrintOnce(opt)
		os.Exit(0)
	}

	t := time.NewTicker(time.Second)
L:
	for {
		select {
		case <-sigc:
			break L
		case <-t.C:
			PrintOnce(opt)
		}
	}

	t.Stop()
}
