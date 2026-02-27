package types

// MainFn Функция-точка запуска и завершения приложения.
type MainFn func() (code uint8, err error)

// ComponentPreferencesFn Тип функции загрузки и обработки настроек компонента.
type ComponentPreferencesFn func(*ComponentInfo) (*ComponentInfo, error)

// ComponentConflictFn Тип функции проверки конфликтов между компонентами.
type ComponentConflictFn func([]*ComponentInfo) error

// ComponentRequiresFn Тип функции проверки зависимостей между компонентами.
type ComponentRequiresFn func([]*ComponentInfo) error

// ComponentSortFn Тип функции сортировки компонентов в соответствии с настройками (before) и (after).
type ComponentSortFn func([]*ComponentInfo) error

// ComponentInitiateFn Тип функции вызова функции Initiate компонента с контролем длительности выполнения.
type ComponentInitiateFn func(*ComponentInfo) error

// ComponentDoFn Тип функции вызова функции Do компонента.
type ComponentDoFn func(*ComponentInfo) error

// ComponentFinalizeFn Тип функции вызова функции Finalize компонента.
type ComponentFinalizeFn func(*ComponentInfo) error
