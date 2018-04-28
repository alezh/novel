package system

type (
	Enging interface {

	}
	System struct {

	}
)


var SystemInfo  = SysInterface()

func SysInterface() Enging {
	return initSystem()
}

func initSystem() *System {
	return new(System)
}

func (sys *System)Init()  {
	
}

func (sys *System)ReInit()  {

}

func (sys *System)Start()  {

}

func (sys *System)Stop()  {

}

func (sys *System)GetConfig()  {
	
}

func (sys *System)SetConfig()  {
	
}


