package main

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const alphabets = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabets[rand.Intn(len(alphabets))]
	}
	return string(b)
}

func TestTail1(t *testing.T) {
	testFile := filepath.Join(os.TempDir(), "alstat-test-teail"+RandomString(8))

	f, err := os.Create(testFile)
	if err != nil {
		panic(err)
	}
	defer os.Remove(testFile)

	lines := []string{"a", "b", "c"}
	f.WriteString(strings.Join(lines, "\n"))

	tails, _ := Tail(testFile, 1)
	expected := lines[2:]
	if !EqSlices(tails, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, tails)
	}

	tails, _ = Tail(testFile, 2)
	expected = lines[1:]
	if !EqSlices(tails, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, tails)
	}

	tails, _ = Tail(testFile, 3)
	expected = lines
	if !EqSlices(tails, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, tails)
	}

	tails, _ = Tail(testFile, 4)
	expected = lines
	if !EqSlices(tails, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, tails)
	}
}

func TestTail2(t *testing.T) {
	testFile := os.TempDir() + "alstat-test-teail" + RandomString(8)

	f, err := os.Create(testFile)
	if err != nil {
		panic(err)
	}
	defer os.Remove(testFile)

	lines := make([]string, 3)
	for i, _ := range lines {
		lines[i] = RandomString(8*2000 + 10)
	}
	f.WriteString(strings.Join(lines, "\n"))

	tails, _ := Tail(testFile, 1)
	expected := lines[2:]
	if !EqSlices(tails, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, tails)
	}

	tails, _ = Tail(testFile, 2)
	expected = lines[1:]
	if !EqSlices(tails, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, tails)
	}

	tails, _ = Tail(testFile, 3)
	expected = lines
	if !EqSlices(tails, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, tails)
	}

	tails, _ = Tail(testFile, 4)
	expected = lines
	if !EqSlices(tails, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, tails)
	}
}

func TestTail3(t *testing.T) {
	testFile := os.TempDir() + "alstat-test-teail" + RandomString(8)

	f, err := os.Create(testFile)
	if err != nil {
		panic(err)
	}
	defer os.Remove(testFile)

	lines := make([]string, 1000)
	for i, _ := range lines {
		lines[i] = RandomString(rand.Intn(1000) + 100)
	}
	f.WriteString(strings.Join(lines, "\n"))

	tails, _ := Tail(testFile, 100)
	expected := lines[900:]
	if !EqSlices(tails, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, tails)
	}

	tails, _ = Tail(testFile, 800)
	expected = lines[200:]
	if !EqSlices(tails, expected) {
		t.Errorf("expected: %v, actual: %v\n", expected, tails)
	}
}
