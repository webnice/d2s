package types

import "io"

// SyncWriter Интерфейс io.Writer с функцией принудительного сброса буфера.
type SyncWriter interface {
	// Writer Наследование стандартного интерфейса io.Writer.
	io.Writer

	// Sync Функция принудительного сброса буфера.
	Sync() error
}
