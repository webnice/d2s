package cfg

// RunlevelSubscriberFn Тип функции для подписки на события изменения уровня работы приложения
// Функция вызывается каждый раз при изменении уровня работы приложения
type RunlevelSubscriberFn func(old uint16, new uint16)

// Тип структуры отправляемой в канал изменения уровня работы приложения
type runLevelUp struct {
	newLever uint16
	done     chan struct{}
}
