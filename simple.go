package main

import (
	"fmt"

	"github.com/foxcpp/ttyprompt/prompt"
	"github.com/foxcpp/ttyprompt/terminal"
)

/*
In simple mode we just ask for password anf write it to stdout.
Nothing more because this is *simple* mode.

finishNotifyChan is used to report errors because mode functions
run asynchronously.
*/
func simpleMode(tty *TTY, flags settings, finishNotifyChan chan error) {
	buf, n, err := prompt.AskForPassword(tty.file, tty.num, flags.simple)
	if err != nil {
		finishNotifyChan <- err
		return
	}
	// In case of signal this will be not executed, but memguard.DestroyAll
	// from main will so we don't care much about it.
	defer buf.Destroy()

	fmt.Println(string(buf.Buffer()[:n]))

	tty.file.WriteString(terminal.TermClear)
	tty.file.WriteString(terminal.TermReset)

	finishNotifyChan <- nil
}
