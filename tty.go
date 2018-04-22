package main

import (
	"os"
	"strconv"

	"github.com/foxcpp/ttyprompt/terminal"
)

// TTY struct is a thin tuple for resources and information associated with virtual terminal.
type TTY struct {
	file *os.File
	num  int
	path string
}

func (t *TTY) free() {
	t.file.WriteString(terminal.TermClear)
	t.file.WriteString(terminal.TermReset)

	t.file.Close()
}

func getTTY(num int) (res *TTY, err error) {
	res = new(TTY)
	res.num = num
	res.path = "/dev/tty" + strconv.Itoa(num)

	res.file, err = os.OpenFile(res.path, os.O_RDWR, os.FileMode(0))
	if err != nil {
		return
	}

	res.file.WriteString("ttyprompt acquired this TTY\n")

	return
}
