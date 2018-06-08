package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/vektra/tai64n"
)

func main() {
	X, err := xgbutil.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal(err)
	}

	for {
		time.Sleep(time.Second)
		if isLocked() {
			continue
		}

		w, err := ewmh.ActiveWindowGet(X)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Coudln't ActiveWindowGet: %s\n", err)
			continue
		}
		name, err := ewmh.WmNameGet(X, w)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Coudln't WmNameGet: %s\n", err)
			continue
		}
		fmt.Printf("%s\t%s\n", tai64n.Now().Label(), name)
	}
}

func isLocked() bool {
	cmd := exec.Command("pgrep", "-c", "i3lock")
	err := cmd.Run()
	return err == nil
}
