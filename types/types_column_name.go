// Package types
package types

import (
	"strings"
	"unicode"
)

// ColumnName Тип названия колонки.
type ColumnName string

// CamelCase Преобразование имени колонки в верблюжью нотацию и совместимый с golang формат.
func (cn *ColumnName) CamelCase() (ret string) {
	const underscore = '_'
	var (
		name  string
		runes []rune
		c     rune
		i     int
		ok    bool
	)

	name = cn.lint()
	runes = []rune(name)
	for i, c = range runes {
		ok = unicode.IsLetter(c) || unicode.IsDigit(c)
		if i == 0 {
			ok = unicode.IsLetter(c)
		}
		if !ok {
			runes[i] = underscore
		}
	}
	ret = string(runes)

	return
}

func (cn *ColumnName) lint() (ret string) {
	const underscore = '_'
	var (
		name, word, u string
		allLower, eow bool
		runes         []rune
		r             rune
		i, w, n       int
		ok            bool
	)

	if name = string(*cn); name == string(underscore) {
		ret = name
		return
	}
	for len(name) > 0 && name[0] == underscore {
		name = name[1:]
	}
	if len(name) > 0 && unicode.IsDigit(rune(name[0])) {
		if _, ok = NumberToWordMap[rune(name[0])]; ok {
			name = NumberToWordMap[rune(name[0])] + string(underscore) + name[1:]
		}
	}
	allLower = true
	for _, r = range name {
		if !unicode.IsLower(r) {
			allLower = false
			break
		}
	}
	if allLower {
		runes = []rune(name)
		if u = strings.ToUpper(name); cn.isAbbreviation(u) {
			copy(runes[0:], []rune(u))
		} else {
			runes[0] = unicode.ToUpper(runes[0])
		}
		ret = string(runes)
		return
	}
	runes = []rune(name)
	for i+1 <= len(runes) {
		eow = false
		switch {
		case i+1 == len(runes):
			eow = true
		case runes[i+1] == underscore:
			eow, n = true, 1
			for i+n+1 < len(runes) && runes[i+n+1] == underscore {
				n++
			}
			if i+n+1 < len(runes) && unicode.IsDigit(runes[i]) && unicode.IsDigit(runes[i+n+1]) {
				n--
			}
			copy(runes[i+1:], runes[i+n+1:])
			runes = runes[:len(runes)-n]
		case unicode.IsLower(runes[i]) && !unicode.IsLower(runes[i+1]):
			eow = true
		}
		if i++; !eow {
			continue
		}
		word = string(runes[w:i])
		if u = strings.ToUpper(word); cn.isAbbreviation(u) {
			copy(runes[w:], []rune(u))

		} else if strings.ToLower(word) == word {
			runes[w] = unicode.ToUpper(runes[w])
		}
		w = i
	}
	ret = string(runes)

	return
}

func (cn *ColumnName) isAbbreviation(s string) (ret bool) {
	for i := range Abbreviations {
		if ret = strings.EqualFold(Abbreviations[i], s); ret {
			return
		}
	}

	return
}
