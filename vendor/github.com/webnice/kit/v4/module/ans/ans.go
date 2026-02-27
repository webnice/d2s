package ans

import kitTypes "github.com/webnice/kit/v4/types"

// New Конструктор объекта сущности пакета, возвращается интерфейс пакета.
func New(logger kitTypes.Logger) Interface {
	var ans = &impl{
		logger: logger,
	}
	ans.essence = newEssence(ans)

	return ans
}

// Gist Интерфейс служебных методов.
func (ans *impl) Gist() Essence { return ans.essence }

func (ans *impl) logWarningf(pattern string, args ...any) {
	if ans.logger == nil {
		return
	}
	ans.logger.Warningf(pattern, args...)
}

func (ans *impl) logErrorf(pattern string, args ...any) {
	if ans.logger == nil {
		return
	}
	ans.logger.Errorf(pattern, args...)
}
