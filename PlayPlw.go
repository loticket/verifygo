package verifygo

import (
	"strings"
)

//排列5投注
type PlayPlw struct {
	BasePlay //投注基本信息
}

var (
	pBall = ball{1, 10, map[string]uint{"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9}}
)

//验证子玩法和投注方式
func (plw *PlayPlw) CheckPlaytype() bool {
	if plw.numLottery.PlayType != 1 {
		return false
	}

	if plw.numLottery.LotType != 1 && plw.numLottery.LotType != 2 {
		return false
	}

	return true
}

//验证投注的号码是否正确
func (plw *PlayPlw) PlayCheck() bool {

	var checkRes bool = false

	if plw.numLottery.LotType == 1 {
		checkRes = plw.plwSingle()
	} else if plw.numLottery.LotType == 2 {
		checkRes = plw.plwComplex()
	} else {
		checkRes = false
	}

	return checkRes
}

//单式
func (plw *PlayPlw) plwSingle() bool {
	if 2*plw.numLottery.Multiple != plw.numLottery.Money {
		return false
	}
	if strings.Count(plw.numLottery.LotNum, ",") != 4 {
		return false
	}

	if plw.numLottery.BetNum != 1 {
		return false
	}

	var lotNumsArr []string = strings.Split(plw.numLottery.LotNum, ";")
	return plw.BallSame(lotNumsArr, pBall.names)

}

//检查每一注的基本信息

//复试
func (plw *PlayPlw) plwComplex() bool {
	if strings.Count(plw.numLottery.LotNum, ";") != 4 {
		return false
	}

	if plw.numLottery.LotType == 1 && plw.numLottery.BetNum != 1 {
		return false
	}

	if plw.numLottery.LotType == 2 && plw.numLottery.BetNum == 1 {
		return false
	}

	//总共7位数 要判断是不是复试
	var lotNumsArr []string = strings.Split(plw.numLottery.LotNum, ";")
	var lotNums [5]int
	var flag bool = true
	for i := 0; i < 5; i++ {
		var tempPosition []string = strings.Split(lotNumsArr[i], ",")
		lotNums[i] = len(tempPosition)
		flag = plw.BallInt(tempPosition, pBall.names)
		if !flag {
			break
		}
	}

	if !flag {
		return false
	}

	//计算注数
	var allZhu int = 1
	for j := 0; j < 5; j++ {
		allZhu *= lotNums[j]
	}

	plw.ticketNum = allZhu

	return flag
}

func NewPlayPlw() LotteryInterface {
	return &PlayPlw{}
}

func init() {
	Register("plw", NewPlayPlw)
}
