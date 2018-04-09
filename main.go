package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/awnumar/memguard"
)

func getOwnerUID(path string) int {
	fi, err := os.Stat(path)
	if err != nil {
		return 0
	}

	return int(fi.Sys().(*syscall.Stat_t).Uid)
}

func main() {
	memguard.DisableUnixCoreDumps()
	defer memguard.DestroyAll()

	// TODO: Parse command-line options.
	prompt := "Application developer forgot to set more useful prompt text."

	tty, err := getTTY(20)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to get target tty access:", err)
		return
	}
	defer tty.free()

	// TODO: Polkit agent mode.
	// TODO: Pinentry mode.
	resNotify := make(chan error)

	go simpleMode(tty, prompt, resNotify)

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	select {
	case sig := <-sigs:
		fmt.Fprintln(os.Stderr, "Signal received:", sig)
	case err := <-resNotify:
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
}
