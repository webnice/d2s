// ****************************************************************************************************************** //
//
//    Функции работы с командной строкой вынесены в отдельный пакет и закрыты интерфейсом с целью стабилизировать
// энтропию которую создают авторы сторонних библиотек, бросающие свои творения, назначающие их устаревшими или
// вовсе удаляющие их из общего доступа.
//    В настоящий момент CLI основан на библиотеке kong, но за данным автором уже имеются не приятные и деструктивные
// действия связанные с его творениями...
//
// В дальнейшем планируется заменить kong на более корректную и качественную реализацию CLI либо написать свою версию,
// с возможностями и подходами используемыми в kong, но в качественной реализации и без изъянов kong.
//
// Данна надстройка создана для максимальной изоляции приложения от используемой библиотеки CLI.
// А так же для реализации недостающих в библиотеке, но требуемых, возможностей по настройке и локализации CLI.
//
// Автор kong отклонил и удалил все pull-request запросы. Объяснением было то что они слишком большие и сложные.
// Но скорее всего это было проявление англосаксонской ненависти и русофобии и не способность сказать причину честно.
//
// ****************************************************************************************************************** //

package cli

import (
	"bytes"
	"fmt"
	"os"
	runtimeDebug "runtime/debug"
	"strings"

	kitModuleCfgCliKong "github.com/webnice/kit/v4/module/cfg/cli/kong"
	kitModuleCfgConst "github.com/webnice/kit/v4/module/cfg/const"
	kitTypes "github.com/webnice/kit/v4/types"
)

// New Конструктор объекта пакета.
func New(bootstrap *kitTypes.BootstrapConfiguration) Interface {
	var cli = &impl{
		bootstrap: bootstrap,
		k2gOut:    &bytes.Buffer{},
		k2gErr:    &bytes.Buffer{},
	}
	return cli
}

// Errors Ошибки известного состояния, которые могут вернуть функции пакета.
func (cli *impl) Errors() *Error { return Errors() }

// Constant Установка названия переменных окружения, используемых отдельно от структуры конфигурации.
func (cli *impl) Constant(env ConstantEnvironmentName) Interface { cli.env = env; return cli }

// Bootstrap Первоначальная инициализация минимальной конфигурации приложения.
func (cli *impl) Bootstrap(b *kitTypes.BootstrapDefaultValue) (err error) {
	// Получение названия приложения из переменной окружения или из функции значения по умолчанию.
	cli.applicationName = os.Getenv(cli.env.Destination(kitModuleCfgConst.EnvironmentApplicationName))
	if cli.applicationName == "" {
		cli.applicationName = b.ApplicationName()
	}
	// Сохранение первоначальных начальных значений конфигурации во внутренней переменной.
	cli.bootstrapDefaultValue = b

	return
}

// Эта говно-поделка=kong паникует вместо возврата ошибки, ловим панику.
func (cli *impl) newKong(bs *kitTypes.BootstrapConfiguration) (err error) {
	var (
		option []kitModuleCfgCliKong.Option
		e      any
	)

	defer func() {
		if e = recover(); e != nil {
			err = fmt.Errorf(panicTemplate, e, string(runtimeDebug.Stack()))
		}
	}()
	// Подготовка опций.
	option = cli.KongBuildOption()
	// Инициализация CLI.
	if cli.k2g, err = kitModuleCfgCliKong.New(bs, option...); err != nil {
		return
	}
	// Добавление динамических глобальных флагов приложения.
	if err = cli.KongDynamicFlag(); err != nil {
		return
	}

	return
}

// Init Инициализация командного интерфейса и загрузка переменных окружения.
func (cli *impl) Init() (help *bytes.Buffer, description string, err error) {
	var n int

	// Создание объекта библиотеки CLI с загрузкой базовой минимальной конфигурации в объект cli.bootstrap.
	if err = cli.newKong(cli.bootstrap); err != nil {
		description, err = err.Error(), cli.Errors().UnexpectedError.Bind()
		return
	}
	// Переназначение имён переменных окружения для некоторых параметров bootstrap конфигурации.
	cli.KongUpdateEnvironmentName()
	// Разбор всех аргументов, флагов и команды командной строки и переменных окружения.
	cli.ctx, err = cli.k2g.Parse(os.Args[1:])
	// Разбор полученных ошибок, преобразование в более структурированных ошибки которые можно сравнивать.
	if err != nil {
		if e, ok := err.(*kitModuleCfgCliKong.ParseError); ok {
			switch e.MustErr().Anchor() {
			case kitModuleCfgCliKong.Errors().Expected("").Anchor():
				description = fmt.Sprintf("%s", e.MustErr().Args()...)
				err = cli.Errors().RequiredCommand.Bind()
			case kitModuleCfgCliKong.Errors().ExpectedOneOf("").Anchor():
				description = fmt.Sprintf("%s", e.MustErr().Args()...)
				err = cli.Errors().RequiredCommand.Bind()
			case kitModuleCfgCliKong.Errors().UnexpectedArgument("").Anchor():
				description = fmt.Sprintf("%q", e.MustErr().Args()...)
				err = cli.Errors().UnknownCommand.Bind()
			case kitModuleCfgCliKong.Errors().UnknownFlag("").Anchor():
				description = fmt.Sprintf("%q", e.MustErr().Args()[0])
				err = cli.Errors().UnknownCommand.Bind()
			case kitModuleCfgCliKong.Errors().UnexpectedArgument("").Anchor():
				description = fmt.Sprintf("%q", e.MustErr().Args()[1])
				err = cli.Errors().UnknownArgument.Bind()
			case kitModuleCfgCliKong.Errors().DecodeValueError("", nil).Anchor(),
				kitModuleCfgCliKong.Errors().DecodeValueEnv(nil, "", "").Anchor():
				description = err.Error()
				err = cli.Errors().NotCorrectArgument.Bind()
			case kitModuleCfgCliKong.Errors().HelperErrorDidYouMean(nil, "").Anchor(),
				kitModuleCfgCliKong.Errors().HelperErrorDidYouMeanOneOf(nil, []string{}).Anchor():
				description = fmt.Sprintf(
					perhapsYouWantedTo, e.MustErr().Args()[0], fmt.Sprintf("%s", e.MustErr().Args()[1:]...),
				)
				err = cli.Errors().UnknownArgument.Bind()
			case kitModuleCfgCliKong.Errors().MissingFlags("").Anchor():
				description = fmt.Sprintf("%q", e.MustErr().Args()[0])
				err = cli.Errors().RequiredFlag.Bind()
			default:
				description, err = err.Error(), cli.Errors().UnexpectedError.Bind()
			}
		}
	}
	// Библиотекой CLI получены флаги или аргументы которые предназначены исключительно для CLI.
	// Создаётся ошибка, при которой дальнейшей запуск приложения не будет выполняться.
	if cli.k2gIsExit {
		return cli.helpDisplayedError()
	}
	// Если ошибки нет, получение команды командной строки приложения, либо команды по умолчанию.
	if err == nil {
		// Команда может быть составной, но команда не может содержать пробелы.
		cli.applicationCommand = strings.Split(cli.ctx.Command(), delimiterSpace)
		for n = range cli.applicationCommand {
			cli.applicationCommand[n] = strings.TrimSpace(cli.applicationCommand[n])
		}
	}

	return
}

// Создаётся ошибка, при которой дальнейшей запуск приложения не будет выполняться, для отображения помощи по CLI.
func (cli *impl) helpDisplayedError() (help *bytes.Buffer, description string, err error) {
	const newLine = "\n"
	var buf []byte

	// Буфер вывода сообщений библиотекой CLI, перехватываются из STDOUT и STDERR.
	help = &bytes.Buffer{}
	buf = cli.k2gOut.Bytes()
	cli.k2gOut.Reset()
	if _, _ = help.Write(buf); len(buf) > 0 {
		_, _ = help.WriteString(newLine)
	}
	buf = cli.k2gErr.Bytes()
	cli.k2gErr.Reset()
	if _, _ = help.Write(buf); len(buf) > 0 {
		_, _ = help.WriteString(newLine)
	}
	err = cli.Errors().HelpDisplayed.Bind()

	return
}

// Command Команда приложения, первая часть.
func (cli *impl) Command() (ret string) {
	if len(cli.applicationCommand) > 0 {
		ret = cli.applicationCommand[1]
	}

	return
}

// CommandFull Полная команда приложения.
func (cli *impl) CommandFull() []string { return cli.applicationCommand }
