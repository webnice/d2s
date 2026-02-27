package dye

import "fmt"

// Black Добавление чёрного ANSI цвета в последовательность.
func (dye *impl) Black() Interface {
	switch dye.isBright {
	case true:
		dye.seq = append(dye.seq, colorANSIBrightBlack.Sequence(dye.isBackground))
	default:
		dye.seq = append(dye.seq, colorANSIBlack.Sequence(dye.isBackground))
	}
	dye.isBright, dye.isBackground = false, false

	return dye
}

// Red Добавление красного ANSI цвета в последовательность.
func (dye *impl) Red() Interface {
	switch dye.isBright {
	case true:
		dye.seq = append(dye.seq, colorANSIBrightRed.Sequence(dye.isBackground))
	default:
		dye.seq = append(dye.seq, colorANSIRed.Sequence(dye.isBackground))
	}
	dye.isBright, dye.isBackground = false, false

	return dye
}

// Green Добавление зелёного ANSI цвета в последовательность.
func (dye *impl) Green() Interface {
	switch dye.isBright {
	case true:
		dye.seq = append(dye.seq, colorANSIBrightGreen.Sequence(dye.isBackground))
	default:
		dye.seq = append(dye.seq, colorANSIGreen.Sequence(dye.isBackground))
	}
	dye.isBright, dye.isBackground = false, false

	return dye
}

// Yellow Добавление жёлтого ANSI цвета в последовательность.
func (dye *impl) Yellow() Interface {
	switch dye.isBright {
	case true:
		dye.seq = append(dye.seq, colorANSIBrightYellow.Sequence(dye.isBackground))
	default:
		dye.seq = append(dye.seq, colorANSIYellow.Sequence(dye.isBackground))
	}
	dye.isBright, dye.isBackground = false, false

	return dye
}

// Blue Добавление синего ANSI цвета в последовательность.
func (dye *impl) Blue() Interface {
	switch dye.isBright {
	case true:
		dye.seq = append(dye.seq, colorANSIBrightBlue.Sequence(dye.isBackground))
	default:
		dye.seq = append(dye.seq, colorANSIBlue.Sequence(dye.isBackground))
	}
	dye.isBright, dye.isBackground = false, false

	return dye
}

// Magenta Добавление пурпурного ANSI цвета в последовательность.
func (dye *impl) Magenta() Interface {
	switch dye.isBright {
	case true:
		dye.seq = append(dye.seq, colorANSIBrightMagenta.Sequence(dye.isBackground))
	default:
		dye.seq = append(dye.seq, colorANSIMagenta.Sequence(dye.isBackground))
	}
	dye.isBright, dye.isBackground = false, false

	return dye
}

// Cyan Добавление бирюзового ANSI цвета в последовательность.
func (dye *impl) Cyan() Interface {
	switch dye.isBright {
	case true:
		dye.seq = append(dye.seq, colorANSIBrightCyan.Sequence(dye.isBackground))
	default:
		dye.seq = append(dye.seq, colorANSICyan.Sequence(dye.isBackground))
	}
	dye.isBright, dye.isBackground = false, false

	return dye
}

// White Добавление белого ANSI цвета в последовательность.
func (dye *impl) White() Interface {
	switch dye.isBright {
	case true:
		dye.seq = append(dye.seq, colorANSIBrightWhite.Sequence(dye.isBackground))
	default:
		dye.seq = append(dye.seq, colorANSIWhite.Sequence(dye.isBackground))
	}
	dye.isBright, dye.isBackground = false, false

	return dye
}

// HexRgb Конвертация HEX RGB цвета в формате #000000 в ANSI цвет и добавление в последовательность.
func (dye *impl) HexRgb(hexRgb string) (err error) {
	var arr []string

	if arr = rexHexColor.FindStringSubmatch(hexRgb); len(arr) != 2 {
		err = fmt.Errorf(tplBadParam, hexRgb)
		return
	}
	dye.seq = append(dye.seq, dye.profile.Color(arr[1]).Sequence(dye.isBackground))

	return
}
