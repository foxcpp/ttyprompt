package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/awnumar/memguard"
	"github.com/foxcpp/ttyprompt/terminal"
	flag "github.com/spf13/pflag"
)

func emulatePinentryOptions() (ttyNum int, err error) {
	flag.BoolP("debug", "d", false, "No-op")
	flag.StringP("display", "D", "", "No-op")
	ttyName := flag.StringP("ttyname", "T", "/dev/tty20", "Set the tty terminal node name; only /dev/tty* supported")
	flag.StringP("ttytype", "N", "", "No-op; always 'linux'")
	flag.StringP("lc-ctype", "C", "", "No-op")
	flag.StringP("lc-messages", "M", "", "No-op")
	flag.Int64P("timeout", "o", 0, "No-op; ttyprompt doesn't supports timeouts")
	flag.BoolP("no-global-grab", "g", false, "No-op")
	flag.BoolP("parent-wid", "W", false, "No-op")
	flag.StringP("colors", "c", "", "No-op")
	flag.StringP("ttyalert", "a", "", "No-op")

	// Hide "pflag: help requested" if --help used.
	flag.ErrHelp = errors.New("")
	flag.Parse()

	if !strings.HasPrefix(*ttyName, "/dev/tty") {
		return -1, errors.New("only virtual terminals supported for -T argument")
	}

	ttyNum, err = strconv.Atoi((*ttyName)[8:])
	if err != nil {
		return -1, errors.New("only virtual terminals supported for -T argument")
	}

	return ttyNum, nil
}

func main() {
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	memguard.DisableUnixCoreDumps()
	defer memguard.DestroyAll()

	pinentry := false
	ttyNum := 20
	settings := terminal.DialogSettings{
		Title:       "Experimental! Do not use in production!",
		Description: "Here goes more detailed request dialog",
		Prompt:      "Enter PIN:",
	}
	if !strings.HasSuffix(os.Args[0], "pinentry") {
		flag.IntVarP(&ttyNum, "tty", "t", 20, "Number of VT (TTY) to use")
		flag.StringVar(&settings.Title, "title", "", "Title text (simple mode only)")
		flag.StringVarP(&settings.Description, "desc", "d", "", "Detailed description (simple mode only)")
		flag.StringVar(&settings.Prompt, "prompt", "Enter PIN:", "Prompt text (simple mode only)")

		flag.BoolVar(&pinentry, "pinentry", false, "Enable pinentry emulation mode")

		// Hide "pflag: help requested" if --help used.
		flag.ErrHelp = errors.New("")
		flag.Parse()
	} else {
		pinentry = true
		var err error
		ttyNum, err = emulatePinentryOptions()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			exitCode = 2
			return
		}
	}

	tty, err := getTTY(ttyNum)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to get target tty access:", err)
		exitCode = 2
		return
	}
	defer tty.free()

	// TODO: Polkit agent mode.
	resNotify := make(chan error)

	if pinentry {
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
