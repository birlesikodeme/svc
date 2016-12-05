package svc

type Server interface {
	Start()
	Stop()
}

func RunAsService(name string, server Server) {
	server.Start()
}
