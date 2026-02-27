package dye

import (
	"fmt"
	"strings"
)

// New Конструктор объекта пакета, возвращается интерфейс пакета.
// Если переданы последовательности, тогда они добавляются "до".
func New(seq ...string) Interface {
	var dye = &impl{profile: getColorProfile()}

	dye.Add(seq...)

	return dye
}

// Add Добавление последовательности в конец текущей последовательности.
func (dye *impl) Add(seq ...string) Interface { dye.seq = append(dye.seq, seq...); return dye }

// Source Возвращает текущую последовательность в исходном виде как срез строк.
func (dye *impl) Source() []string { return dye.seq }

// Done Завершение и возврат последовательности.
func (dye *impl) Done() Sequence {
	var ret string

	ret = strings.Join(dye.seq, ";")
	ret = fmt.Sprintf("%s%sm", seqCSI, ret)

	return sequence(ret)
}

// Bold Добавление в последовательность переключения текста "Жирный".
func (dye *impl) Bold() Interface { dye.seq = append(dye.seq, seqBold); return dye }

// Faded Добавление в последовательность переключения текста "Блёклый".
func (dye *impl) Faded() Interface { dye.seq = append(dye.seq, seqFaded); return dye }

// Italic Добавление в последовательность переключения текста "Курсив".
func (dye *impl) Italic() Interface { dye.seq = append(dye.seq, seqItalic); return dye }

// Underline Добавление в последовательность переключения текста "Подчёркнутый один раз".
func (dye *impl) Underline() Interface { dye.seq = append(dye.seq, seqUnderline); return dye }

// Reverse Добавление в последовательность переключения текста "Инвертирование цвета".
func (dye *impl) Reverse() Interface { dye.seq = append(dye.seq, seqReverse); return dye }

// CrossOut Добавление в последовательность переключения текста "Зачёркнутый".
func (dye *impl) CrossOut() Interface { dye.seq = append(dye.seq, seqCrossOut); return dye }

// ResetBold Добавление в последовательность переключения текста "Сбросить жирный".
func (dye *impl) ResetBold() Interface { dye.seq = append(dye.seq, seqResetBold); return dye }

// ResetFaded Добавление в последовательность переключения текста "Сбросить блёклый".
func (dye *impl) ResetFaded() Interface { dye.seq = append(dye.seq, seqResetFaded); return dye }

// ResetItalic Добавление в последовательность переключения текста "Сбросить курсив".
func (dye *impl) ResetItalic() Interface { dye.seq = append(dye.seq, seqResetItalic); return dye }

// ResetUnderline Добавление в последовательность переключения текста "Сбросить подчёркнутый один раз".
func (dye *impl) ResetUnderline() Interface { dye.seq = append(dye.seq, seqResetUnderline); return dye }

// ResetReverse Добавление в последовательность переключения текста "Сбросить инвертирование цвета".
func (dye *impl) ResetReverse() Interface { dye.seq = append(dye.seq, seqResetReverse); return dye }

// ResetCrossOut Добавление в последовательность переключения текста "Сбросить зачёркнутый".
func (dye *impl) ResetCrossOut() Interface { dye.seq = append(dye.seq, seqResetCrossOut); return dye }

// Reset Добавление в последовательность "Сброс", переход в нормальный режим.
func (dye *impl) Reset() Interface { dye.seq = append(dye.seq, seqReset); return dye }

// Background Флаг, говорящий что следующий цвет будет цветом фона, только для ANSI цветов.
func (dye *impl) Background() Interface { dye.isBackground = true; return dye }

// Foreground Флаг, говорящий что следующий цвет будет цветом текста, только для ANSI цветов, по умолчанию.
func (dye *impl) Foreground() Interface { dye.isBackground = false; return dye }

// Bright Флаг, указывающий что следующий цвет будет ярким, только для ANSI цветов.
func (dye *impl) Bright() Interface { dye.isBright = true; return dye }

// Normal Флаг, указывающий что следующий цвет будет нормальной яркости, только для ANSI цветов.
func (dye *impl) Normal() Interface { dye.isBright = false; return dye }
