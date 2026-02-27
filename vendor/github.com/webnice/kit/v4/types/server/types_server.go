package server

import "strings"

var (
	TUnknown = TServer{u: 0, s: ""}
	TUdp     = TServer{u: 1, s: "udp"}
	TTcp     = TServer{u: 2, s: "tcp"}
	TWeb     = TServer{u: 3, s: "web"}
	TGrpc    = TServer{u: 4, s: "grpc"}
)

var _, _ = TServerParseUint16(0), TServerParseString("")

// TServer Тип сервера.
type TServer struct {
	u uint16
	s string
}

// TServerParseUint16 Преобразование uint16 в тип.
func TServerParseUint16(i uint16) (ret TServer) {
	switch i {
	case TWeb.Uint16():
		ret = TWeb
	case TGrpc.Uint16():
		ret = TGrpc
	case TTcp.Uint16():
		ret = TTcp
	case TUdp.Uint16():
		ret = TUdp
	default:
		ret = TUnknown
	}

	return ret
}

// TServerParseString Преобразование строки в тип.
func TServerParseString(s string) (ret TServer) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case TWeb.String():
		ret = TWeb
	case TGrpc.String():
		ret = TGrpc
	case TTcp.String():
		ret = TTcp
	case TUdp.String():
		ret = TUdp
	default:
		ret = TUnknown
	}

	return ret
}

// Uint16 Представление типа в качестве числа.
func (sto *TServer) Uint16() uint16 { return sto.u }

// String Реализация интерфейса Stringify.
func (sto *TServer) String() string { return sto.s }
