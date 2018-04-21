package main

import (
	"time"

	assuan "github.com/foxcpp/go-assuan/common"
	"github.com/foxcpp/go-assuan/pinentry"
	"github.com/foxcpp/ttyprompt/terminal"
)

func pinentryMode(tty *TTY, settings terminal.DialogSettings, finishNotifyChan chan error) {
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

	getPIN := func(opts pinentry.Settings) (string, *assuan.Error) {
		defer tty.file.WriteString(terminal.TermClear)

		if len(opts.Error) != 0 {
			opts.Desc += "\n\nERROR: " + opts.Error
		}
		if len(opts.Prompt) != 0 {
			opts.Prompt = "Enter PIN:"
		}
		if len(opts.Title) != 0 {
			opts.Title = "pinentry mode"
		}

		buf, n, err := terminal.AskForPassword(tty.file, terminal.DialogSettings{
			Title:       opts.Title,
			Description: opts.Desc,
			Prompt:      opts.Prompt,
		})

		if err != nil {
			return "", &assuan.Error{
				assuan.ErrSrcPinentry, assuan.ErrAssuanServerFault,
				"pinentry", err.Error(),
			}
		}

		return string(buf.Buffer()[:n]), nil
	}
	// Confirm?
	// Message?

	tty.file.WriteString("Running in pinentry mode, waiting for requests...\n")
	err = pinentry.Serve(pinentry.Callbacks{getPIN, nil, nil}, "ttyprompt v0.1.0")
	finishNotifyChan <- err
}
