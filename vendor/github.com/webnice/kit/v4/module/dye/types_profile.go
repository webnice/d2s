package dye

import (
	"image/color"
	"strconv"
	"strings"

	"github.com/webnice/kit/v4/module/dye/colorful"
)

// Профиль терминала.
type terminalProfile int

func (tp terminalProfile) Convert(c Color) Color {
	if tp == terminalAscii {
		return colorNone{}
	}

	switch v := c.(type) {
	case colorAnsi:
		return v

	case colorAnsi256:
		if tp == terminalANSI {
			return ansi256ToANSIColor(v)
		}
		return v

	case colorRgb:
		h, err := colorful.Hex(string(v))
		if err != nil {
			return nil
		}
		if tp < terminalTrueColor {
			ac := hexToANSI256Color(h)
			if tp == terminalANSI {
				return ansi256ToANSIColor(ac)
			}
			return ac
		}
		return v
	}

	return c
}

func (tp terminalProfile) Color(s string) Color {
	if len(s) == 0 {
		return nil
	}

	var c Color
	if strings.HasPrefix(s, keyHash) {
		c = colorRgb(s)
	} else {
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil
		}

		if i < 16 {
			c = colorAnsi(i)
		} else {
			c = colorAnsi256(i)
		}
	}

	return tp.Convert(c)
}

func (tp terminalProfile) FromColor(c color.Color) Color {
	col, _ := colorful.MakeColor(c)
	return tp.Color(col.Hex())
}
