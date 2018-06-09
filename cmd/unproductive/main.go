package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
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
		fmt.Printf("%d\t%s\t%t\t%s\n",
			time.Now().Unix(), ssid(), runningVPN(), name)
	}
}

func isLocked() bool {
	cmd := exec.Command("pgrep", "-c", "i3lock")
	err := cmd.Run()
	return err == nil
}

func ssid() string {
	out, err := exec.Command("iwgetid", "-r").Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't get ssid: %s\n", err)
		return ""
	}
	return strings.TrimSuffix(string(out), "\n")
}

func runningVPN() bool {
	cmd := exec.Command("pgrep", "-c", "openvpn")
	err := cmd.Run()
	return err == nil
}
