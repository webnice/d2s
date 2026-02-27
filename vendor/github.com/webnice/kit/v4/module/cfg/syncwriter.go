package cfg

import (
	"os"

	kitTypes "github.com/webnice/kit/v4/types"
)

// Конструктор объекта syncWriter реализующий интерфейс kitTypes.SyncWriter.
func newSyncWriter(parent *impl) kitTypes.SyncWriter {
	var ret = &syncWriter{
		parent: parent,
		wrh:    os.Stderr,
	}

	return ret
}

// Интерфейс io.Writer.
func (sw *syncWriter) Write(buf []byte) (n int, err error) {
	if !sw.parent.isForkWorker {
		n, err = sw.wrh.Write(buf)
		return
	}

	// TODO: Сделать вывод данных для режима fork worker.

	return
}

// Sync Интерфейс kitTypes.SyncWriter.
func (sw *syncWriter) Sync() (err error) {
	if !sw.parent.isForkWorker {
		err = sw.wrh.Sync()
		return
	}

	// TODO: Сделать вывод данных для режима fork worker.

	return
}
