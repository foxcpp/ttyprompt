package terminal

import (
	"fmt"
	"log"
	"os"
)

const (
	lockedMode  os.FileMode = 0000
	defaultMode             = 620
)

func UnlockTTY(tty *os.File) {
	info, _ := tty.Stat()
	if info.Mode()&os.ModePerm != lockedMode {
		return
	}
	log.Println("Setting", tty.Name(), "access mode to 0620...")
	err := os.Chmod(tty.Name(), 0620)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to recover old access mode on", tty.Name(), ":", err)
	}
}

func LockTTY(tty *os.File) {
	info, _ := tty.Stat()
	if info.Mode()&os.ModePerm != defaultMode {
		log.Println("Skipping chmod lock because of non-default permissions...")
		return
	}

	log.Println("Setting", tty.Name(), "access mode to 0000...")
	err := os.Chmod(tty.Name(), 0000)
	if err != nil {
		panic("failed to recover old access mode on " + tty.Name() + ": " + err.Error())
	}
}
