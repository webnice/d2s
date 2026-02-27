package cfg

import (
	"time"

	kitModuleLogLevel "github.com/webnice/kit/v4/module/log/level"
	kitTypes "github.com/webnice/kit/v4/types"
)

// Получение объекта интерфейса Recorder.
func newRecorder(parent *impl) *recorder {
	var rer = &recorder{parent: parent}

	return rer
}

// Key Ключи логирования дополняющие лог.
func (rer *recorder) Key(keys ...kitTypes.LoggerKey) kitTypes.Logger {
	var (
		rec *record
		n   int
		key string
	)

	// Создание объекта.
	rec = newRec(rer)
	rec.stackBack--
	// Копирование ключей.
	for n = range keys {
		for key = range keys[n] {
			rec.keys[key] = keys[n][key]
		}
	}

	return rec
}

// Time Переопределение времени записи переданным значением.
func (rer *recorder) Time(timestamp time.Time) kitTypes.Logger {
	var rec *record

	// Создание объекта.
	rec = newRec(rer)
	rec.stackBack--
	if !timestamp.IsZero() {
		// Изменение зоны времени с локальной на UTC.
		if timestamp.Location() != time.UTC {
			timestamp = timestamp.In(time.UTC)
		}
		// Присвоение значения времени
		rec.timestamp = timestamp
	}

	return rec
}

// StackBackCorrect Дополнительная коррекция просмотра стека вызовов для поиска функции вызвавшей логирование.
// Используется когда целевая функция в стеке вызова смещена больше или меньше стандартного значения.
func (rer *recorder) StackBackCorrect(stepBack int) kitTypes.Logger {
	var rec *record

	// Создание объекта.
	rec = newRec(rer)
	rec.stackBack--
	rec.stackBackCorrect = stepBack

	return rec
}

// Fatality Изменение режима фатальности по умолчанию для записи лога.
// Допустимые значения:
//   - true  - после вывода записи в лог, приложение получает сигнал немедленного завершения.
//   - false - отменяет завершение приложения, указанное для уровня логирования по умолчанию, например,
//     для записи лога Fatal, можно отменить завершение приложения: log.Fatality(false).Fatal(...).
func (rer *recorder) Fatality(fy bool) kitTypes.Logger {
	var rec *record

	// Создание объекта.
	rec = newRec(rer)
	rec.stackBack--
	rec.fatality = &fy

	return rec
}

// Fatal Уровень 0: уровень предсмертных сообщений - система не стабильна, продолжение работы невозможно.
func (rer *recorder) Fatal(ag ...any) { newRec(rer).Fatal(ag...) }

// Fatalf Уровень 0: уровень предсмертных сообщений - система не стабильна, продолжение работы невозможно.
func (rer *recorder) Fatalf(pt string, ag ...any) { newRec(rer).Fatalf(pt, ag...) }

// Alert Уровень 1: уровень сообщений тревоги - система нестабильна, но может частично продолжить работу.
func (rer *recorder) Alert(ag ...any) { newRec(rer).Alert(ag...) }

// Alertf Уровень 1: уровень сообщений тревоги - система нестабильна, но может частично продолжить работу.
func (rer *recorder) Alertf(pt string, ag ...any) { newRec(rer).Alertf(pt, ag...) }

// Critical Уровень 2: уровень критических ошибок - часть функционала системы работает не корректно.
func (rer *recorder) Critical(ag ...any) { newRec(rer).Critical(ag...) }

// Criticalf Уровень 2: уровень критических ошибок - часть функционала системы работает не корректно.
func (rer *recorder) Criticalf(pt string, ag ...any) { newRec(rer).Criticalf(pt, ag...) }

// Error Уровень 3: уровень не критических ошибок - ошибки не прерывающие работу приложения.
func (rer *recorder) Error(ag ...any) { newRec(rer).Error(ag...) }

// Errorf Уровень 3: уровень не критических ошибок - ошибки не прерывающие работу приложения.
func (rer *recorder) Errorf(pt string, ag ...any) { newRec(rer).Errorf(pt, ag...) }

// Warning Уровень 4: уровень сообщений с предупреждениями.
func (rer *recorder) Warning(ag ...any) { newRec(rer).Warning(ag...) }

// Warningf Уровень 4: уровень сообщений с предупреждениями.
func (rer *recorder) Warningf(pt string, ag ...any) { newRec(rer).Warningf(pt, ag...) }

// Notice Уровень 5: уровень штатных информационных сообщений, требующих повышенного внимания.
func (rer *recorder) Notice(ag ...any) { newRec(rer).Notice(ag...) }

// Noticef Уровень 5: уровень штатных информационных сообщений, требующих повышенного внимания.
func (rer *recorder) Noticef(pt string, ag ...any) { newRec(rer).Noticef(pt, ag...) }

// Info Уровень 6: сообщения информационного характера описывающие шаги выполнения алгоритмов приложения.
func (rer *recorder) Info(ag ...any) { newRec(rer).Info(ag...) }

// Infof Уровень 6: сообщения информационного характера описывающие шаги выполнения алгоритмов приложения.
func (rer *recorder) Infof(pt string, ag ...any) { newRec(rer).Infof(pt, ag...) }

// Debug Уровень 7: уровень отладочных сообщений.
func (rer *recorder) Debug(ag ...any) { newRec(rer).Debug(ag...) }

// Debugf Уровень 7: уровень отладочных сообщений.
func (rer *recorder) Debugf(pt string, ag ...any) { newRec(rer).Debugf(pt, ag...) }

// Trace Уровень 8: уровень максимально подробной трассировки.
func (rer *recorder) Trace(ag ...any) { newRec(rer).Trace(ag...) }

// Tracef Уровень 8: уровень максимально подробной трассировки.
func (rer *recorder) Tracef(pt string, ag ...any) { newRec(rer).Tracef(pt, ag...) }

// MessageWithLevel Отправка сообщения в лог с указанием уровня логирования.
func (rer *recorder) MessageWithLevel(lv kitModuleLogLevel.Level, pt string, ag ...any) {
	newRec(rer).MessageWithLevel(lv, pt, ag...)
}
