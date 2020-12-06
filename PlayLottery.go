package verifygo

import (
	"fmt"
)

//每一种玩法都要实现这个接口
type LotteryInterface interface {
	CreateTicket(NumLottery)       //传参数到信息
	GetTicketNum() int             //获取计算的金额
	GetMultiple() int              //获取投注倍数
	GetBetNum() int                //获取投注金额
	GetMoney() int                 //获取当前票的金额
	CheckPlaytype() bool           //检查子玩法是否正确
	PlayCheck() bool               //检查投注信息是否正确
	GetNumberLottery() NumLottery  //获取当前验证票的信息
	GetSpliteTicket() []NumLottery //获取根据规则彩票的票数
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

//获取工厂模式+策略模式
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
