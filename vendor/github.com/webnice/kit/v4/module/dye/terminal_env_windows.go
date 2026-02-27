//go:build windows
// +build windows

package dye

import (
	"os"
	"strconv"

	"golang.org/x/sys/windows"
)

func colorProfile() terminalProfile {
	var (
		err                     error
		winVersion, buildNumber uint32
		conVersion              string
		cv                      int64
	)

	if os.Getenv(envAnsiConEmu) == keyOn {
		return terminalTrueColor
	}
	winVersion, _, buildNumber = windows.RtlGetNtVersionNumbers()
	if buildNumber < 10586 || winVersion < 10 {
		if os.Getenv(envAnsiCon) != "" {
			conVersion = os.Getenv(envAnsiConVersion)
			if cv, err = strconv.ParseInt(conVersion, 10, 64); err != nil || cv < 181 {
				return terminalANSI
			}
			return terminalANSI256
		}
		return terminalAscii
	}
	if buildNumber < 14931 {
		return terminalANSI256
	}

	return terminalTrueColor
}

func foregroundColor() Color { return ANSIColor(7) }

func backgroundColor() Color { return ANSIColor(0) }

// WindowsAnsiConsoleEnable включает обработку виртуальных терминалов в Windows.
// Это позволит использовать последовательности ANSI в консоли Windows.
// Возвращает исходный режим консоли и ошибку.
func WindowsAnsiConsoleEnable() (ret uint32, err error) {
	var (
		handle  windows.Handle
		vtpMode uint32
	)

	if handle, err = windows.GetStdHandle(windows.STD_OUTPUT_HANDLE); err != nil {
		return
	}
	if err = windows.GetConsoleMode(handle, &ret); err != nil {
		return
	}
	// Документация: https://docs.microsoft.com/en-us/windows/console/console-virtual-terminal-sequences
	if ret&windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING != windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING {
		vtpMode = ret | windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING
		if err = windows.SetConsoleMode(handle, vtpMode); err != nil {
			return
		}
	}

	return ret, nil
}

// WindowsAnsiConsoleRestore Восстанавливает режим консоли в предыдущее состояние.
func WindowsAnsiConsoleRestore(mode uint32) error {
	handle, err := windows.GetStdHandle(windows.STD_OUTPUT_HANDLE)
	if err != nil {
		return err
	}

	return windows.SetConsoleMode(handle, mode)
}
