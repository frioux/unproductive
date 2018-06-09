package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Entry struct {
	Count   int
	Entries map[string]*Entry
}

func (e Entry) Sum() int {
	ret := e.Count

	for _, e := range e.Entries {
		ret += e.Sum()
	}

	return ret
}

func (e Entry) Render(name []string) string {
	ret := "" // switch to []byte
	for k, i := range e.Entries {
		ret += i.Render(append(name, k))
	}

	ret += fmt.Sprintf("% 7d %s\n", e.Sum(), strings.Join(name, "/"))

	return ret
}

func NewEntry() *Entry {
	return &Entry{Entries: make(map[string]*Entry)}
}

func main() {
	s := bufio.NewScanner(os.Stdin)
	var path []string
	entries := NewEntry()

	for s.Scan() {
		err := json.Unmarshal([]byte(s.Text()), &path)
		if err != nil {
			fmt.Fprintln(os.Stderr, "decoding stdin:", err)
			os.Exit(1)
		}

		cur := entries
		for _, seg := range path {
			if found, ok := cur.Entries[seg]; ok {
				cur = found
			} else {
				cur.Entries[seg] = NewEntry()
				cur = cur.Entries[seg]
			}
		}
		cur.Count++

	}

	fmt.Println(entries.Render([]string{}))
	// fmt.Printf("%#v\n", entries)
	if err := s.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}
}
