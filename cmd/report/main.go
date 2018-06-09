package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

var showPercents bool

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

func (e Entry) RenderWithPercent(name []string) string {
	return e.renderWithPercent(name, e.Sum())
}

func (e Entry) renderWithPercent(name []string, rootTotal int) string {
	ret := "" // switch to []byte
	for k, i := range e.Entries {
		ret += i.renderWithPercent(append(name, k), rootTotal)
	}

	s := e.Sum()
	ret += fmt.Sprintf("%d\t(%d%%)\t%s\n", s, int(100*s/rootTotal), strings.Join(name, "/"))

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
	flag.BoolVar(&showPercents, "show-percents", false, "Set to include percentage of total")
	flag.Parse()
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

	if showPercents {
		fmt.Print(entries.RenderWithPercent([]string{}))
	} else {
		fmt.Println(entries.Render([]string{}))
	}
	if err := s.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}
}
