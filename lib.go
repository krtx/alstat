package main

import (
	"os"
	"strings"
)

// Check if slices are equal
func EqSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if b[i] != v {
			return false
		}
	}
	return true
}

// Reverse string slice
func Reverse(s []string) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Read n lines of tails of file
func Tail(name string, n int) (res []string, err error) {
	var bufSize int64 = 8 * 1024

	file, err := os.Open(name)
	if err != nil {
		return res, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return res, err
	}

	currentIndex := stat.Size()
	reachHead := false

	// 以下の下線部を読み込んだとき、"cdef" はすぐには res には加えない
	// ...\nabcdef\nghijkl
	//        ^^^^^^^^^^^^
	// この中途半端な文字列を保存するための変数
	head := ""

	for {
		currentIndex -= bufSize

		if currentIndex < 0 {
			reachHead = true
			file.Seek(0, 0)
		} else {
			file.Seek(currentIndex, 0)
		}

		buf := make([]byte, bufSize)
		file.Read(buf)

		if currentIndex < 0 {
			// 読みすぎた分を削除
			buf = buf[:bufSize+currentIndex]
		}

		lines := strings.Split(string(buf)+head, "\n")

		// 残りの読み込むべき行数
		l := n - len(res)
		if reachHead {
			if l > 0 {
				if len(lines) < l {
					res = append(lines, res...)
				} else {
					res = append(lines[len(lines)-l:], res...)
				}
			}
			break
		}

		head = lines[0]
		res = append(lines[1:], res...)

		if len(res) == n {
			break
		} else if len(res) > n {
			excess := len(res) - n
			res = res[excess:]
			break
		}
	}

	return res, nil
}
