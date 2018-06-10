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

func (e Entry) Render(name []string, callbacks []func(Entry) string) string {
	ret := "" // switch to []byte
	for k, i := range e.Entries {
		ret += i.Render(append(name, k), callbacks)
	}

	inners := make([]string, len(callbacks))
	for i, cb := range callbacks {
		inners[i] = cb(e)
	}
	ret += fmt.Sprintf("%d\t%s\t%s\n", e.Sum(), strings.Join(inners, "\t"), strings.Join(name, "/"))

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

	cbs := []func(Entry) string{}
	if showPercents {
		total := entries.Sum()
		cbs = append(cbs, func(e Entry) string {
			s := e.Sum()
			return fmt.Sprintf("(%d%%)", int(100*s/total))
		})
	}

	fmt.Println(entries.Render([]string{}, cbs))

	if err := s.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
		os.Exit(1)
	}
}
