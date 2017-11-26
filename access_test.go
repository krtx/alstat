package main

import (
	"regexp"
	"testing"
)

func EqAccess(a1, a2 Access) bool {
	return a1.index == a2.index && EqSlices(a1.fields, a2.fields)
}

func TestAccess1(t *testing.T) {
	acc := AccessAggregation{}
	opt := Options{
		keyLabels:  []string{"method"},
		sumLabels:  []string{},
		keyRegexps: []*regexp.Regexp{nil},
	}

	acc.AddLine(opt, "method:POST	status:200")
	acc.AddLine(opt, "method:GET	status:200")

	keysExpected := []Access{
		Access{fields: []string{"POST"}, index: 0},
		Access{fields: []string{"GET"}, index: 1},
	}

	for i, key := range keysExpected {
		if !EqAccess(key, acc.keys[i]) {
			t.Errorf("expected: %v, actual: %v\n", key, acc.keys[i])
		}
	}
}

func TestAccess2(t *testing.T) {
	acc := AccessAggregation{}
	opt := Options{
		keyLabels:  []string{"method", "status"},
		sumLabels:  []string{"time"},
		keyRegexps: []*regexp.Regexp{nil, nil},
	}

	acc.AddLine(opt, "method:POST	status:200	time:1")
	acc.AddLine(opt, "method:GET	status:200	time:2")
	acc.AddLine(opt, "method:POST	status:404	time:3")
	acc.AddLine(opt, "method:GET	status:404	time:4")
	acc.AddLine(opt, "method:POST	status:200	time:5")

	keysExpected := []Access{
		Access{fields: []string{"POST", "200"}, index: 0},
		Access{fields: []string{"GET", "200"}, index: 1},
		Access{fields: []string{"POST", "404"}, index: 2},
		Access{fields: []string{"GET", "404"}, index: 3},
	}

	for i, key := range keysExpected {
		if !EqAccess(key, acc.keys[i]) {
			t.Errorf("expected: %v, actual: %v\n", key, acc.keys[i])
		}
	}

	countsExpected := []int{2, 1, 1, 1}
	for i, count := range countsExpected {
		if count != acc.counts[i] {
			t.Errorf("expected: %v, actual %v\n", countsExpected, acc.counts)
			break
		}
	}

	sumsExpected := []int{6, 2, 3, 4}
	for i, sum := range sumsExpected {
		if len(acc.sums[i]) != 1 || sum != acc.sums[i][0] {
			t.Errorf("expected: %v, actual %v\n", sumsExpected, acc.sums)
			break
		}
	}
}
