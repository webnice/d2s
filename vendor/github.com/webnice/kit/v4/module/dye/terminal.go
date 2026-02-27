package dye

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/webnice/kit/v4/module/dye/colorful"
	"github.com/webnice/kit/v4/module/dye/isatty"
)

func isTTY(fd uintptr) bool {
	if len(os.Getenv("CI")) > 0 {
		return false
	}

	return isatty.IsTerminal(fd)
}

// Определение поддерживаемого цвета терминала.
func getColorProfile() terminalProfile {
	if !isTTY(os.Stdout.Fd()) {
		return terminalAscii
	}

	return colorProfile()
}

func cliColorForced() bool {
	if forced := os.Getenv(envCliColorForce); forced != "" {
		return forced != "0"
	}

	return false
}

func xTermColor(s string) (colorRgb, error) {
	if len(s) < 24 || len(s) > 25 {
		return colorRgb(""), errInvalidColor
	}

	switch {
	case strings.HasSuffix(s, "\a"):
		s = strings.TrimSuffix(s, "\a")
	case strings.HasSuffix(s, "\033"):
		s = strings.TrimSuffix(s, "\033")
	case strings.HasSuffix(s, "\033\\"):
		s = strings.TrimSuffix(s, "\033\\")
	default:
		return colorRgb(""), errInvalidColor
	}

	s = s[4:]

	prefix := prefixRgb
	if !strings.HasPrefix(s, prefix) {
		return colorRgb(""), errInvalidColor
	}
	s = strings.TrimPrefix(s, prefix)

	h := strings.Split(s, "/")
	hex := fmt.Sprintf("#%s%s%s", h[0][:2], h[1][:2], h[2][:2])
	return colorRgb(hex), nil
}

func ansi256ToANSIColor(c colorAnsi256) colorAnsi {
	var r int
	md := math.MaxFloat64

	h, _ := colorful.Hex(colorAnsiHex[c])
	for i := 0; i <= 15; i++ {
		hb, _ := colorful.Hex(colorAnsiHex[i])
		d := h.DistanceHSLuv(hb)

		if d < md {
			md = d
			r = i
		}
	}

	return colorAnsi(r)
}

func hexToANSI256Color(c colorful.Color) colorAnsi256 {
	v2ci := func(v float64) int {
		if v < 48 {
			return 0
		}
		if v < 115 {
			return 1
		}
		return int((v - 35) / 40)
	}

	// Calculate the nearest 0-based color index at 16..231
	r := v2ci(c.R * 255.0) // 0..5 each
	g := v2ci(c.G * 255.0)
	b := v2ci(c.B * 255.0)
	ci := 36*r + 6*g + b /* 0..215 */

	// Calculate the represented colors back from the index
	i2cv := [6]int{0, 0x5f, 0x87, 0xaf, 0xd7, 0xff}
	cr := i2cv[r] // r/g/b, 0..255 each
	cg := i2cv[g]
	cb := i2cv[b]

	// Calculate the nearest 0-based gray index at 232..255
	var grayIdx int
	average := (r + g + b) / 3
	if average > 238 {
		grayIdx = 23
	} else {
		grayIdx = (average - 3) / 10 // 0..23
	}
	gv := 8 + 10*grayIdx // same value for r/g/b, 0..255

	// Return the one which is nearer to the original input rgb value
	c2 := colorful.Color{R: float64(cr) / 255.0, G: float64(cg) / 255.0, B: float64(cb) / 255.0}
	g2 := colorful.Color{R: float64(gv) / 255.0, G: float64(gv) / 255.0, B: float64(gv) / 255.0}
	colorDist := c.DistanceHSLuv(c2)
	grayDist := c.DistanceHSLuv(g2)

	if colorDist <= grayDist {
		return colorAnsi256(16 + ci)
	}
	return colorAnsi256(232 + grayIdx)
}
