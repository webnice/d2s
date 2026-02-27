package cli

import (
	"fmt"
	"reflect"
	"strings"

	kitModuleCfgCliKong "github.com/webnice/kit/v4/module/cfg/cli/kong"
	kitTypes "github.com/webnice/kit/v4/types"
)

// KongBuildOption Подготовка опций.
func (cli *impl) KongBuildOption() (ret []kitModuleCfgCliKong.Option) {
	ret = append(ret, cli.KongExit())
	ret = append(ret, kitModuleCfgCliKong.Name(cli.applicationName))
	ret = append(ret, kitModuleCfgCliKong.Description(cli.bootstrapDefaultValue.ApplicationDescription()))
	ret = append(ret, kitModuleCfgCliKong.ExplicitGroups(cli.KongBuildCommandGroups()))
	ret = append(ret, cli.KongVars()...)
	// Настройка локализованной помощи.
	ret = append(ret, cli.KongHelpHelper()...)
	// Добавление динамических команд в опции.
	ret = append(ret, cli.KongDynamicCommand()...)
	// Перехват вывода в STDOUT, STDERR в буфер реализующий интерфейс io.Writer.
	ret = append(ret, kitModuleCfgCliKong.Writers(cli.k2gOut, cli.k2gErr))
	// Другие опции, частично или полностью не работающие в этой поделке.
	//option = append(option, kong.NoDefaultHelp())
	//option = append(option, kong.UsageOnError())
	//option = append(option, kong.ShortHelp(cli.DefaultShortHelpPrinter))
	//option = append(option, kong.Help(cli.DefaultHelpPrinter))

	return
}

// KongBuildCommandGroups Подготовка групп команд.
func (cli *impl) KongBuildCommandGroups() (ret []kitModuleCfgCliKong.Group) {
	var (
		commandGroup []kitTypes.CommandGroup
		n            int
	)

	if cli.bootstrapDefaultValue.CommandGroup != nil {
		if commandGroup = cli.bootstrapDefaultValue.CommandGroup(); len(commandGroup) > 0 {
			ret = make([]kitModuleCfgCliKong.Group, 0, len(commandGroup))
			for n = range commandGroup {
				if commandGroup[n].Description == "" {
					continue
				}
				ret = append(ret, kitModuleCfgCliKong.Group{
					Key:         commandGroup[n].Key,
					Title:       commandGroup[n].Title,
					Description: commandGroup[n].Description,
				})
			}
		}
	}

	return
}

// KongExit Функция завершения приложения.
func (cli *impl) KongExit() (ret kitModuleCfgCliKong.Option) {
	ret = kitModuleCfgCliKong.Exit(func(i int) {
		cli.k2gIsExit = true
		cli.k2gExitCode = i
		return
	})

	return
}

// KongDynamicCommand Создание динамических команд командной строки в обёртке kong.
// В kong теги работают частично...
func (cli *impl) KongDynamicCommand() (ret []kitModuleCfgCliKong.Option) {
	const (
		defaultCommandTag = "default"
		hiddenCommand     = "hidden"
	)
	var (
		value      any
		valueEmpty struct{}
		tags       []string
		n          int
	)

	// Добавление динамических команд в опции
	for n = range cli.command {
		if value = cli.command[n].Value; value == nil {
			value = &valueEmpty
		}
		tags = make([]string, 0)
		if cli.command[n].IsDefault {
			tags = append(tags, defaultCommandTag)
		}
		if cli.command[n].IsHidden {
			tags = append(tags, hiddenCommand)
		}
		ret = append(ret, kitModuleCfgCliKong.DynamicCommand(
			cli.command[n].Command,
			cli.command[n].Description,
			cli.command[n].GroupKey,
			value,
			tags...,
		))
	}

	return
}

// KongDynamicFlag Создание динамических флагов командной строки и вставка во внутренности kong.
// Автор kong отказался принимать pull-request с дополнительными возможностями.
func (cli *impl) KongDynamicFlag() (err error) {
	const (
		errDuplicateShortKey = "Короткий ключ %q флага %q уже используется у другого флага."
		errDuplicateFlagName = "Флаг с именем %q уже зарегистрирован."
	)
	var (
		flg      *kitModuleCfgCliKong.Flag
		rv       reflect.Value
		n, i     int
		typeName string
	)

	for n = range cli.flag {
		// Проверка совпадения короткого ключа с уже зарегистрированными.
		for i = range cli.k2g.Model.Flags {
			if cli.flag[n].ShortKey != 0 && cli.k2g.Model.Flags[i].Short == cli.flag[n].ShortKey {
				err = fmt.Errorf(errDuplicateShortKey, "-"+string(cli.flag[n].ShortKey), cli.flag[n].Flag)
				return
			}
		}
		// Проверка совпадения имени флага с уже зарегистрированными.
		for i = range cli.k2g.Model.Flags {
			if cli.flag[n].Flag != "" && cli.k2g.Model.Flags[i].Value.Name == cli.flag[n].Flag {
				err = fmt.Errorf(errDuplicateFlagName, cli.flag[n].Flag)
				return
			}
		}
		rv = reflect.ValueOf(cli.flag[n].Value).Elem()
		typeName = rv.Type().Name()
		flg = &kitModuleCfgCliKong.Flag{
			Short:       cli.flag[n].ShortKey,
			Env:         cli.flag[n].Environment,
			PlaceHolder: cli.flag[n].Placeholder,
			Hidden:      cli.flag[n].IsHidden,
			Value: &kitModuleCfgCliKong.Value{
				Name:         cli.flag[n].Flag,
				Help:         cli.flag[n].Description,
				OrigHelp:     cli.flag[n].Description,
				Target:       rv,
				Default:      fmt.Sprintf("%s", cli.flag[n].Value),
				DefaultValue: reflect.ValueOf(cli.flag[n].Value),
				Mapper:       cli.k2g.Registry.ForValue(rv),
				Tag: &kitModuleCfgCliKong.Tag{
					Name:     cli.flag[n].Flag,
					Help:     cli.flag[n].Description,
					TypeName: typeName,
					Env:      cli.flag[n].Environment,
					Short:    cli.flag[n].ShortKey,
					Required: cli.flag[n].IsRequired,
				},
				Required: cli.flag[n].IsRequired,
			},
		}
		cli.k2g.Model.Flags = append(cli.k2g.Model.Flags, cli.k2g.MakeFlag(flg)...)
	}

	return
}

// KongVars Создание переменных для аргументов помощи и команд командной строки.
func (cli *impl) KongVars() (ret []kitModuleCfgCliKong.Option) {
	const (
		varApplicationTargetlevel = "ApplicationTargetlevel"
		varApplicationDebug       = "ApplicationDebug"
		varApplicationName        = "ApplicationName"
		varHomeDirectory          = "HomeDirectory"
		varWorkingDirectory       = "WorkingDirectory"
		varTempDirectory          = "TempDirectory"
		varCacheDirectory         = "CacheDirectory"
		varConfigDirectory        = "ConfigDirectory"
		varLogDirectory           = "LogDirectory"
		varConfigFile             = "ConfigFile"
		varPidFile                = "PidFile"
		varStateFile              = "StateFile"
		varSocketFile             = "SocketFile"
		varLogLevel               = "LogLevel"
	)

	ret = append(ret, kitModuleCfgCliKong.Vars{
		varApplicationTargetlevel: cli.bootstrapDefaultValue.ApplicationTargetlevel(),
		varApplicationDebug:       cli.bootstrapDefaultValue.ApplicationDebug(),
		varApplicationName:        cli.bootstrapDefaultValue.ApplicationName(),
		varHomeDirectory:          cli.bootstrapDefaultValue.HomeDirectory(),
		varWorkingDirectory:       cli.bootstrapDefaultValue.WorkingDirectory(),
		varTempDirectory:          cli.bootstrapDefaultValue.TempDirectory(),
		varCacheDirectory:         cli.bootstrapDefaultValue.CacheDirectory(),
		varConfigDirectory:        cli.bootstrapDefaultValue.ConfigDirectory(),
		varLogDirectory:           cli.bootstrapDefaultValue.LogDirectory(),
		varConfigFile:             cli.bootstrapDefaultValue.ConfigFile(),
		varPidFile:                cli.bootstrapDefaultValue.PidFile(),
		varStateFile:              cli.bootstrapDefaultValue.StateFile(),
		varSocketFile:             cli.bootstrapDefaultValue.SocketFile(),
		varLogLevel:               cli.bootstrapDefaultValue.LogLevel(),
	})

	return
}

// KongHelpHelper Настройка отображения помощи по командам командной строки.
func (cli *impl) KongHelpHelper() (ret []kitModuleCfgCliKong.Option) {
	ret = append(ret, kitModuleCfgCliKong.ShortUsageOnError())
	ret = append(ret, kitModuleCfgCliKong.HelpDisplaySetup(
		helpKey, helpDescription, helpDescription, helpShortKey, false),
	)
	ret = append(ret, kitModuleCfgCliKong.UsageHelperFunc(func(name string, summary string) string {
		return fmt.Sprintf(usageHelperTemplate, name, strings.TrimSpace(summary))
	}))
	ret = append(ret, kitModuleCfgCliKong.RunArgumentHelperFunc(func(name string, summary string) string {
		return fmt.Sprintf(runArgumentHelperTemplate, name, strings.TrimSpace(summary))
	}))
	ret = append(ret, kitModuleCfgCliKong.RunCommandArgumentHelperFunc(func(name string, summary string) string {
		return fmt.Sprintf(runCommandArgumentHelperTemplate, name, strings.TrimSpace(summary))
	}))
	ret = append(ret, kitModuleCfgCliKong.HelpCommandsLabelSetup(helpCommandsLabel))
	ret = append(ret, kitModuleCfgCliKong.HelpFlagsLabelSetup(helpFlagsLabel))
	ret = append(ret, kitModuleCfgCliKong.HelpArgumentsLabelSetup(helpArgumentsLabel))

	return
}

// KongUpdateEnvironmentName Переназначение имён переменных окружения для некоторых параметров bootstrap конфигурации.
func (cli *impl) KongUpdateEnvironmentName() {
	var (
		n   int
		env *ConstantEnvironment
		tmp string
	)

	for n = range cli.k2g.Model.Flags {
		if tmp = strings.TrimSpace(strings.ToUpper(cli.k2g.Model.Flags[n].Env)); tmp == "" {
			continue
		}
		if env = cli.env.MustFindByAnchor(tmp); env.Anchor != tmp {
			continue
		}
		cli.k2g.Model.Flags[n].Env = env.Destination
		if cli.k2g.Model.Flags[n].Tag != nil {
			cli.k2g.Model.Flags[n].Tag.Env = env.Destination
		}
	}
}
