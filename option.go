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
	labels         []string
	labelRegexps   []*regexp.Regexp
	sumLabels      []string
	printSeparator bool
	printRate      bool
	interval       int
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	if i == nil {
		return ""
	} else {
		return strings.Join(*i, ",")
	}
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

	flag.IntVar(&opt.n, "n", 1000, "number of lines to tail")
	flag.BoolVar(&opt.printSeparator, "sep", false, "print separator")
	flag.BoolVar(&opt.printRate, "rate", false, "print rate")
	flag.Var(&ls, "l", "labels")
	flag.Var(&sumLabels, "sum", "labels to sum up their values")
	flag.IntVar(&opt.interval, "c", 1, "interval")

	flag.Parse()

	if len(flag.Args()) != 1 {
		flag.Usage()
		os.Exit(1)
	}

	opt.logName = flag.Args()[0]

	opt.labels = make([]string, len(ls))
	opt.labelRegexps = make([]*regexp.Regexp, len(ls))
	for i, l := range ls {
		pos := strings.IndexRune(l, ':')
		if pos == -1 {
			opt.labels[i] = l
		} else {
			opt.labels[i] = l[:pos]
			re, err := regexp.Compile(l[pos+1:])
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", err)
				os.Exit(1)
			}
			opt.labelRegexps[i] = re
		}
	}

	// reject repeated labels
	s := make([]string, len(opt.labels))
	copy(s, opt.labels)
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

	fmt.Printf("%v\n", opt)

	if len(opt.labels) < 1 || opt.n <= 0 {
		flag.Usage()
		os.Exit(1)
	}
}
