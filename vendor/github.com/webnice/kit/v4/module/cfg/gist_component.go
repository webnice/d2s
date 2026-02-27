package cfg

import (
	"math"
	"reflect"
	"sort"
	"strings"
	"time"

	"github.com/webnice/dic"
	kitModuleCfgCli "github.com/webnice/kit/v4/module/cfg/cli"
	kitTypes "github.com/webnice/kit/v4/types"
)

// ComponentName Получение уникального имени пакета компоненты.
func (essence *gist) ComponentName(obj any) (ret string) {
	var rt reflect.Type

	if rt = reflect.TypeOf(obj); rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}
	ret = rt.PkgPath()

	return
}

// ComponentNames Возвращает список зарегистрированных компонентов.
func (essence *gist) ComponentNames() (ret []string) {
	if !essence.parent.comp.IsClose {
		return
	}
	essence.parent.comp.ComponentMutex.Lock()
	for _, c := range essence.parent.comp.Component {
		ret = append(ret, c.ComponentName)
	}
	essence.parent.comp.ComponentMutex.Unlock()

	return
}

// ComponentPreferences Функция-менеджер загрузки и обработки настроек компонентов.
func (essence *gist) ComponentPreferences(fn kitTypes.ComponentPreferencesFn) (code uint8, err error) {
	var (
		ierr     dic.IError
		src, res *kitTypes.ComponentInfo
	)

	essence.parent.comp.ComponentMutex.Lock()
	defer essence.parent.comp.ComponentMutex.Unlock()
	if !essence.parent.comp.IsClose {
		code = essence.parent.Errors().ComponentPreferencesCallBeforeCompleting.CodeU8().Get()
		essence.ErrorAppend(essence.parent.Errors().ComponentPreferencesCallBeforeCompleting.Bind())
		return
	}
	for _, src = range essence.parent.comp.Registered {
		if res, err = fn(src); err != nil {
			if ierr = Errors().Unbind(err); ierr != nil {
				code = ierr.CodeU8().Get()
			}
			return
		}
		essence.parent.comp.Component = append(essence.parent.comp.Component, res)
	}
	// Очистка значений среза временного списка зарегистрированных компонентов.
	essence.parent.comp.Registered = essence.parent.comp.Registered[:0]

	return
}

// ComponentCheckConflict Проверка конфликтов между всеми зарегистрированными компонентами.
func (essence *gist) ComponentCheckConflict(fn kitTypes.ComponentConflictFn) (code uint8, err error) {
	var ierr dic.IError

	essence.parent.comp.ComponentMutex.Lock()
	defer essence.parent.comp.ComponentMutex.Unlock()
	if err = fn(essence.parent.comp.Component); err != nil {
		if ierr = Errors().Unbind(err); ierr != nil {
			code = ierr.CodeU8().Get()
		}
		return
	}

	return
}

// ComponentRequiresCheck Проверка зависимостей между всеми зарегистрированными компонентами.
func (essence *gist) ComponentRequiresCheck(fn kitTypes.ComponentRequiresFn) (code uint8, err error) {
	var ierr dic.IError

	essence.parent.comp.ComponentMutex.Lock()
	defer essence.parent.comp.ComponentMutex.Unlock()
	if err = fn(essence.parent.comp.Component); err != nil {
		if ierr = Errors().Unbind(err); ierr != nil {
			code = ierr.CodeU8().Get()
		}
		return
	}

	return
}

// ComponentSort Сортировка зарегистрированных компонентов в соответствии с настройками (before) и (after).
func (essence *gist) ComponentSort(fn kitTypes.ComponentSortFn) (code uint8, err error) {
	var ierr dic.IError

	essence.parent.comp.ComponentMutex.Lock()
	defer essence.parent.comp.ComponentMutex.Unlock()
	if err = fn(essence.parent.comp.Component); err != nil {
		if ierr = Errors().Unbind(err); ierr != nil {
			code = ierr.CodeU8().Get()
		}
		return
	}

	return
}

// ComponentMapRunlevel Построение шагов переключения уровня выполнения приложения (runlevel).
func (essence *gist) ComponentMapRunlevel(targetlevel uint16) (code uint8, err error) {
	var (
		n   int
		l   uint16
		tmp map[uint16]bool
	)

	essence.parent.comp.ComponentMutex.Lock()
	defer essence.parent.comp.ComponentMutex.Unlock()
	// Очистка карты уровней
	essence.parent.mapRunLevel = essence.parent.mapRunLevel[:0]
	tmp = make(map[uint16]bool)
	tmp[defaultRunlevel], tmp[targetlevel], tmp[math.MaxUint16] = true, true, true
	// Обход всех зарегистрированных компонентов, получение уровня отличного от нуля
	for n = range essence.parent.comp.Component {
		if essence.parent.comp.Component[n].Runlevel != 0 {
			tmp[essence.parent.comp.Component[n].Runlevel] = true
		}
	}
	// Увеличение памяти для хранения среза карты уровней, если изначально выделено недостаточно.
	if cap(essence.parent.mapRunLevel) < len(tmp) {
		essence.parent.mapRunLevel = make([]uint16, 0, len(tmp))
	}
	// Копирование получившейся карты из map в массив.
	for l = range tmp {
		essence.parent.mapRunLevel = append(essence.parent.mapRunLevel, l)
	}
	// Сортировка по значению уровня выполнения приложения.
	sort.Slice(essence.parent.mapRunLevel, func(i, j int) bool {
		return essence.parent.mapRunLevel[i] < essence.parent.mapRunLevel[j]
	})

	return
}

// ComponentInitiate Вызов функции Initiate у всех зарегистрированных компонентов в прямом порядке.
func (essence *gist) ComponentInitiate(fn kitTypes.ComponentInitiateFn) (code uint8, err error) {
	var (
		ierr dic.IError
		n    int
	)

	essence.parent.comp.ComponentMutex.Lock()
	defer essence.parent.comp.ComponentMutex.Unlock()
	for n = range essence.parent.comp.Component {
		// Пропускаем отключённые компоненты.
		if essence.parent.comp.Component[n].IsDisable {
			continue
		}
		// Пропускаем уже проинициализированные компоненты.
		if essence.parent.comp.Component[n].IsInitiate {
			continue
		}
		if err = fn(essence.parent.comp.Component[n]); err != nil {
			if ierr = Errors().Unbind(err); ierr != nil {
				code = ierr.CodeU8().Get()
			}
			break
		}
	}

	return
}

// ComponentDoFindByCommand Функция поиска компоненты по команде.
func (essence *gist) ComponentDoFindByCommand(n int) (ret bool) {
	var i int

	for i = range essence.parent.comp.Component[n].Command {
		// Компонента подходит если у неё полностью совпадает команда, а так же, если у компоненты есть
		// пустая команда, тогда она тоже запускается как для своей команды, так и для всех других команд.
		if strings.EqualFold(essence.parent.comp.Component[n].Command[i], essence.parent.Command()) ||
			essence.parent.comp.Component[n].Command[i] == "" {
			ret = true
		}
	}
	// Если у компоненты вообще нет команд, то она подходит для всех команд.
	if len(essence.parent.comp.Component[n].Command) == 0 {
		ret = true
	}

	return
}

// ComponentDo Вызов функции Do у всех зарегистрированных компонентов в прямом порядке для указанного уровня приложения.
func (essence *gist) ComponentDo(runlevel uint16, fn kitTypes.ComponentDoFn) (code uint8, err error) {
	var (
		bcfw *kitTypes.BootstrapConfigurationForkWorker
		ierr dic.IError
		n    int
	)

	bcfw = essence.parent.ForkWorker()
	essence.parent.comp.ComponentMutex.Lock()
	defer essence.parent.comp.ComponentMutex.Unlock()
	for n = range essence.parent.comp.Component {
		// Пропускаем отключённые компоненты.
		if essence.parent.comp.Component[n].IsDisable {
			continue
		}
		// Пропускаем не инициализированные компоненты.
		if !essence.parent.comp.Component[n].IsInitiate {
			continue
		}
		// Пропускаем ранее запущенные компоненты.
		if essence.parent.comp.Component[n].IsDo {
			continue
		}
		// Пропускаем компоненты с уровнем запуска не равному текущему уровню работы приложения.
		if essence.parent.comp.Component[n].Runlevel > 0 && essence.parent.comp.Component[n].Runlevel != runlevel {
			continue
		}
		// Выбор на основе флага isForkWorker.
		switch essence.parent.isForkWorker {
		case true:
			switch {
			// Режим forkWorker и указанная компонента, имеет приоритет над командой приложения.
			case bcfw.Component != "":
				// Запускается только компонента указанная для запуска.
				if !strings.EqualFold(bcfw.Component, essence.parent.comp.Component[n].ComponentName) {
					continue
				}
			default:
				// Если приложение запущено с командой, поиск компоненты с указанной командой, пропуск остальных.
				if essence.parent.Command() != "" && !essence.ComponentDoFindByCommand(n) {
					continue
				}
			}
		default:
			// Если приложение запущено с командой, поиск компоненты с указанной командой, пропуск остальных.
			if essence.parent.Command() != "" && !essence.ComponentDoFindByCommand(n) {
				continue
			}
		}
		// Если у компоненты уровень по умолчанию (=0), присваиваем уровень на котором был выполнен запуск компоненты.
		if essence.parent.comp.Component[n].Runlevel == 0 {
			essence.parent.comp.Component[n].Runlevel = runlevel
		}
		// Помечаем то что компонента была запущена.
		essence.parent.comp.Component[n].IsDo = true
		// Выполнение запуска компоненты с использованием внешней функции запуска.
		if err = fn(essence.parent.comp.Component[n]); err != nil {
			if ierr = Errors().Unbind(err); ierr != nil {
				code = ierr.CodeU8().Get()
			}
			break
		}
	}

	return
}

// ComponentFinalizeWarningTimeout Возвращает время отводимое на выполнение функции Finalize(), до печати в лог
// сообщения о долгой работе функции.
func (essence *gist) ComponentFinalizeWarningTimeout() (ret time.Duration) {
	ret = essence.parent.finalizeWarningTimeout
	return
}

// ComponentFinalize Вызов функции Finalize у всех зарегистрированных компонентов в обратном порядке.
func (essence *gist) ComponentFinalize(fn kitTypes.ComponentFinalizeFn) (code uint8, err error) {
	var (
		ierr dic.IError
		n    int
	)

	essence.parent.comp.ComponentMutex.Lock()
	defer essence.parent.comp.ComponentMutex.Unlock()
	for n = len(essence.parent.comp.Component) - 1; n >= 0; n-- {
		// Пропускаем отключённые компоненты.
		if essence.parent.comp.Component[n].IsDisable {
			continue
		}
		// Пропускаем не инициализированные компоненты.
		if !essence.parent.comp.Component[n].IsInitiate {
			continue
		}
		// Пропускаем не запущенные компоненты.
		if !essence.parent.comp.Component[n].IsDo {
			continue
		}
		// Пропускаем ранее "финализированные" компоненты.
		if essence.parent.comp.Component[n].IsFinalize {
			continue
		}
		if err = fn(essence.parent.comp.Component[n]); err != nil {
			if ierr = Errors().Unbind(err); ierr != nil {
				code = ierr.CodeU8().Get()
			}
			break
		}
	}

	return
}

// ComponentCommandRegister Регистрация команды и группы команд компоненты.
func (essence *gist) ComponentCommandRegister(cc kitTypes.ComponentCommand) {
	var (
		groupKey string
		foundPos int
		n        int
	)

	groupKey, foundPos = cc.GroupKey, -1
	for n = range essence.parent.comp.Groups {
		if strings.EqualFold(essence.parent.comp.Groups[n].Key, groupKey) {
			foundPos = n
			break
		}
	}
	if foundPos >= 0 {
		if essence.parent.comp.Groups[foundPos].Title == "" {
			essence.parent.comp.Groups[foundPos].Title = cc.GroupTitle
		}
		if essence.parent.comp.Groups[foundPos].Description == "" {
			essence.parent.comp.Groups[foundPos].Description = cc.GroupDescription
		}
	} else {
		essence.parent.comp.Groups = append(essence.parent.comp.Groups, &kitTypes.CommandGroup{
			Key:         cc.GroupKey,
			Title:       cc.GroupTitle,
			Description: cc.GroupDescription,
		})
	}
	essence.parent.cli.RegisterCommand(&kitModuleCfgCli.Command{
		Command:     cc.Command,
		Description: cc.Description,
		GroupKey:    cc.GroupKey,
		IsDefault:   cc.IsDefault,
		IsHidden:    cc.IsHidden,
		Value:       cc.Value,
	})
}

// ComponentFlagRegister Регистрация глобального флага компоненты приложения.
func (essence *gist) ComponentFlagRegister(cf kitTypes.ComponentFlag) {
	essence.parent.cli.RegisterFlag(&kitModuleCfgCli.Flag{
		ShortKey:    cf.ShortKey,
		Flag:        cf.Flag,
		Description: cf.Description,
		Environment: cf.Environment,
		Placeholder: cf.Placeholder,
		IsRequired:  cf.IsRequired,
		IsHidden:    cf.IsHidden,
		Value:       cf.Value,
	})
}
