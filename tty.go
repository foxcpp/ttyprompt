package main

import (
	"fmt"
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
	err := os.Chmod(t.path, 0620)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to recover old access mode on", t.path, ":", err)
	}

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

	os.Chmod(res.path, 0000)
	res.file.WriteString("ttyprompt acquired this TTY\n")

	return
}
