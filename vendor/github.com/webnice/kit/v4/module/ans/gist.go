package ans

// Создание объекта и возвращение интерфейса Essence.
func newEssence(parent *impl) *gist {
	var ece = &gist{
		parent: parent,
	}

	return ece
}

// Debug Присвоение нового значения режима отладки.
func (ece *gist) Debug(debug bool) Essence { ece.parent.debug = debug; return ece }
