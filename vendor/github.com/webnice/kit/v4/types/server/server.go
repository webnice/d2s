package server

// Server Структура с информацией о сервере.
type Server struct {
	// T Тип сервера, web, grpc, tcp, udp.
	T TServer

	// Web Конфигурация ВЕБ сервера.
	Web *WebConfiguration

	// Grpc Конфигурация GRPC сервера.
	Grpc *GrpcConfiguration

	// Tcp Конфигурация TCP сервера.
	//Tcp *TcpConfiguration

	// Udp Конфигурация UDP сервера.
	//Udp *UdpConfiguration
}
