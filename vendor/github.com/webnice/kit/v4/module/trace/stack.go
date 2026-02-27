package trace

import "runtime"

// StackShort Загрузка короткого стека текущей горутины.
func StackShort() (ret []byte) {
	var n int

	ret = make([]byte, defaultBufferLength)
	for {
		if n = runtime.Stack(ret, false); n < len(ret) {
			ret = ret[:n]
			break
		}
		ret = make([]byte, len(ret)*2)
	}

	return
}

// StackFull Загрузка полного стека вызовов всех горутин приложения.
func StackFull() (ret []byte) {
	var n int

	ret = make([]byte, defaultBufferLength)
	for {
		if n = runtime.Stack(ret, true); n < len(ret) {
			ret = ret[:n]
			break
		}
		ret = make([]byte, len(ret)*2)
	}

	return
}
