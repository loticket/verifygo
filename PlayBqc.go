package verifygo

import (
	"strings"
)

//足彩半全场
type PlayBqc struct {
	BasePlay
}

//验证子玩法和投注方式
func (play *PlayBqc) CheckPlaytype() bool {
	if play.numLottery.PlayType != 1 {
		return false
	}

	if play.numLottery.LotType != 1 && play.numLottery.LotType != 2 {
		return false
	}

	return true
}

//检查每一注的基本信息
func (play *PlayBqc) PlayCheck() bool {
	var checkRes bool = false
	if play.numLottery.LotType == 1 { //单式
		checkRes = play.checkBallSingle()
	} else if play.numLottery.LotType == 2 { //复试
		checkRes = play.checkBallReexamine()
	}
	return checkRes
}

//单式检查
func (play *PlayBqc) checkBallSingle() bool {
	if (2 * play.numLottery.Multiple) != play.numLottery.Money {
		return false
	}

	if play.numLottery.BetNum != 1 {
		return false
	}

	if strings.Count(play.numLottery.LotNum, ",") != 11 {
		return false
	}
	var lotNumsArr []string = strings.Split(play.numLottery.LotNum, ",")
	if len(lotNumsArr) != 12 {
		return false
	}

	var ballFlag bool = true

	for i := 0; i < 12; i++ {
		if lotNumsArr[i] != "3" && lotNumsArr[i] != "1" && lotNumsArr[i] != "0" {
			ballFlag = false
			break
		}
	}
	if !ballFlag {
		return false
	}

	play.ticketNum = 1

	return true
}

//复试检查
func (play *PlayBqc) checkBallReexamine() bool {
	if strings.Count(play.numLottery.LotNum, ",") != 11 {
		return false
	}
	var lotNumsArr []string = strings.Split(play.numLottery.LotNum, ",")
	if len(lotNumsArr) != 12 {
		return false
	}

	var ballFlag bool = true
	var betZhu int = 1

	for i := 0; i < 12; i++ {
		var tempPosition []string = strings.Split(lotNumsArr[i], " ")
		if len(tempPosition) == 1 {
			if tempPosition[0] != "3" && tempPosition[0] != "1" && tempPosition[0] != "0" {
				ballFlag = false
			}
		} else {
			for _, ball := range tempPosition {
				if ball != "3" && ball != "1" && ball != "0" {
					ballFlag = false
					break
				}
			}
		}

		if !ballFlag {
			break
		}

		betZhu = betZhu * len(tempPosition)

	}
	if !ballFlag {
		return false
	}

	play.ticketNum = betZhu

	return true
}

//新建半全场对象
func NewPlayBqc() LotteryInterface {
	return &PlayBqc{}
}

func init() {
	Register("bqc", NewPlayBqc)
}
