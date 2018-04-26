package prompt

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/awnumar/memguard"
	"github.com/foxcpp/ttyprompt/terminal"
)

type DialogSettings struct {
	Title       string
	Description string
	Prompt      string
	PassChar    string
}

var FirstTty = -1

/*
ReadPassword configures TTY to non-canonical mode and reads password
byte-by-byte showing '*' for each character.
*/
func ReadPassword(tty *os.File, output []byte, echoChar string) ([]byte, error) {
	log.Println("Reading password...")
	cursor := output[0:1]
	readen := 0
	for {
		n, err := tty.Read(cursor)
		if n != 1 {
			return nil, errors.New("ReadPassword: invalid read size when not in canonical mode")
		}
		if err != nil {
			return nil, errors.New("ReadPassword: " + err.Error())
		}
		if cursor[0] == '\n' {
			break
		}
		if cursor[0] == '\x7F' /* DEL */ {
			if readen != 0 {
				_, err := tty.WriteString("\b \b")
				if err != nil {
					return nil, errors.New("ReadPassword: " + err.Error())
				}
				readen--
				cursor = output[readen : readen+1]
			}
			continue
		}
		_, err = tty.WriteString(echoChar)
		if err != nil {
			return nil, errors.New("ReadPassword: " + err.Error())
		}
		readen++
		cursor = output[readen : readen+1]
	}

	return output[0:readen], nil
}

func switchToOriginalVT(outputTty *os.File, num int) error {
	if FirstTty == -1 {
		FirstTty = num
	}

	if err := terminal.SwitchVTThrough(outputTty.Fd(), FirstTty); err != nil {
		fmt.Fprintln(os.Stderr, "failed to switch TTY back:", err)
		outputTty.WriteString("\nOops! We can't switch TTYs. Do it manually (i.e. Ctrl+Alt+F1).\n")
		time.Sleep(time.Second * 5)
		return err
	}
	return nil
}

func WritePrompt(tty *os.File, settings DialogSettings) error {
	ctx, err := CaptureExeCtx()
	if err != nil {
		return err
	}

	fullPrompt := terminal.TermClear + terminal.TermReset
	fullPrompt += "ttyprompt v0.1 | " + settings.Title + "\n"
	fullPrompt += "\n"
	fullPrompt += ctx.String()
	fullPrompt += "\n\n"
	fullPrompt += settings.Description
	fullPrompt += "\n"
	fullPrompt += settings.Prompt + " "
	_, err = tty.WriteString(fullPrompt)
	if err != nil {
		return errors.New("WritePrompt: " + err.Error())
	}
	return nil
}

/*
AskForPassword does everything needed to get a password from user through specified TTY.

Returned values are: Pointer to password buffer, length of password, error if any.
*/
func AskForPassword(tty *os.File, ttyNum int, settings DialogSettings) (*memguard.LockedBuffer, int, error) {
	log.Println("Requesting password...")
	firsttty, err := terminal.CurrentVT()
	if err != nil {
		return nil, 0, err
	}

	if err := terminal.SwitchVTThrough(tty.Fd(), ttyNum); err != nil {
		return nil, 0, err
	}
	defer switchToOriginalVT(tty, firsttty)

	settings.Description += "\n\nJust press <Enter> to reject request."
	if WritePrompt(tty, settings); err != nil {
		return nil, 0, err
	}

	origTermios, err := terminal.TurnOnRawIO(tty)
	if err != nil {
		return nil, 0, err
	}
	defer terminal.TcSetAttr(tty.Fd(), &origTermios)

	bufHandle, err := memguard.NewMutable(2048)
	if err != nil {
		return nil, 0, errors.New("AskForPassword: " + err.Error())
	}

	slice, err := ReadPassword(tty, bufHandle.Buffer(), settings.PassChar)
	if err != nil {
		return nil, 0, err
	}
	if len(slice) == 0 {
		return nil, 0, errors.New("AskForPassword: prompt rejected")
	}
	return bufHandle, len(slice), nil
}

func AskToConfirm(tty *os.File, ttyNum int, settings DialogSettings) (bool, error) {
	log.Println("Requesting confirmation...")
	firsttty, err := terminal.CurrentVT()
	if err != nil {
		return false, err
	}

	if err := terminal.SwitchVTThrough(tty.Fd(), ttyNum); err != nil {
		return false, err
	}
	defer switchToOriginalVT(tty, firsttty)

	if WritePrompt(tty, settings); err != nil {
		return false, err
	}

	origTermios, err := terminal.TurnOnRawIO(tty)
	if err != nil {
		return false, err
	}
	defer terminal.TcSetAttr(tty.Fd(), &origTermios)

	chr := make([]byte, 1)
	if tty.Read(chr); err != nil {
		return false, err
	}

	if chr[0] == 'Y' || chr[0] == 'y' {
		return true, nil
	}

	return false, nil
}

func ShowMessage(tty *os.File, ttyNum int, settings DialogSettings) error {
	log.Println("Showing message...")
	firsttty, err := terminal.CurrentVT()
	if err != nil {
		return err
	}

	if err := terminal.SwitchVTThrough(tty.Fd(), ttyNum); err != nil {
		return err
	}
	defer switchToOriginalVT(tty, firsttty)

	settings.Prompt = "Press any key."
	if WritePrompt(tty, settings); err != nil {
		return err
	}

	chr := make([]byte, 1)
	if tty.Read(chr); err != nil {
		return err
	}
	return nil
}
