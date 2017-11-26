package main

import (
	"sort"
	"strconv"
	"strings"
)

type Access struct {
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

// add one ltsv line
func (a *AccessAggregation) AddLine(opt Options, line string) {
	key := make([]string, len(opt.keyLabels))
	sums := make([]int, len(opt.sumLabels))

	// parse a ltsv line
	for _, lvalue := range strings.Split(line, "\t") {
		pos := strings.IndexRune(lvalue, ':')
		if pos == -1 {
			// ignore broken values
			continue
		}

		// construct the key
		for i, label := range opt.keyLabels {
			if label == lvalue[:pos] {
				if opt.keyRegexps[i] != nil {
					key[i] = opt.keyRegexps[i].FindString(lvalue[pos+1:])
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

	// add a key
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
