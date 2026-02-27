package dye

import (
	"errors"
	"time"
)

const (
	// Не цветной терминал.
	terminalAscii = terminalProfile(iota)

	// Терминал ANSI, цвет 4 бита.
	terminalANSI

	// Терминал ANSI256, цвет 8 бит.
	terminalANSI256

	// Терминал TrueColor, цвет 24 бита.
	terminalTrueColor
)

const (
	envCliColorForce  = "CLICOLOR_FORCE"
	envTerm           = "TERM"
	envColorTerm      = "COLORTERM"
	envTermProgram    = "TERM_PROGRAM"
	envColorFgbg      = "COLORFGBG"
	envAnsiCon        = "ANSICON"
	envAnsiConVersion = "ANSICON_VER"
	envAnsiConEmu     = "ConEmuANSI"
)

// Определение последовательностей.
const (
	seqCSI            = "\x1b[" // Начало последовательности управления.
	seqReset          = "0"     // Сброс, переход в нормальный режим.
	seqBold           = "1"     // Жирный.
	seqFaded          = "2"     // Блёклый.
	seqItalic         = "3"     // Курсив.
	seqUnderline      = "4"     // Подчёркнутый один раз.
	seqReverse        = "7"     // Инвертирование цвета.
	seqCrossOut       = "9"     // Зачёркнутый.
	seqResetBold      = "21"    // Сбросить жирный.
	seqResetFaded     = "22"    // Сбросить блёклый.
	seqResetItalic    = "23"    // Сбросить курсив.
	seqResetUnderline = "24"    // Сбросить подчёркнутый один раз.
	seqResetReverse   = "27"    // Сбросить инвертирование цвета.
	seqResetCrossOut  = "29"    // Сбросить зачёркнутый.
)

const (
	// Время ожидания для OSC запросов.
	timeoutOSC       = 5 * time.Second
	key24Bit         = "24bit"
	keyTrueColor     = "truecolor"
	keyScreen        = "screen"
	keyTmux          = "tmux"
	keyYes           = "yes"
	keyTrue          = "true"
	keyXtermKitty    = "xterm-kitty"
	keyLinux         = "linux"
	key256Color      = "256color"
	keyColor         = "color"
	keyAnsi          = "ansi"
	keyOn            = "ON"
	keyHash          = "#"
	prefixForeground = "38"
	prefixBackground = "48"
	prefixRgb        = ";rgb:"
)

var (
	errStatusReport   = errors.New("не удалось получить отчёт статуса терминала")
	errReturnedNoData = errors.New("возвращён ответ - нет данных")
	errTimeout        = errors.New("время истекло")
	errInvalidColor   = errors.New("не верный цвет")
)

/*

    // Незадействованные последовательности.

	//seqOSC            = "\x1b]" // Последовательность OSC "команда операционной системы".
	//seqBlink          = "5"     // Мерцание медленно.
	//seqBlinkFast      = "6"     // Мерцание часто.
	//seqInvisible      = "8"     // Невидимый.
	//seqMainFont       = "10"    // Основной шрифт (по умолчанию).
	//seqResetBold      = "21"    // Сбросить жирный или двойное подчёркивание.
	//seqResetBlink     = "25"    // Сбросить мерцание.
	//seqResetInvisible = "28"    // Сбросить невидимый.
	//seqFramed         = "51"    // Обрамлённый.
	//seqSurrounded     = "52"    // Окружённый.
	//seqOverline       = "53"    // Надчёркнутый.
	//seqResetFramed    = "54"    // Сбросить обрамлённый и окружённый.
	//seqResetOverline  = "55"    // Сбросить надчёркнутый.

	// 11-19   - Альтернативные шрифты.
	// 20      - Не поддерживается нигде.
	// 30-37   - Цвет текста обычной яркости.
	// 38      - Зарезервировано для дополнительных цветов.
	// 39      - Цвет текста по умолчанию.
	// 40-47   - Цвет фона обычной яркости.
	// 48      - Зарезервировано для установки расширенного цвета фона.
	// 39      - Цвет фона по умолчанию.
	// 50      - Зарезервировано.
	// 56-59   - Зарезервировано.
	// 90–97   - Цвет текста повышенной яркости.
	// 100–107 - Цвет фона повышенной яркости.
	// OSC     - Последовательность используется для установки заголовка окна, изменения цветов экрана.

*/
