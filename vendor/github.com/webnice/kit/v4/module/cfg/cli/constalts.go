package cli

const (
	panicTemplate                    = "Паника: %q\nСтек вызовов, в момент паники:\n%s."
	helpKey                          = "help"
	helpDescription                  = "Отображение помощи по командам и параметрам приложения."
	helpShortKey                     = 'h'
	usageHelperTemplate              = "Используйте: %s %s"
	runArgumentHelperTemplate        = "Выполните: \"%s --%s\" для получения подробной информации."
	runCommandArgumentHelperTemplate = "Выполните: \"%s <command> --%s\" для получения подробной информации о команде."
	helpCommandsLabel                = "Команды:"
	helpFlagsLabel                   = "Флаги:"
	helpArgumentsLabel               = "Аргументы:"
	perhapsYouWantedTo               = "%s, возможно, вы хотели указать: %s."
	delimiterSpace                   = " "
)
