package terminal

import (
	"errors"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"

	"github.com/awnumar/memguard"
)

type DialogSettings struct {
	Title       string
	Description string
	Prompt      string
}

var FirstTty = -1

/*
TurnOnRawIO sets flags suitable for raw I/O (no echo, per-character input, etc)
and returns original flags.
*/
func TurnOnRawIO(tty *os.File) (orig Termios, err error) {
	log.Println("Turning on raw I/O on", tty.Name())

	termios, err := TcGetAttr(tty.Fd())
	if err != nil {
		return Termios{}, errors.New("TurnOnRawIO: failed to get flags: " + err.Error())
	}
	termiosOrig := *termios

	termios.Lflag &^= syscall.ECHO
	termios.Lflag &^= syscall.ICANON
	termios.Iflag &^= syscall.IXON
	termios.Iflag |= syscall.IUTF8
	err = TcSetAttr(tty.Fd(), termios)
	if err != nil {
		return Termios{}, errors.New("TurnOnRawIO: flags to set flags: " + err.Error())
	}
	return termiosOrig, nil
}

func UnlockTTY(tty *os.File) {
	log.Println("Setting", tty.Name(), "access mode to 0620...")
	err := os.Chmod(tty.Name(), 0620)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to recover old access mode on", tty.Name(), ":", err)
	}
}

func LockTTY(tty *os.File) {
	log.Println("Setting", tty.Name(), "access mode to 0000...")
	err := os.Chmod(tty.Name(), 0000)
	if err != nil {
		panic("failed to recover old access mode on " + tty.Name() + ": " + err.Error())
	}
}

/*
ReadPassword configures TTY to non-canonical mode and reads password
byte-by-byte showing '*' for each character.
*/
func ReadPassword(tty *os.File, output []byte) ([]byte, error) {
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
		_, err = tty.WriteString("*")
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

	if err := SwitchVTThrough(outputTty.Fd(), FirstTty); err != nil {
		fmt.Fprintln(os.Stderr, "failed to switch TTY back:", err)
		outputTty.WriteString("\nOops! We can't switch TTYs. Do it manually (i.e. Ctrl+Alt+F1).\n")
		time.Sleep(time.Second * 5)
		return err
	}
	return nil
}

func WritePrompt(tty *os.File, settings DialogSettings) error {
	fullPrompt := TermClear + TermReset
	fullPrompt += "ttyprompt v0.1 | " + settings.Title + "\n"
	fullPrompt += "================================================================================\n"
	fullPrompt += "[" + time.Now().String() + "]\n"
	fullPrompt += "\n"
	fullPrompt += settings.Description
	fullPrompt += "\n"
	fullPrompt += settings.Prompt + " "
	_, err := tty.WriteString(fullPrompt)
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
	firsttty, err := CurrentVT()
	if err != nil {
		return nil, 0, err
	}

	if err := SwitchVTThrough(tty.Fd(), ttyNum); err != nil {
		return nil, 0, err
	}
	defer switchToOriginalVT(tty, firsttty)

	settings.Description += "\n\nJust press <Enter> to reject request."
	if WritePrompt(tty, settings); err != nil {
		return nil, 0, err
	}

	origTermios, err := TurnOnRawIO(tty)
	if err != nil {
		return nil, 0, err
	}
	defer TcSetAttr(tty.Fd(), &origTermios)

	bufHandle, err := memguard.NewMutable(2048)
	if err != nil {
		return nil, 0, errors.New("AskForPassword: " + err.Error())
	}

	slice, err := ReadPassword(tty, bufHandle.Buffer())
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
	firsttty, err := CurrentVT()
	if err != nil {
		return false, err
	}

	if err := SwitchVTThrough(tty.Fd(), ttyNum); err != nil {
		return false, err
	}
	defer switchToOriginalVT(tty, firsttty)

	if WritePrompt(tty, settings); err != nil {
		return false, err
	}

	origTermios, err := TurnOnRawIO(tty)
	if err != nil {
		return false, err
	}
	defer TcSetAttr(tty.Fd(), &origTermios)

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
	firsttty, err := CurrentVT()
	if err != nil {
		return err
	}

	if err := SwitchVTThrough(tty.Fd(), ttyNum); err != nil {
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
