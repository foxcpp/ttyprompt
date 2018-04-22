package main

import (
	"fmt"

	"github.com/foxcpp/ttyprompt/terminal"
)

/*
In simple mode we just ask for password anf write it to stdout.
Nothing more because this is *simple* mode.

finishNotifyChan is used to report errors because mode functions
runs asynchronously.
*/
func simpleMode(tty *TTY, settings terminal.DialogSettings, finishNotifyChan chan error) {
	buf, n, err := terminal.AskForPassword(tty.file, tty.num, settings)
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
