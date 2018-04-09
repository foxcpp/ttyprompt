package terminal

import (
	"errors"
	"os"
	"syscall"
	"time"

	"github.com/awnumar/memguard"
)

/*
ReadPassword configures TTY to non-canonical mode and reads password
byte-by-byte showing '*' for each character.
*/
func ReadPassword(tty *os.File, output []byte) ([]byte, error) {
	termios, err := TcGetAttr(tty.Fd())
	if err != nil {
		return nil, errors.New("failed to configure TTY: " + err.Error())
	}
	termiosOrig := *termios
	defer TcSetAttr(tty.Fd(), &termiosOrig)

	termios.Lflag &^= syscall.ECHO
	termios.Lflag &^= syscall.ICANON
	termios.Iflag &^= syscall.IXON
	//termios.Iflag &^= syscall.IGNBRK
	//termios.Iflag &^= syscall.BRKINT
	termios.Iflag |= syscall.IUTF8
	err = TcSetAttr(tty.Fd(), termios)
	if err != nil {
		return nil, errors.New("ReadPassword: " + err.Error())
	}

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

/*
AskForPassword does everything needed to get a password from user through specified TTY.

Returned values are: Pointer to password buffer, length of password, error if any.
*/
func AskForPassword(tty *os.File, prompt string) (*memguard.LockedBuffer, int, error) {
	fullPrompt := TermClear + TermReset
	fullPrompt += "ttyprompt v0.1\n"
	fullPrompt += "################################################################################\n"
	fullPrompt += "[" + time.Now().String() + "]\n"
	fullPrompt += "\n"
	fullPrompt += prompt
	fullPrompt += "\n\n"
	fullPrompt += "Just press <Enter> to reject request.\n"
	fullPrompt += ": "
	_, err := tty.WriteString(fullPrompt)
	if err != nil {
		return nil, 0, errors.New("AskForPassword: " + err.Error())
	}

	bufHandle, err := memguard.NewMutable(2048)
	if err != nil {
		return nil, 0, errors.New("AskForPassword: " + err.Error())
	}

	slice, err := ReadPassword(tty, bufHandle.Buffer())
	if err != nil {
		return nil, 0, errors.New("AskForPassword: " + err.Error())
	}
	if len(slice) == 0 {
		return nil, 0, errors.New("AskForPassword: prompt rejected")
	}
	return bufHandle, len(slice), nil
}
