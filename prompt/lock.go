package prompt

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

func UnlockTTY(tty *os.File) {
	log.Println("Unlocking", tty.Name()+"...")
	syscall.FcntlFlock(tty.Fd(), syscall.F_UNLCK, &syscall.Flock_t{
		Type: syscall.F_WRLCK,
		Pid:  int32(os.Getpid()),
	})
}

func LockTTY(tty *os.File) {
	log.Println("Locking", tty.Name()+"...")
	err := syscall.FcntlFlock(tty.Fd(), syscall.F_SETLK, &syscall.Flock_t{
		Type: syscall.F_WRLCK,
		Pid:  int32(os.Getpid()),
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, "Another ttyprompt instance is using requested TTY prompt. Waiting for it to finish...")
	}
	syscall.FcntlFlock(tty.Fd(), syscall.F_SETLKW, &syscall.Flock_t{
		Type: syscall.F_WRLCK,
		Pid:  int32(os.Getpid()),
	})
}
