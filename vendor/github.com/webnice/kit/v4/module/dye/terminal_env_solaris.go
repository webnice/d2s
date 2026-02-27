//go:build solaris || illumos
// +build solaris illumos

package dye

import "golang.org/x/sys/unix"

func isForeground(fd int) (ret bool) {
	var (
		err     error
		pgrp, g int
	)

	if pgrp, err = unix.IoctlGetInt(fd, unix.TIOCGPGRP); err != nil {
		return
	}
	if g, err = unix.Getpgrp(); err != nil {
		return
	}
	ret = pgrp == g

	return
}
