package servicehandler

type ServiceHandler struct {
	sv ServiceManager
}

type ServiceManager interface {
	Run()
	RegisterRoute()
	OnShutdown()
}

func NewServiceHandler(svManager ServiceManager) *ServiceHandler {
	return &ServiceHandler{
		sv: svManager,
	}
}

func (sh *ServiceHandler) RunService() {
	sh.sv.RegisterRoute()
	sh.sv.Run()
}
