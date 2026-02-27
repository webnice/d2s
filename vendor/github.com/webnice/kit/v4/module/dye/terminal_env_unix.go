//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris
// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package dye

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

func colorProfile() terminalProfile {
	var (
		term      string
		colorTerm string
	)

	term, colorTerm = os.Getenv(envTerm), os.Getenv(envColorTerm)
	switch strings.ToLower(colorTerm) {
	case key24Bit:
		fallthrough
	case keyTrueColor:
		if strings.HasPrefix(term, keyScreen) {
			// Терминал tmux поддерживает TrueColor, экран только ANSI256
			if os.Getenv(envTermProgram) != keyTmux {
				return terminalANSI256
			}
		}
		return terminalTrueColor
	case keyYes:
		fallthrough
	case keyTrue:
		return terminalANSI256
	}
	switch term {
	case keyXtermKitty:
		return terminalTrueColor
	case keyLinux:
		return terminalANSI
	}
	if strings.Contains(term, key256Color) {
		return terminalANSI256
	}
	if strings.Contains(term, keyColor) {
		return terminalANSI
	}
	if strings.Contains(term, keyAnsi) {
		return terminalANSI
	}

	return terminalAscii
}

func foregroundColor() (ret Color) {
	var (
		err       error
		s         string
		c         []string
		colorFGBG string
		i         int
	)

	if s, err = termStatusReport(10); err == nil {
		if ret, err = xTermColor(s); err == nil {
			return
		}
	}
	colorFGBG = os.Getenv(envColorFgbg)
	if strings.Contains(colorFGBG, ";") {
		c = strings.Split(colorFGBG, ";")
		if i, err = strconv.Atoi(c[0]); err == nil {
			ret = colorAnsi(i)
			return
		}
	}
	ret = colorAnsi(7)

	return
}

func backgroundColor() (ret Color) {
	var (
		err       error
		s         string
		c         []string
		colorFGBG string
		i         int
	)

	if s, err = termStatusReport(11); err == nil {
		if ret, err = xTermColor(s); err == nil {
			return
		}
	}
	colorFGBG = os.Getenv(envColorFgbg)
	if strings.Contains(colorFGBG, ";") {
		c = strings.Split(colorFGBG, ";")
		if i, err = strconv.Atoi(c[len(c)-1]); err == nil {
			ret = colorAnsi(i)
			return
		}
	}
	ret = colorAnsi(0)

	return
}

func waitForData(fd uintptr, timeout time.Duration) (err error) {
	var (
		tv      unix.Timeval
		readFds unix.FdSet
		n       int
	)

	tv = unix.NsecToTimeval(int64(timeout))
	readFds.Set(int(fd))
	for {
		if n, err = unix.Select(int(fd)+1, &readFds, nil, nil, &tv); err == unix.EINTR {
			continue
		}
		if err != nil {
			return
		}
		if n == 0 {
			err = errTimeout
			return
		}
		break
	}

	return
}

func readNextByte(f *os.File) (ret byte, err error) {
	var (
		b [1]byte
		n int
	)

	if err = waitForData(f.Fd(), timeoutOSC); err != nil {
		return
	}
	if n, err = f.Read(b[:]); err != nil {
		return
	}
	if n == 0 {
		err = errReturnedNoData
	}
	ret = b[0]

	return
}

// Читает либо ответ OSC, либо ответ позиции курсора:
//   - OSC ответ: "\x1b]11;rgb:1111/1111/1111\x1b\\"
//   - Ответ позиции курсора: "\x1b[42;1R"
func readNextResponse(fd *os.File) (response string, isOSC bool, err error) {
	const (
		squareBracketOpen, squareBracketClose = '[', ']'
		keyBel, keyR, key033                  = '\a', 'R', '\033'
	)
	var (
		start, tpe, b byte
		oscResponse   bool
	)

	if start, err = readNextByte(fd); err != nil {
		return
	}
	for start != key033 {
		if start, err = readNextByte(fd); err != nil {
			return
		}
	}
	response += string(start)
	if tpe, err = readNextByte(fd); err != nil {
		return
	}
	switch response += string(tpe); tpe {
	case squareBracketOpen:
		oscResponse = false
	case squareBracketClose:
		oscResponse = true
	default:
		err = errStatusReport
		return
	}
	for {
		if b, err = readNextByte(os.Stdout); err != nil {
			return
		}
		switch response += string(b); oscResponse {
		case true:
			if b == keyBel || strings.HasSuffix(response, string(key033)) {
				isOSC = true
				return
			}
		default:
			if b == keyR {
				return
			}
		}
		if len(response) > 25 {
			break
		}
	}
	err = errStatusReport

	return
}

func termStatusReport(sequence int) (ret string, err error) {
	var (
		term   string
		noEcho unix.Termios
		t      *unix.Termios
		isOSC  bool
	)

	if term = os.Getenv(envTerm); strings.HasPrefix(term, keyScreen) {
		err = errStatusReport
		return
	}
	if !isForeground(unix.Stdout) {
		err = errStatusReport
		return
	}
	if t, err = unix.IoctlGetTermios(unix.Stdout, tcgetattr); err != nil {
		err = errStatusReport
		return
	}
	defer func() { _ = unix.IoctlSetTermios(unix.Stdout, tcsetattr, t) }()
	noEcho = *t
	noEcho.Lflag = noEcho.Lflag &^ unix.ECHO
	noEcho.Lflag = noEcho.Lflag &^ unix.ICANON
	if err = unix.IoctlSetTermios(unix.Stdout, tcsetattr, &noEcho); err != nil {
		err = errStatusReport
		return
	}
	fmt.Printf("\033]%d;?\033\\", sequence)
	fmt.Printf("\033[6n")
	if ret, isOSC, err = readNextResponse(os.Stdout); err != nil {
		return
	}
	if !isOSC {
		err = errStatusReport
		return
	}
	if _, _, err = readNextResponse(os.Stdout); err != nil {
		return
	}

	return
}
