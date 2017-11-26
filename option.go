package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

type Options struct {
	logName        string
	n              int
	keyLabels      []string
	sumLabels      []string
	keyRegexps     []*regexp.Regexp // used to extract keys
	printSeparator bool
	printRate      bool
	interval       int
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return ""
}

func (i *arrayFlags) Set(v string) error {
	*i = append(*i, v)
	return nil
}

func (opt *Options) Load() {
	ls := arrayFlags{}
	sumLabels := arrayFlags{}

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] LOGFILE\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.IntVar(&opt.n, "n", 1000, "number of tail lines to read")
	flag.BoolVar(&opt.printSeparator, "sep", false, "print separator")
	flag.BoolVar(&opt.printRate, "rate", false, "print rate")
	flag.Var(&ls, "l", "labels (-l can be used multiple times)")
	flag.Var(&sumLabels, "sum", "labels to sum up their values (-s can be used multiple times)")
	flag.IntVar(&opt.interval, "c", 1, "interval seconds")

	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	opt.logName = flag.Args()[0]

	opt.keyLabels = make([]string, len(ls))
	opt.keyRegexps = make([]*regexp.Regexp, len(ls))
	for i, l := range ls {
		pos := strings.IndexRune(l, ':')
		if pos == -1 {
			opt.keyLabels[i] = l
		} else {
			opt.keyLabels[i] = l[:pos]
			re, err := regexp.Compile(l[pos+1:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", err)
				os.Exit(1)
			}
			opt.keyRegexps[i] = re
		}
	}

	// reject repeated keyLabels
	s := make([]string, len(opt.keyLabels))
	copy(s, opt.keyLabels)
	sort.Strings(s)
	for i := 0; i < len(s)-1; i++ {
		if s[i] == s[i+1] {
			fmt.Fprintf(os.Stderr, "Invalid option: repeated labels found: %s\n", s[i])
			fmt.Fprintf(os.Stderr, "Usage: \n")
			flag.PrintDefaults()
			os.Exit(1)
		}
	}

	opt.sumLabels = make([]string, len(sumLabels))
	for i, _ := range sumLabels {
		opt.sumLabels[i] = sumLabels[i]
	}

	if len(opt.keyLabels) < 1 || opt.n <= 0 {
		flag.Usage()
		os.Exit(1)
	}
}
