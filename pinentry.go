package main

import (
	"fmt"
	"log"

	"github.com/awnumar/memguard"
	assuan "github.com/foxcpp/go-assuan/common"
	"github.com/foxcpp/go-assuan/pinentry"
	"github.com/foxcpp/ttyprompt/prompt"
	"github.com/foxcpp/ttyprompt/terminal"
)

func eq(a, b []byte) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

func askForPasswd(tty *TTY, opts pinentry.Settings) (string, *assuan.Error) {
	dopts := prompt.DialogSettings{opts.Title, opts.Desc, opts.Prompt, opts.Opts.InvisibleChar}
	orig := dopts
	var firstbuf *memguard.LockedBuffer
	firstlen := 0
	for {
		buf, n, err := prompt.AskForPassword(tty.file, tty.num, dopts)

		if err != nil {
			if err.Error() == "AskForPassword: prompt rejected" {
				return "", &assuan.Error{
					assuan.ErrSrcPinentry, assuan.ErrCanceled,
					"pinentry", err.Error(),
				}
			}
			return "", &assuan.Error{
				assuan.ErrSrcPinentry, assuan.ErrAssuanServerFault,
				"pinentry", err.Error(),
			}
		}

		if len(opts.RepeatPrompt) != 0 {
			if firstbuf != nil {
				if eq(buf.Buffer()[:n], firstbuf.Buffer()[:firstlen]) {
					return string(buf.Buffer()[:n]), nil
				}

				dopts = orig
				dopts.Description = "\n\n" + opts.RepeatError

				firstbuf = nil
				firstlen = 0
				continue
			}

			dopts = orig
			dopts.Prompt = opts.RepeatPrompt
			firstbuf = buf
			firstlen = n
		} else {
			return string(buf.Buffer()[:n]), nil
		}
	}
}

func getPIN(tty *TTY, opts pinentry.Settings) (string, *assuan.Error) {
	defer tty.file.WriteString(terminal.TermClear + terminal.TermReset)
	if len(opts.Error) != 0 {
		opts.Desc += "\n\nERROR: " + opts.Error
	}
	if len(opts.Prompt) == 0 {
		opts.Prompt = "Enter PIN:"
	}
	if len(opts.Title) == 0 {
		opts.Title = "pinentry mode"
	}
	if len(opts.RepeatError) == 0 {
		opts.RepeatError = "Passwords do not match"
	}
	if len(opts.Opts.InvisibleChar) == 0 {
		opts.Opts.InvisibleChar = "*"
	}

	return askForPasswd(tty, opts)
}

func confirm(tty *TTY, opts pinentry.Settings) (bool, *assuan.Error) {
	defer tty.file.WriteString(terminal.TermClear + terminal.TermReset)
	if len(opts.Error) != 0 {
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

	res, err := prompt.AskToConfirm(tty.file, tty.num, prompt.DialogSettings{
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
	defer tty.file.WriteString(terminal.TermClear + terminal.TermReset)
	if len(opts.Error) != 0 {
		opts.Desc += "\n\nERROR: " + opts.Error
	}
	if len(opts.Title) == 0 {
		opts.Title = "pinentry mode"
	}

	err := prompt.ShowMessage(tty.file, tty.num, prompt.DialogSettings{
		Title:       opts.Title,
		Description: opts.Desc,
		Prompt:      opts.Prompt,
	})

	if err != nil {
		return &assuan.Error{
			assuan.ErrSrcPinentry, assuan.ErrAssuanServerFault,
			"pinentry", err.Error(),
		}
	}

	return nil
}

func pinentryMode(tty *TTY, flags settings, finishNotifyChan chan error) {
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
	log.Println("Accepting pinentry requests on stdin")
	err := pinentry.Serve(pinentry.Callbacks{getPINfunc, confirmFunc, msgFunc}, "ttyprompt v0.1.0")
	finishNotifyChan <- err
}
