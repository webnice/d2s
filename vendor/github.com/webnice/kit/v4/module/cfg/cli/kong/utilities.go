package kong

import (
	"fmt"
	"regexp"
)

func valueOrDefaultValue(value any, defaultValue any) any {
	switch value.(type) {
	case string:
		if value == "" {
			value = defaultValue
		}
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		if value == 0 {
			value = defaultValue
		}
	case bool:
		if value == false {
			value = defaultValue
		}
	}

	return value
}

func oneSpace(template string, args ...any) (ret string) {
	var (
		rex *regexp.Regexp
	)

	rex = regexp.MustCompile(` +`)
	return rex.ReplaceAllString(fmt.Sprintf(template, args...), " ")
}

// Очистка среза цитированных строк от повторяющихся и пустых строк значений.
func uniqueQuotedNotEmpty(src []string) (ret []string) {
	var (
		tmp map[string]bool
		ok  bool
		n   int
	)

	tmp, ret = make(map[string]bool), make([]string, 0, len(src))
	for n = range src {
		if _, ok = tmp[src[n]]; ok || src[n] == "" || src[n] == `""` {
			continue
		}
		tmp[src[n]], ret = true, append(ret, src[n])
	}

	return
}
