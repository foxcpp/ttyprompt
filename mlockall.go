// +build !nomlock

package main

import (
	"fmt"
	"log"
	"os"
	"syscall"
)

func lockProcMem() {
	log.Println("Locking process memory...")
	if err := syscall.Mlockall(syscall.MCL_CURRENT | syscall.MCL_FUTURE); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to lock ttyprompt's memory. Refusing to continue.")
		os.Exit(3)
	}
}
