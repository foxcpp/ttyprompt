// +build nomlock

package main

import (
	"fmt"
	"os"
)

func lockProcMem() {
	fmt.Fprintln(os.Stderr, "WARNING: Running without process memory lock. Targeted attack may lead to passphrase leak!")
}
