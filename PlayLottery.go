package verifygo

import (
	"fmt"
)

//每一种玩法都要实现这个接口
type LotteryInterface interface {
	CreateTicket(int, int, string, int, int, int)
	Verification() (bool, error)
	GetTicketNum() int
	CheckPlaytype() bool
	PlayCheck() bool
	GetSpliteTicket()
}

type Instance func() LotteryInterface

var adapters = make(map[string]Instance)

func Register(name string, adapter Instance) {
	if adapter == nil {
		panic("cache: Register adapter is nil")
	}
	if _, ok := adapters[name]; ok {
		panic("cache: Register called twice for adapter " + name)
	}
	adapters[name] = adapter
}

func NewLottery(adapterName string) (adapter LotteryInterface, err error) {
	instanceFunc, ok := adapters[adapterName]
	if !ok {
		err = fmt.Errorf("cache: unknown adapter name %q (forgot to import?)", adapterName)
		return
	}
	adapter = instanceFunc()
	if err != nil {
		adapter = nil
	}
	return
}
