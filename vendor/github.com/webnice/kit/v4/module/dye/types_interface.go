package dye

// Interface Интерфейс пакета.
type Interface interface {
	// Add Добавление последовательности в конец текущей последовательности.
	Add(seq ...string) Interface

	// Source Возвращает текущую последовательность в исходном виде как срез строк.
	Source() []string

	// Done Завершение и возврат последовательности.
	Done() Sequence

	// ПЕРЕКЛЮЧАТЕЛИ.

	// Bold Добавление в последовательность переключения текста "Жирный".
	Bold() Interface

	// Faded Добавление в последовательность переключения текста "Блёклый".
	Faded() Interface

	// Italic Добавление в последовательность переключения текста "Курсив".
	Italic() Interface

	// Underline Добавление в последовательность переключения текста "Подчёркнутый один раз".
	Underline() Interface

	// Reverse Добавление в последовательность переключения текста "Инвертирование цвета".
	Reverse() Interface

	// CrossOut Добавление в последовательность переключения текста "Зачёркнутый".
	CrossOut() Interface

	// ResetBold Добавление в последовательность переключения текста "Сбросить жирный".
	ResetBold() Interface

	// ResetFaded Добавление в последовательность переключения текста "Сбросить блёклый".
	ResetFaded() Interface

	// ResetItalic Добавление в последовательность переключения текста "Сбросить курсив".
	ResetItalic() Interface

	// ResetUnderline Добавление в последовательность переключения текста "Сбросить подчёркнутый один раз".
	ResetUnderline() Interface

	// ResetReverse Добавление в последовательность переключения текста "Сбросить инвертирование цвета".
	ResetReverse() Interface

	// ResetCrossOut Добавление в последовательность переключения текста "Сбросить зачёркнутый".
	ResetCrossOut() Interface

	// Reset Добавление в последовательность "Сброс", переход в нормальный режим.
	Reset() Interface

	// ОПЦИИ.

	// Background Флаг, указывающий что следующий цвет будет цветом фона.
	Background() Interface

	// Foreground Флаг, указывающий что следующий цвет будет цветом текста.
	Foreground() Interface

	// Bright Флаг, указывающий что следующий цвет будет ярким.
	Bright() Interface

	// Normal Флаг, указывающий что следующий цвет будет нормальной яркости.
	Normal() Interface

	// ЦВЕТА.

	// Black Добавление чёрного ANSI цвета в последовательность.
	Black() Interface

	// Red Добавление красного ANSI цвета в последовательность.
	Red() Interface

	// Green Добавление зелёного ANSI цвета в последовательность.
	Green() Interface

	// Yellow Добавление жёлтого ANSI цвета в последовательность.
	Yellow() Interface

	// Blue Добавление синего ANSI цвета в последовательность.
	Blue() Interface

	// Magenta Добавление пурпурного ANSI цвета в последовательность.
	Magenta() Interface

	// Cyan Добавление бирюзового ANSI цвета в последовательность.
	Cyan() Interface

	// White Добавление белого ANSI цвета в последовательность.
	White() Interface

	// HexRgb Конвертация HEX RGB цвета в формате #000000 в ANSI цвет и добавление в последовательность.
	HexRgb(hexRgb string) (err error)
}

// Sequence Интерфейс готовой последовательности.
type Sequence interface {
	// String Возвращение последовательности в виде строки.
	String() string

	// Byte Возвращение последовательности в виде среза байт.
	Byte() []byte
}
