package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/awnumar/memguard"
	"github.com/foxcpp/ttyprompt/terminal"
)

func getOwnerUID(path string) int {
	fi, err := os.Stat(path)
	if err != nil {
		return 0
	}

	return int(fi.Sys().(*syscall.Stat_t).Uid)
}

func main() {
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	memguard.DisableUnixCoreDumps()
	defer memguard.DestroyAll()

	// TODO: Parse command-line options.
	settings := terminal.DialogSettings{
		Title:       "Experimental! Do not use in production!",
		Description: "Here goes more detailed request dialog",
		Prompt:      "Enter PIN:",
	}

	tty, err := getTTY(20)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to get target tty access:", err)
		return
	}
	defer tty.free()

	// TODO: Polkit agent mode.
	// TODO: Pinentry mode.
	resNotify := make(chan error)

	go simpleMode(tty, settings, resNotify)

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	select {
	case sig := <-sigs:
		fmt.Fprintln(os.Stderr, "Signal received:", sig)
	case err := <-resNotify:
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			exitCode = 1
		}
	}
}
