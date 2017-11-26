package main

import (
	"fmt"
	"os"
	"os/signal"
	"sort"
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

func PrintOnce(opt Options) {
	lines, err := Tail(opt.logName, opt.n)
	if err != nil {
		if _, ok := err.(*os.PathError); ok {
			fmt.Fprintf(os.Stderr, "No such file: %s\n", opt.logName)
			os.Exit(1)
		}
		panic(err)
	}

	countLabels := make([]int, 0)
	accs := make([]Access, 0)
	for _, line := range lines {
		fields := make([]string, len(opt.labels))

		for _, lvalue := range strings.Split(line, "\t") {
			pos := strings.IndexRune(lvalue, ':')
			if pos == -1 {
				// ignore broken values
				continue
			}
			for i, label := range opt.labels {
				if label == lvalue[:pos] {
					if opt.labelRegexps[i] != nil {
						fields[i] = opt.labelRegexps[i].FindString(lvalue[pos+1:])
					} else {
						fields[i] = lvalue[pos+1:]
					}
					break
				}
			}
		}

		keyExists := false
		for i, k := range accs {
			if EqSlices(k.fields, fields) {
				countLabels[i]++
				keyExists = true
				break
			}
		}

		if !keyExists {
			accs = append(accs, Access{fields: fields, index: len(accs)})
			countLabels = append(countLabels, 1)
		}
	}

	// calculate width
	width := make([]int, len(opt.labels))
	for i, label := range opt.labels {
		if width[i] < len(label) {
			width[i] = len(label)
		}
	}
	for _, acc := range accs {
		for i, f := range acc.fields {
			if width[i] < len(f) {
				width[i] = len(f)
			}
		}
	}

	// calculate totals
	totals := make(map[string]int)
	for i, acc := range accs {
		if _, ok := totals[acc.fields[0]]; !ok {
			totals[acc.fields[0]] = 0
		}
		totals[acc.fields[0]] += countLabels[i]
	}

	// make separator
	sepLength := 6
	for _, w := range width {
		sepLength += w + 2
	}
	if opt.printRate {
		sepLength += 9
	}
	var separatorBytes = make([]byte, sepLength)
	for i := 0; i < sepLength; i++ {
		separatorBytes[i] = '-'
	}
	separator := string(separatorBytes)

	sort.Sort(ByFields(accs))

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
		fmt.Println("   (rate)")
	} else {
		fmt.Println("")
	}

	fmt.Println(separator)

	firstField := accs[0].fields[0]
	for _, acc := range accs {
		// print separator
		if opt.printSeparator && firstField != acc.fields[0] {
			fmt.Println(separator)
			firstField = acc.fields[0]
		}

		// print values
		for i, f := range acc.fields {
			if len(f) == 0 {
				fmt.Printf("%-*s  ", width[i], "*")
			} else {
				fmt.Printf("%-*s  ", width[i], f)
			}
		}
		fmt.Printf("%6d", countLabels[acc.index])

		if opt.printRate {
			rate := float64(countLabels[acc.index]) / float64(totals[acc.fields[0]]) * 100.0
			fmt.Printf("  %6.2f%%", rate)
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
