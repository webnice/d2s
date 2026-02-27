package cfg

import (
	"os"
	"os/user"
	"path"
	"path/filepath"
	"reflect"
	"runtime"
	"strconv"
	"strings"
)

// Загрузка флага debug до инициализации CLI и конфигурации.
func initApplicationDebug() (ret bool) {
	const (
		keyDebugShort = "-d"
		keyDebugLong  = "--debug"
	)
	var (
		err      error
		value    string
		ok       bool
		arg      string
		argDebug string
	)

	if value, ok = os.LookupEnv(envNameAnchorForApplicationDebug); ok {
		if ret, err = strconv.ParseBool(value); err == nil {
			return
		}
	}
	for _, arg = range os.Args[1:] {
		for _, argDebug = range []string{keyDebugShort, keyDebugLong} {
			if strings.EqualFold(arg, argDebug) {
				ret = true
			}
		}
	}

	return
}

// AbsolutePath Функция преобразует относительный путь в абсолютный путь к директории или файлу.
// Учитывается символ '~' обозначающий домашнюю директорию текущего пользователя.
// Так же раскрываются все символические ссылки в пути оригинальных файлов и директорий.
func AbsolutePath(pth string) (ret string) {
	const wDrv = ":"
	var (
		err error
		sep string
		tmp string
	)

	sep = string(os.PathSeparator)
	if ret = strings.TrimSpace(pth); len(ret) == 0 {
		return
	}
	if ret[0] == runeTilda {
		ret = path.Join(Get().User().HomeDir, ret[1:])
	}
	switch runtime.GOOS {
	case keyOsWindows:
		// Если не "*:\\"
		if len(ret) >= 3 && ret[1:3] != wDrv+sep || len(ret) < 3 {
			if tmp, err = os.Getwd(); err == nil {
				ret = path.Join(tmp, ret)
			}
		}
	default:
		if ret[0] != runeSlash {
			if tmp, err = os.Getwd(); err == nil {
				ret = path.Join(tmp, ret)
			}
		}
	}
	if tmp, err = filepath.EvalSymlinks(ret); err == nil {
		ret = tmp
	}
	ret = path.Clean(ret)

	return
}

// Получение полного названия функции.
func getFuncFullName(obj any) (ret string) {
	var (
		rv   reflect.Value
		rt   reflect.Type
		star string
	)

	if rv = indirect(reflect.ValueOf(obj)); !rv.CanAddr() {
		ret = runtime.FuncForPC(rv.Pointer()).Name()
	} else {
		rt = indirectType(reflect.TypeOf(obj))
		if rt.Name() == "" {
			if rt.Kind() == reflect.Pointer {
				star = "*"
			}
		}
		if rt.Name() != "" {
			if rt.PkgPath() == "" {
				ret = star + rt.Name()
			} else {
				ret = star + rt.PkgPath() + "." + rt.Name()
			}
		}
	}

	return
}

// Получение текущего пользователя операционной системы.
func currentUser(cfg *impl) (ret *user.User) {
	var err error

	if cfg.user != nil {
		return cfg.user
	}
	if ret, err = user.Current(); err != nil {
		ret = &user.User{HomeDir: func() (homeDir string) {
			if homeDir, err = os.UserHomeDir(); err != nil {
				homeDir, err = keySlash, nil
			}
			return
		}()}
		cfg.Gist().ErrorAppend(Errors().GetCurrentUser.Bind(err))
		return
	}

	return
}
