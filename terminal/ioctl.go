package terminal

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

const (
	vtActivate uintptr = 0x5606
)

func SwitchVT(num int) error {
	log.Println("Switching to VT", num, "using current VT")

	// Workaround (nouveau?) bug: if we switch VT's too fast we will crash X.Org.
	time.Sleep(250 * time.Millisecond)

	tty, err := CurrentVT()
	if err != nil {
		return errors.New("SwitchVT: " + err.Error())
	}
	ttyf, err := os.Open("/dev/tty" + strconv.Itoa(tty))
	if err != nil {
		return errors.New("SwitchVT: " + err.Error())
	}

	_, _, errnop := syscall.Syscall(syscall.SYS_IOCTL, ttyf.Fd(), vtActivate, uintptr(num))
	errno := syscall.Errno(errnop)
	if errno != 0 {
		return errors.New("SwitchVT: ioctl: " + errno.Error())
	}
	return nil
}

// SwitchVTThrough is same as SwitchVT but allows you to specify TTY descriptor (you may use any TTY).
func SwitchVTThrough(fd uintptr, num int) error {
	log.Println("Switching to VT", num, "using FD", fd)

	// Workaround (nouveau?) bug: if we switch VT's too fast we will crash X.Org.
	time.Sleep(250 * time.Millisecond)

	_, _, errnop := syscall.Syscall(syscall.SYS_IOCTL, fd, vtActivate, uintptr(num))
	errno := syscall.Errno(errnop)
	if errno != 0 {
		return errors.New("SwitchVT: ioctl: " + errno.Error())
	}
	return nil
}

func CurrentVT() (int, error) {
	file, err := ioutil.ReadFile("/sys/devices/virtual/tty/tty0/active")
	if err != nil {
		return 0, errors.New("CurrentVT: " + err.Error())
	}
	filestr := string(file)

	if !strings.HasPrefix(filestr, "tty") {
		return 0, errors.New("CurrentVT: invalid active vt name: " + filestr)
	}

	i, err := strconv.Atoi(strings.TrimSuffix(filestr[3:], "\n"))
	if err != nil {
		return 0, errors.New("CurrentVT: invalid active vt name: " + filestr)
	}

	return i, nil
}
