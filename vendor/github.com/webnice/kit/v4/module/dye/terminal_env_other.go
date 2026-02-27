//go:build js || plan9 || aix
// +build js plan9 aix

package dye

func colorProfile() Profile {
	return ANSI256
}

func foregroundColor() Color {
	return ANSIColor(7)
}

func backgroundColor() Color {
	return ANSIColor(0)
}
