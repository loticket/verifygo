package verifygo

import (
	"errors"
)

type PlayStrategy struct {
	lottery LotteryInterface
}

func NewPlayStrategy() *PlayStrategy {
	return &PlayStrategy{}
}

func (play *PlayStrategy) Lottery(plays LotteryInterface) *PlayStrategy {
	play.lottery = plays
	return play
}

func (play *PlayStrategy) SetTicket(playtype int, lottype int, lotnum string, money int, betnum int, multiple int) {
	var playTicket NumLottery = NumLottery{
		PlayType: playtype,
		LotType:  lottype,
		LotNum:   lotnum,
		Money:    money,
		BetNum:   betnum,
		Multiple: multiple,
	}

	play.CreateTicket(playTicket)
}

func (play *PlayStrategy) CreateTicket(ticket NumLottery) {
	play.lottery.CreateTicket(ticket)
}

//验证投注的号码是否正确
func (play *PlayStrategy) Verification() (bool, error) {
	var flags bool = play.lottery.CheckPlaytype()
	if !flags {
		return false, errors.New("子玩法或者投注方式错误")
	}

	flags = play.lottery.PlayCheck()

	if !flags {
		return false, errors.New("投注号码错误")
	}

	var ticketNum int = play.lottery.GetTicketNum()
	var multiple int = play.lottery.GetMultiple()
	var money int = play.lottery.GetMoney()

	if ticketNum != play.lottery.GetBetNum() {
		return false, errors.New("注数计算错误")
	}

	//单注不能超过两万
	if money/multiple > 20000 {
		return false, errors.New("单注金额不能超过2万")
	}

	if ticketNum*2*multiple != money {
		return false, errors.New("金额计算错误")
	}

	return true, nil
}

//按照倍数或者钱数拆票
func (play *PlayStrategy) GetSpliteTicket() []NumLottery {
	var multi []int = play.spliteMultiple()
	var ticketAll []NumLottery = play.lottery.GetSpliteTicket()
	var oneMoney int = 2
	var newTicket []NumLottery = make([]NumLottery, 0)
	for i := 0; i < len(ticketAll); i++ {
		var tempTikcket NumLottery = ticketAll[i]
		for _, val := range multi {
			tempTikcket.Multiple = val
			tempTikcket.Money = val * tempTikcket.BetNum * oneMoney
			newTicket = append(newTicket, tempTikcket)
		}
	}

	return newTicket
}

func (play *PlayStrategy) spliteMultiple() []int {
	var ticketMoney int = play.lottery.GetMoney()
	var ticketMultiple int = play.lottery.GetMultiple()
	if ticketMoney <= MAXMONEY && ticketMultiple <= MAXMULTIPLE {
		return []int{1}
	}

	//如果这张票超过两万 就按照倍数拆票
	//按照倍数钞票
	var singleMoney int = ticketMoney / ticketMultiple //单注票的价格

	//求出最大倍数
	var maxMultiple int = MAXMONEY / singleMoney
	var singleMultiple int = MAXMULTIPLE
	if maxMultiple < MAXMULTIPLE {
		singleMultiple = maxMultiple
	}

	var ticketMutiple []int

	var allMultiple int = ticketMultiple

	for {
		if allMultiple == 0 || allMultiple < 0 {
			break
		}
		var tempMul int = 0
		if allMultiple > singleMultiple {
			tempMul = singleMultiple
		} else {
			tempMul = allMultiple
		}

		ticketMutiple = append(ticketMutiple, tempMul)

		allMultiple = allMultiple - singleMultiple

	}

	return ticketMutiple

}
