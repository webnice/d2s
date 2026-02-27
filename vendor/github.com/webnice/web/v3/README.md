# web

[![Go Reference](https://pkg.go.dev/badge/github.com/webnice/web/v3.svg)](https://pkg.go.dev/github.com/webnice/web/v3)
[![Go Report Card](https://goreportcard.com/badge/github.com/webnice/web/v3)](https://goreportcard.com/report/github.com/webnice/web/v3)
[![Coverage Status](https://coveralls.io/repos/github/webnice/web/badge.svg?branch=master)](https://coveralls.io/github/webnice/web?branch=master)

#### Описание

Это микро сборка для создания управляемого веб сервера на основе net/http.

#### Поддержка

    tcp, tcp4, tcp6  - Сервер поднимается на указанном Host:Port;
    unix, unixpacket - Сервер поднимается на указанном unix/unixpacket;
    socket           - Сервер поднимается на юникс-сокет;
    systemd          - Сервер поднимается на сокет, который открывает служба systemd, получае сокет по имени и
                       взаимодействует с сокет службой systemd;

    Поддержка прокси-протокол версии 1 и 2, реализованных в HAProxy и nginx.

#### Зависимости

    github.com/webnice/dic
    github.com/webnice/net
    github.com/pires/go-proxyproto

#### Подключение
```bash
go get github.com/webnice/web/v3
```
