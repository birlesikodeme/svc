package svc

import (
	"time"

	"github.com/mosteknoloji/glog"

	"golang.org/x/sys/windows/svc"
)

type Server interface {
	Start()
	Stop()
}

type service struct {
	server Server
}

func (w *service) Execute(args []string, req <-chan svc.ChangeRequest, stat chan<- svc.Status) (svcSpecificEC bool, exitCode uint32) {
	stat <- svc.Status{State: svc.StartPending}
	go func() {
		w.server.Start()
	}()

	stat <- svc.Status{State: svc.Running, Accepts: svc.AcceptStop | svc.AcceptShutdown}
loop:
	for {
		switch r := <-req; r.Cmd {
		case svc.Stop, svc.Shutdown:
			w.server.Stop()
			break loop

		case svc.Interrogate:
			stat <- r.CurrentStatus
			time.Sleep(100 * time.Millisecond)
			stat <- r.CurrentStatus
		}
	}

	stat <- svc.Status{State: svc.StopPending}
	return
}

func RunAsService(name string, server Server) {
	if err := svc.Run(name, &service{server: server}); err != nil {
		glog.Error(err.Error())
	}
}
