package prompt

import (
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
	syscall.FcntlFlock(tty.Fd(), syscall.F_SETLKW, &syscall.Flock_t{
		Type: syscall.F_WRLCK,
		Pid:  int32(os.Getpid()),
	})
}
