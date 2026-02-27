package cfg

import (
	"time"

	kitModuleLog "github.com/webnice/kit/v4/module/log"
	kitModuleLogLevel "github.com/webnice/kit/v4/module/log/level"
	kitModuleTrace "github.com/webnice/kit/v4/module/trace"
	kitTypes "github.com/webnice/kit/v4/types"
)

// Получение объекта интерфейса kitTypes.Logger.
func newRec(rer *recorder) (rec *record) {
	const stackBack = 2
	rec = &record{
		recorder:  rer,
		traceInfo: kitTypes.NewTraceInfo(),
		stackBack: stackBack,
		keys:      make(map[string]any),
		fatality:  nil,
	}

	return
}

// Key Ключи логирования дополняющие лог.
func (rec *record) Key(keys ...kitTypes.LoggerKey) kitTypes.Logger {
	var (
		n   int
		key string
	)

	// Копирование ключей
	for n = range keys {
		for key = range keys[n] {
			rec.keys[key] = keys[n][key]
		}
	}

	return rec
}

// Time Переопределение времени записи переданным значением.
func (rec *record) Time(timestamp time.Time) kitTypes.Logger {
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
func (rec *record) StackBackCorrect(stepBack int) kitTypes.Logger {
	rec.stackBackCorrect = stepBack

	return rec
}

// Fatality Изменение режима фатальности по умолчанию для записи лога.
// Допустимые значения:
//   - true  - после вывода записи в лог, приложение получает сигнал немедленного завершения.
//   - false - отменяет завершение приложения, указанное для уровня логирования по умолчанию, например,
//     для записи лога Fatal, можно отменить завершение приложения: log.Fatality(false).Fatal(...).
func (rec *record) Fatality(fy bool) kitTypes.Logger { rec.fatality = &fy; return rec }

// Fatal Уровень 0: уровень предсмертных сообщений - система не стабильна, продолжение работы невозможно.
func (rec *record) Fatal(args ...any) {
	var defaultFatality = true

	kitModuleTrace.Short(rec.traceInfo, rec.stackBack+rec.stackBackCorrect)
	// Присвоение времени записи.
	if rec.timestamp.IsZero() {
		rec.timestamp = time.Now().In(time.UTC)
	}
	// Режим фатальности.
	if rec.fatality == nil {
		rec.fatality = &defaultFatality
	}
	rec.messageSend(kitModuleLogLevel.Fatal, "", args...)
}

// Fatalf Уровень 0: уровень предсмертных сообщений - система не стабильна, продолжение работы невозможно.
func (rec *record) Fatalf(pattern string, args ...any) {
	var defaultFatality = true

	kitModuleTrace.Short(rec.traceInfo, rec.stackBack+rec.stackBackCorrect)
	// Присвоение времени записи.
	if rec.timestamp.IsZero() {
		rec.timestamp = time.Now().In(time.UTC)
	}
	// Режим фатальности.
	if rec.fatality == nil {
		rec.fatality = &defaultFatality
	}
	rec.messageSend(kitModuleLogLevel.Fatal, pattern, args...)
}

// Alert Уровень 1: уровень сообщений тревоги - система нестабильна, но может частично продолжить работу.
func (rec *record) Alert(args ...any) {
	kitModuleTrace.Short(rec.traceInfo, rec.stackBack+rec.stackBackCorrect)
	// Присвоение времени записи.
	if rec.timestamp.IsZero() {
		rec.timestamp = time.Now().In(time.UTC)
	}
	rec.messageSend(kitModuleLogLevel.Alert, "", args...)
}

// Alertf Уровень 1: уровень сообщений тревоги - система нестабильна, но может частично продолжить работу.
func (rec *record) Alertf(pattern string, args ...any) {
	kitModuleTrace.Short(rec.traceInfo, rec.stackBack+rec.stackBackCorrect)
	// Присвоение времени записи.
	if rec.timestamp.IsZero() {
		rec.timestamp = time.Now().In(time.UTC)
	}
	rec.messageSend(kitModuleLogLevel.Alert, pattern, args...)
}

// Critical Уровень 2: уровень критических ошибок - часть функционала системы работает не корректно.
func (rec *record) Critical(args ...any) {
	kitModuleTrace.Short(rec.traceInfo, rec.stackBack+rec.stackBackCorrect)
	// Присвоение времени записи.
	if rec.timestamp.IsZero() {
		rec.timestamp = time.Now().In(time.UTC)
	}
	rec.messageSend(kitModuleLogLevel.Critical, "", args...)
}

// Criticalf Уровень 2: уровень критических ошибок - часть функционала системы работает не корректно.
func (rec *record) Criticalf(pattern string, args ...any) {
	kitModuleTrace.Short(rec.traceInfo, rec.stackBack+rec.stackBackCorrect)
	// Присвоение времени записи.
	if rec.timestamp.IsZero() {
		rec.timestamp = time.Now().In(time.UTC)
	}
	rec.messageSend(kitModuleLogLevel.Critical, pattern, args...)
}

// Error Уровень 3: уровень не критических ошибок - ошибки не прерывающие работу приложения.
func (rec *record) Error(args ...any) {
	kitModuleTrace.Short(rec.traceInfo, rec.stackBack+rec.stackBackCorrect)
	// Присвоение времени записи.
	if rec.timestamp.IsZero() {
		rec.timestamp = time.Now().In(time.UTC)
	}
	rec.messageSend(kitModuleLogLevel.Error, "", args...)
}

// Errorf Уровень 3: уровень не критических ошибок - ошибки не прерывающие работу приложения.
func (rec *record) Errorf(pattern string, args ...any) {
	kitModuleTrace.Short(rec.traceInfo, rec.stackBack+rec.stackBackCorrect)
	// Присвоение времени записи.
	if rec.timestamp.IsZero() {
		rec.timestamp = time.Now().In(time.UTC)
	}
	rec.messageSend(kitModuleLogLevel.Error, pattern, args...)
}

// Warning Уровень 4: уровень сообщений с предупреждениями.
func (rec *record) Warning(args ...any) {
	kitModuleTrace.Short(rec.traceInfo, rec.stackBack+rec.stackBackCorrect)
	// Присвоение времени записи.
	if rec.timestamp.IsZero() {
		rec.timestamp = time.Now().In(time.UTC)
	}
	rec.messageSend(kitModuleLogLevel.Warning, "", args...)
}

// Warningf Уровень 4: уровень сообщений с предупреждениями.
func (rec *record) Warningf(pattern string, args ...any) {
	kitModuleTrace.Short(rec.traceInfo, rec.stackBack+rec.stackBackCorrect)
	// Присвоение времени записи.
	if rec.timestamp.IsZero() {
		rec.timestamp = time.Now().In(time.UTC)
	}
	rec.messageSend(kitModuleLogLevel.Warning, pattern, args...)
}

// Notice Уровень 5: уровень штатных информационных сообщений, требующих повышенного внимания.
func (rec *record) Notice(args ...any) {
	kitModuleTrace.Short(rec.traceInfo, rec.stackBack+rec.stackBackCorrect)
	// Присвоение времени записи.
	if rec.timestamp.IsZero() {
		rec.timestamp = time.Now().In(time.UTC)
	}
	rec.messageSend(kitModuleLogLevel.Notice, "", args...)
}

// Noticef Уровень 5: уровень штатных информационных сообщений, требующих повышенного внимания.
func (rec *record) Noticef(pattern string, args ...any) {
	kitModuleTrace.Short(rec.traceInfo, rec.stackBack+rec.stackBackCorrect)
	// Присвоение времени записи.
	if rec.timestamp.IsZero() {
		rec.timestamp = time.Now().In(time.UTC)
	}
	rec.messageSend(kitModuleLogLevel.Notice, pattern, args...)
}

// Info Уровень 6: сообщения информационного характера описывающие шаги выполнения алгоритмов приложения.
func (rec *record) Info(args ...any) {
	kitModuleTrace.Short(rec.traceInfo, rec.stackBack+rec.stackBackCorrect)
	// Присвоение времени записи.
	if rec.timestamp.IsZero() {
		rec.timestamp = time.Now().In(time.UTC)
	}
	rec.messageSend(kitModuleLogLevel.Info, "", args...)
}

// Infof Уровень 6: сообщения информационного характера описывающие шаги выполнения алгоритмов приложения.
func (rec *record) Infof(pattern string, args ...any) {
	kitModuleTrace.Short(rec.traceInfo, rec.stackBack+rec.stackBackCorrect)
	// Присвоение времени записи.
	if rec.timestamp.IsZero() {
		rec.timestamp = time.Now().In(time.UTC)
	}
	rec.messageSend(kitModuleLogLevel.Info, pattern, args...)
}

// Debug Уровень 7: уровень отладочных сообщений.
func (rec *record) Debug(args ...any) {
	kitModuleTrace.Short(rec.traceInfo, rec.stackBack+rec.stackBackCorrect)
	// Присвоение времени записи.
	if rec.timestamp.IsZero() {
		rec.timestamp = time.Now().In(time.UTC)
	}
	rec.messageSend(kitModuleLogLevel.Debug, "", args...)
}

// Debugf Уровень 7: уровень отладочных сообщений.
func (rec *record) Debugf(pattern string, args ...any) {
	kitModuleTrace.Short(rec.traceInfo, rec.stackBack+rec.stackBackCorrect)
	// Присвоение времени записи.
	if rec.timestamp.IsZero() {
		rec.timestamp = time.Now().In(time.UTC)
	}
	rec.messageSend(kitModuleLogLevel.Debug, pattern, args...)
}

// Trace Уровень 8: уровень максимально подробной трассировки.
func (rec *record) Trace(args ...any) {
	kitModuleTrace.Short(rec.traceInfo, rec.stackBack+rec.stackBackCorrect)
	// Присвоение времени записи.
	if rec.timestamp.IsZero() {
		rec.timestamp = time.Now().In(time.UTC)
	}
	rec.messageSend(kitModuleLogLevel.Trace, "", args...)
}

// Tracef Уровень 8: уровень максимально подробной трассировки.
func (rec *record) Tracef(pattern string, args ...any) {
	kitModuleTrace.Short(rec.traceInfo, rec.stackBack+rec.stackBackCorrect)
	// Присвоение времени записи.
	if rec.timestamp.IsZero() {
		rec.timestamp = time.Now().In(time.UTC)
	}
	rec.messageSend(kitModuleLogLevel.Trace, pattern, args...)
}

// MessageWithLevel Отправка сообщения в лог с указанием уровня логирования.
func (rec *record) MessageWithLevel(level kitModuleLogLevel.Level, pattern string, args ...any) {
	kitModuleTrace.Short(rec.traceInfo, rec.stackBack+rec.stackBackCorrect)
	// Присвоение времени записи.
	if rec.timestamp.IsZero() {
		rec.timestamp = time.Now().In(time.UTC)
	}
	rec.messageSend(level, pattern, args...)
}

// Единое место отправки сообщения с полученным стеком вызовов и уровнем логирования.
func (rec *record) messageSend(lv kitModuleLogLevel.Level, pattern string, args ...any) {
	var (
		isLog bool
		msg   *kitModuleLog.Message
		key   string
	)

	// Логирование с фатальностью всегда передаётся в лог, без учёта уровня логирования.
	if rec.fatality != nil && *rec.fatality {
		isLog = true
	}
	// При включённом режиме отладки логи всегда пишутся, без учёта уровня логирования.
	if !isLog && rec.recorder.parent.bootstrapConfiguration.ApplicationDebug {
		isLog = true
	}
	// Выход, если логирование отключено.
	if !isLog && rec.recorder.parent.Loglevel() <= kitModuleLogLevel.Off {
		return
	}
	// Получение объекта для записи в лог из бассейна.
	msg = rec.recorder.parent.logger.MessageGet()
	// Копирование всех данных лога в полученный объект.
	msg.Timestamp, msg.Level = rec.timestamp, lv
	msg.Pattern.WriteString(pattern)
	msg.Argument = make([]any, 0, len(args))
	msg.Argument = append(msg.Argument, args...)
	for key = range rec.keys {
		msg.Keys[key] = rec.keys[key]
	}
	msg.Trace.Copy(rec.traceInfo)
	if rec.fatality != nil {
		msg.Fatality = *rec.fatality
	}
	// Отправка объекта в лог.
	rec.recorder.parent.logger.Message(msg)
}
