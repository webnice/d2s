package cli

import "strings"

// ConstantEnvironment Названия специальных переменных окружения.
type ConstantEnvironment struct {
	// Anchor Название якоря переменной.
	Anchor string

	// Destination Требуемое название переменной окружения.
	Destination string
}

// ConstantEnvironmentName Срез названий специальных переменных окружения.
type ConstantEnvironmentName []*ConstantEnvironment

// MustFindByAnchor Поиск по якорю, возвращается либо найденный объект, либо пустой объект.
func (cen ConstantEnvironmentName) MustFindByAnchor(name string) (ret *ConstantEnvironment) {
	for n := range cen {
		if strings.EqualFold(cen[n].Anchor, name) {
			ret = cen[n]
			break
		}
	}
	if ret == nil {
		ret = &ConstantEnvironment{}
	}

	return
}

// Destination Требуемое название переменной окружения.
func (cen ConstantEnvironmentName) Destination(name string) string {
	return cen.MustFindByAnchor(name).Destination
}
