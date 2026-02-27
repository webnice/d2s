package cli

// RegisterCommand Регистрации динамических команд приложения.
func (cli *impl) RegisterCommand(cmd *Command) { cli.command = append(cli.command, cmd) }

// RegisterFlag Регистрация динамических глобальных флагов приложения.
func (cli *impl) RegisterFlag(flg *Flag) { cli.flag = append(cli.flag, flg) }
