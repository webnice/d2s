package log

import (
	"sync/atomic"
	"time"

	kitModuleLogLevel "github.com/webnice/kit/v4/module/log/level"
	kitModuleTrace "github.com/webnice/kit/v4/module/trace"
)

// Write Реализация интерфейса io.Writer.
func (log *logger) Write(b []byte) (n int, err error) {
	const stackBack = 3
	var (
		msg *Message
		buf []byte
	)

	atomic.AddInt64(&log.waitCounter, 1)
	defer func() { atomic.AddInt64(&log.waitCounter, -1) }()
	buf = make([]byte, len(b))
	if n = copy(buf, b); n == 0 {
		return
	}
	msg = log.MessageGet()
	msg.Timestamp = time.Now().In(time.UTC)
	msg.Level = findLoglevel()
	msg.Fatality = msg.Level == kitModuleLogLevel.Fatal
	kitModuleTrace.Short(msg.Trace, stackBack)
	n, err = msg.Pattern.
		Write(rexSpaceLast.ReplaceAll(rexSpaceFirst.ReplaceAll(buf, []byte{}), []byte{}))
	log.Message(msg)

	return
}

// WriteString Реализация интерфейса io.StringWriter.
func (log *logger) WriteString(s string) (n int, err error) {
	const stackBack = 3
	var (
		msg *Message
		buf []byte
	)

	atomic.AddInt64(&log.waitCounter, 1)
	defer func() { atomic.AddInt64(&log.waitCounter, -1) }()
	buf = []byte(s)
	msg = log.MessageGet()
	msg.Timestamp = time.Now().In(time.UTC)
	msg.Level = findLoglevel()
	msg.Fatality = msg.Level == kitModuleLogLevel.Fatal
	kitModuleTrace.Short(msg.Trace, stackBack)
	n, err = msg.Pattern.
		WriteString(rexSpaceLast.ReplaceAllString(rexSpaceFirst.ReplaceAllString(string(buf), ``), ``))
	log.Message(msg)

	return
}
