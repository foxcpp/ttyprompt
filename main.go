package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/awnumar/memguard"
	"github.com/foxcpp/ttyprompt/prompt"
	"github.com/foxcpp/ttyprompt/terminal"
	flag "github.com/spf13/pflag"

	assuan "github.com/foxcpp/go-assuan/common"
	pinentry "github.com/foxcpp/go-assuan/pinentry"
	assuansrv "github.com/foxcpp/go-assuan/server"
)

type settings struct {
	debugLog bool
	ttyNum   int
	simple   prompt.DialogSettings
	pinentry bool
}

func parseFlags(flags *settings) {
	flag.BoolVarP(&flags.debugLog, "debug", "D", false, "Log debug information to stderr")

	flag.IntVarP(&flags.ttyNum, "tty", "t", 20, "Number of VT (TTY) to use")

	flag.StringVar(&flags.simple.Title, "title", "", "Title text (simple mode only)")
	flag.StringVarP(&flags.simple.Description, "desc", "d", "", "Detailed description (simple mode only)")
	flag.StringVar(&flags.simple.Prompt, "prompt", "Enter password:", "Prompt text (simple mode only)")

	flag.BoolVar(&flags.pinentry, "pinentry", false, "Enable pinentry emulation mode")

	// Hide "pflag: help requested" if --help used.
	flag.ErrHelp = errors.New("")
	flag.Parse()
}

func enableDebugLog() {
	log.SetPrefix("DEBUG(ttyprompt): ")
	log.SetOutput(os.Stderr)
	log.SetFlags(0)

	assuan.Logger.SetOutput(os.Stderr)
	assuansrv.Logger.SetOutput(os.Stderr)
	pinentry.Logger.SetOutput(os.Stderr)
}

func main() {
	// This way we can return with custom exit code and have defer statements
	// executed.
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	memguard.DisableUnixCoreDumps()
	defer memguard.DestroyAll()

	flags := settings{}
	flags.simple.PassChar = "*"
	parseFlags(&flags)

	if flags.debugLog {
		enableDebugLog()
	} else {
		log.SetOutput(ioutil.Discard)
	}

	lockProcMem()

	log.Println("Getting TTY...")
	tty, err := getTTY(flags.ttyNum)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to get target tty access:", err)
		exitCode = 3
		return
	}
	defer tty.free()
	// In case of signal terminal may be left in unclear state.
	defer tty.file.WriteString(terminal.TermClear + terminal.TermReset)

	prompt.LockTTY(tty.file)
	defer prompt.UnlockTTY(tty.file)

	tty.file.WriteString("ttyprompt acquired this TTY\n")

	// TODO: Polkit agent mode.
	resNotify := make(chan error)

	// We need to handle signals asynchronously.
	// Note: It's dangerous because defer statements in mode functions and
	// others called from there WILL NOT EXECUTED.
	if flags.pinentry {
		log.Println("Going into pinentry mode...")
		go pinentryMode(tty, flags, resNotify)
	} else {
		log.Println("Going into simple mode...")
		go simpleMode(tty, flags, resNotify)
	}

	sigs := make(chan os.Signal)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	select {
	case sig := <-sigs:
		fmt.Fprintln(os.Stderr, "Signal received:", sig)
	case err := <-resNotify:
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			if err.Error() == "ReadPassword: prompt rejected" {
				exitCode = 1
			} else {
				exitCode = 3
			}
		}
	}
}
