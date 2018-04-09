package main

import (
	"fmt"
	"os"
	"time"

	"github.com/foxcpp/ttyprompt/terminal"
)

func switchToOriginalVT(outputTty *os.File, num int) error {
	if err := terminal.SwitchVTThrough(outputTty.Fd(), num); err != nil {
		fmt.Fprintln(os.Stderr, "failed to switch TTY back:", err)
		outputTty.WriteString("\nOops! We can't switch TTYs. Do it manually (i.e. Ctrl+Alt+F1).\n")
		return err
	}
	return nil
}

/*
In simple mode we just ask for password anf write it to stdout.
Nothing more because this is *simple* mode.

finishNotifyChan is used to report errors because mode functions
runs asynchronously.
*/
func simpleMode(tty *TTY, settings terminal.DialogSettings, finishNotifyChan chan error) {
	firsttty, err := terminal.CurrentVT()
	if err != nil {
		finishNotifyChan <- err
		return
	}

	if err := terminal.SwitchVTThrough(tty.file.Fd(), tty.num); err != nil {
		finishNotifyChan <- err
		return
	}

	defer func() {
		if err := switchToOriginalVT(tty.file, firsttty); err != nil {
			time.Sleep(time.Second * 5)
		}
	}()

	buf, n, err := terminal.AskForPassword(tty.file, settings)
	if err != nil {
		finishNotifyChan <- err
		return
	}
	defer buf.Destroy()

	fmt.Println(string(buf.Buffer()[:n]))

	tty.file.WriteString(terminal.TermClear)
	tty.file.WriteString(terminal.TermReset)

	finishNotifyChan <- nil
}
