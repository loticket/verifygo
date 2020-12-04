package verifygo

import (
	"errors"
	"strings"
)

type PlayQxc struct {
	BasePlay //投注基本信息
}

var (
	qBall       = ball{1, 10, map[string]uint{"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9}}
	qxcBlueBall = ball{1, 15, map[string]uint{"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9, "10": 10, "11": 11, "12": 12, "13": 13, "14": 14}}
)

//验证投注的号码是否正确
func (pq *PlayQxc) Verification() (bool, error) {
	var flags bool = pq.CheckPlaytype()
	if !flags {
		return false, errors.New("子玩法或者投注方式错误")
	}

	flags = pq.PlayCheck()

	if !flags {
		return false, errors.New("投注号码错误")
	}

	if pq.ticketNum != pq.numLottery.BetNum {
		return false, errors.New("注数计算错误")
	}

	//单注不能超过两万
	if (pq.numLottery.Money / pq.numLottery.Multiple) > 20000 {
		return false, errors.New("单注金额不能超过2万")
	}

	if pq.ticketNum*2*pq.numLottery.Multiple != pq.numLottery.Money {
		return false, errors.New("金额计算错误")
	}

	return true, nil
}

//拆票
func (pq *PlayQxc) GetSpliteTicket() {}

//验证子玩法和投注方式
func (pq *PlayQxc) CheckPlaytype() bool {
	if pq.numLottery.PlayType != 1 {
		return false
	}

	if pq.numLottery.LotType != 1 && pq.numLottery.LotType != 2 {
		return false
	}

	return true
}

//检查每一注的基本信息
func (pq *PlayQxc) PlayCheck() bool {
	var oneLot NumLottery = pq.numLottery
	if strings.EqualFold(oneLot.LotNum, "") {
		return false
	}
	var checkRes bool = false
	if oneLot.LotType == 1 {
		checkRes = pq.qixingcaiSingle()
	} else if oneLot.LotType == 2 {
		checkRes = pq.qixingcaiComplex()
	}

	return checkRes

}

//复试格式验证
func (pq *PlayQxc) qixingcaiSingle() bool {

	if pq.numLottery.BetNum != 1 {
		return false
	}

	if (2 * pq.numLottery.Multiple) != pq.numLottery.Money {
		return false
	}

	return pq.checkBallCommon()
}

//复式检查
func (pq *PlayQxc) qixingcaiComplex() bool {

	if pq.numLottery.BetNum == 1 {
		return false
	}
	return pq.checkBallCommon()
}

//格式总验证
func (pq *PlayQxc) checkBallCommon() bool {

	if strings.Count(pq.numLottery.LotNum, "-") != 1 {
		return false
	}

	var ballArr []string = strings.Split(pq.numLottery.LotNum, "-")

	var blueBallArr []string = strings.Split(ballArr[1], ",")

	var blueBallNum int = len(blueBallArr)

	if blueBallNum < qxcBlueBall.min || blueBallNum > qxcBlueBall.max {
		return false
	}

	if !pq.Ball(blueBallArr, qxcBlueBall.names) {
		return false
	}

	//总共7位数 要判断是不是复试
	var lotNumsArr []string = strings.Split(ballArr[0], ";")
	var lotNums [6]int
	var flag bool = true
	for i := 0; i < 6; i++ {
		var tempPosition []string = strings.Split(lotNumsArr[i], ",")
		lotNums[i] = len(tempPosition)
		flag = pq.Ball(tempPosition, qBall.names)
		if !flag {
			break
		}
	}

	if !flag {
		return flag
	}

	var allZhu int = 1
	for j := 0; j < 6; j++ {
		allZhu *= lotNums[j]
	}

	pq.ticketNum = allZhu * blueBallNum

	return flag
}

func NewPlayQxc() LotteryInterface {
	return &PlayQxc{}
}

func init() {
	Register("qxc", NewPlayQxc)
}
