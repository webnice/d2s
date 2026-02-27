package log

import (
	"bytes"
	"time"

	kitModuleLogLevel "github.com/webnice/kit/v4/module/log/level"
	kitTypes "github.com/webnice/kit/v4/types"
)

// Message Объект сообщения лога.
type Message struct {
	Timestamp time.Time               `json:"timestamp"`          // Время записи.
	Level     kitModuleLogLevel.Level `json:"level"`              // Уровень логирования сообщения.
	Pattern   *bytes.Buffer           `json:"pattern,omitempty"`  // Шаблон сообщения.
	Argument  []any                   `json:"argument,omitempty"` // Аргументы сообщения.
	Keys      map[string]any          `json:"keys,omitempty"`     // Ключи сообщения.
	Trace     *kitTypes.TraceInfo     `json:"trace"`              // Информация о вызове и стеке.
	Fatality  bool                    `json:"fatality"`           // Флаг фатальности.
	Done      chan struct{}           `json:"-"`                  // Канал ожидания завершения. После получения сигнала, объект вернётся в бассейн.
}

// Reset Сброс объекта сообщения лога.
func (msg *Message) Reset() {
	var key string

	_ = msg.Timestamp.UnmarshalBinary([]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0xff, 0xff})
	msg.Level = 0
	msg.Pattern.Reset()
	msg.Argument = msg.Argument[:0]
	for key = range msg.Keys {
		delete(msg.Keys, key)
	}
	msg.Trace.Reset()
	msg.Fatality = false
}
