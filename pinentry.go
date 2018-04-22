package main

import (
	"fmt"
	"time"

	assuan "github.com/foxcpp/go-assuan/common"
	"github.com/foxcpp/go-assuan/pinentry"
	"github.com/foxcpp/ttyprompt/terminal"
)

func getPIN(tty *TTY, opts pinentry.Settings) (string, *assuan.Error) {
	defer tty.file.WriteString(terminal.TermClear)
	if len(opts.Error) != 0 {
		opts.Desc += "\n\nERROR: " + opts.Error
	}
	if len(opts.Prompt) == 0 {
		opts.Prompt = "Enter PIN:"
	}
	if len(opts.Title) == 0 {
		opts.Title = "pinentry mode"
	}

	buf, n, err := terminal.AskForPassword(tty.file, tty.num, terminal.DialogSettings{
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

func confirm(tty *TTY, opts pinentry.Settings) (bool, *assuan.Error) {
	defer tty.file.WriteString(terminal.TermClear)
	if len(opts.Error) == 0 {
		opts.Desc += "\n\nERROR: " + opts.Error
	}
	if len(opts.Title) == 0 {
		opts.Title = "pinentry mode"
	}
	if len(opts.OkBtn) == 0 {
		opts.OkBtn = "OK"
	}
	if len(opts.CancelBtn) == 0 {
		opts.CancelBtn = "Cancel"
	}
	opts.Prompt = fmt.Sprintf("Y = %s, n = %s: ", opts.OkBtn, opts.CancelBtn)

	res, err := terminal.AskToConfirm(tty.file, tty.num, terminal.DialogSettings{
		Title:       opts.Title,
		Description: opts.Desc,
		Prompt:      opts.Prompt,
	})

	if err != nil {
		return false, &assuan.Error{
			assuan.ErrSrcPinentry, assuan.ErrAssuanServerFault,
			"pinentry", err.Error(),
		}
	}

	return res, nil
}

func msg(tty *TTY, opts pinentry.Settings) *assuan.Error {
	defer tty.file.WriteString(terminal.TermClear)
	if len(opts.Error) == 0 {
		opts.Desc += "\n\nERROR: " + opts.Error
	}
	if len(opts.Title) == 0 {
		opts.Title = "pinentry mode"
	}

	err := terminal.ShowMessage(tty.file, tty.num, terminal.DialogSettings{
		Title:       opts.Title,
		Description: opts.Desc,
		Prompt:      opts.Prompt,
	})

	time.Sleep(5 * time.Second)

	if err != nil {
		return &assuan.Error{
			assuan.ErrSrcPinentry, assuan.ErrAssuanServerFault,
			"pinentry", err.Error(),
		}
	}

	return nil
}

func pinentryMode(tty *TTY, settings terminal.DialogSettings, finishNotifyChan chan error) {
	getPINfunc := func(opts pinentry.Settings) (string, *assuan.Error) {
		return getPIN(tty, opts)
	}
	confirmFunc := func(opts pinentry.Settings) (bool, *assuan.Error) {
		return confirm(tty, opts)
	}
	msgFunc := func(opts pinentry.Settings) *assuan.Error {
		return msg(tty, opts)
	}

	tty.file.WriteString("Running in pinentry mode, waiting for requests...\n")
	err := pinentry.Serve(pinentry.Callbacks{getPINfunc, confirmFunc, msgFunc}, "ttyprompt v0.1.0")
	finishNotifyChan <- err
}
