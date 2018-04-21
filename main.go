package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/awnumar/memguard"
	"github.com/foxcpp/ttyprompt/terminal"
	flag "github.com/spf13/pflag"
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

	settings := terminal.DialogSettings{
		Title:       "Experimental! Do not use in production!",
		Description: "Here goes more detailed request dialog",
		Prompt:      "Enter PIN:",
	}
	ttyNum := flag.IntP("tty", "t", 20, "Number of VT (TTY) to use")
	flag.StringVar(&settings.Title, "title", "", "Title text (simple mode only)")
	flag.StringVarP(&settings.Description, "desc", "d", "", "Detailed description (simple mode only)")
	flag.StringVar(&settings.Prompt, "prompt", "Enter PIN:", "Prompt text (simple mode only)")

	pinentry := flag.Bool("pinentry", false, "Enable pinentry emulation mode")

	// Hide "pflag: help requested" if --help used.
	flag.ErrHelp = errors.New("")
	flag.Parse()

	tty, err := getTTY(*ttyNum)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to get target tty access:", err)
		return
	}
	defer tty.free()

	// TODO: Polkit agent mode.
	resNotify := make(chan error)

	if *pinentry {
		go pinentryMode(tty, settings, resNotify)
	} else {
		go simpleMode(tty, settings, resNotify)
	}

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
