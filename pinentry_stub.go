// +build nopinentry

package main

import (
	"errors"
	"fmt"
)

func pinentryMode(_ *TTY, _ settings, errC chan error) {
	fmt.Println("ERR 536870981 Not implemented <User defined source 1>")
	errC <- errors.New("Pinentry mode is not supported by this build.")
}
