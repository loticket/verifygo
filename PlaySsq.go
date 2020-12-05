package verifygo

import (
	"errors"
	"github.com/loticket/verifygo/utils"
	"strings"
)

//双色球投注格式验证
type PlaySsq struct {
	BasePlay //投注基本信息
}

//前驱胆拖  胆【1-5】 拖码 【2-32】
//后驱胆拖  【1-16】

var (
	redSBall = ball{6, 33, map[string]uint{"01": 1, "02": 2, "03": 3, "04": 4, "05": 5, "06": 6, "07": 7, "08": 8, "09": 9,
		"10": 10, "11": 11, "12": 12, "13": 13, "14": 14, "15": 15, "16": 16, "17": 17, "18": 18, "19": 19, "20": 20, "21": 21,
		"22": 22, "23": 23, "24": 24, "25": 25, "26": 26, "27": 27, "28": 28, "29": 29, "30": 30, "31": 31, "32": 32, "33": 33}}

	redSBallDan = ball{1, 5, nil}  //前驱胆码
	redSBallTuo = ball{2, 32, nil} //前驱胆码

	blueSBall = ball{1, 16, map[string]uint{"01": 1, "02": 2, "03": 3, "04": 4, "05": 5, "06": 6, "07": 7, "08": 8, "09": 9, "10": 10, "11": 11, "12": 12, "13": 13, "14": 14, "15": 15, "16": 16}}
)

//验证投注的号码是否正确
func (ps *PlaySsq) Verification() (bool, error) {
	var flags bool = ps.CheckPlaytype()
	if !flags {
		return false, errors.New("子玩法或者投注方式错误")
	}

	flags = ps.PlayCheck()

	if !flags {
		return false, errors.New("投注号码错误")
	}

	if ps.ticketNum != ps.numLottery.BetNum {
		return false, errors.New("注数计算错误")
	}

	//单注不能超过两万
	if (ps.numLottery.Money / ps.numLottery.Multiple) > 20000 {
		return false, errors.New("单注金额不能超过2万")
	}

	if ps.ticketNum*2*ps.numLottery.Multiple != ps.numLottery.Money {
		return false, errors.New("金额计算错误")
	}

	return true, nil
}

//验证子玩法和投注方式
func (ps *PlaySsq) CheckPlaytype() bool {
	if ps.numLottery.PlayType != 1 {
		return false
	}

	if ps.numLottery.LotType != 1 && ps.numLottery.LotType != 2 && ps.numLottery.LotType != 3 {
		return false
	}

	return true
}

//检查每一注的基本信息
func (ps *PlaySsq) PlayCheck() bool {
	var oneLot NumLottery = ps.numLottery
	if strings.EqualFold(oneLot.LotNum, "") {
		return false
	}
	var checkRes bool = true
	if oneLot.LotType == 3 {
		checkRes = ps.checkBallDanTuo()
	} else if oneLot.LotType == 2 {
		checkRes = ps.checkBallNormal()
	} else {
		checkRes = ps.checkBallSingle()
	}

	return checkRes

}

//拆票
func (ps *PlaySsq) GetSpliteTicket() []NumLottery {
	var multi []int = ps.spliteMultiple()
	var newTicket []NumLottery = make([]NumLottery, 0)
	var ticket NumLottery = ps.numLottery
	for _, val := range multi {
		ticket.Multiple = val
		ticket.Money = val * ticket.BetNum * 2
		newTicket = append(newTicket, ticket)
	}
	return newTicket
}

//检查单式格式是否正确
func (ps *PlaySsq) checkBallSingle() bool {
	return ps.checkBallNormal()
}

//验证号码是否正确 -单式或者复式
//01,02,03,04,05,06-01
func (ps *PlaySsq) checkBallNormal() bool {
	if strings.Count(ps.numLottery.LotNum, "-") != 1 {
		return false
	}

	var redBallNum int = 0
	var blueBallNum int = 0

	var ballArr []string = strings.Split(ps.numLottery.LotNum, "-")

	var redBallArr = strings.Split(ballArr[0], ",")
	var blueBallArr = strings.Split(ballArr[1], ",")

	redBallNum = len(redBallArr)
	blueBallNum = len(blueBallArr)

	if redBallNum < redSBall.min || redBallNum > redSBall.max {
		return false
	}

	if blueBallNum < blueSBall.min || blueBallNum > blueSBall.max {
		return false
	}

	//红球的投注号码

	if ok := ps.redBall(redBallArr); !ok {
		return false
	}

	//篮球的号码
	if ok := ps.blueBall(blueBallArr); !ok {
		return false
	}

	//验证注数
	var redComb int = utils.Combination(redBallNum, redSBall.min)
	var blueComb int = utils.Combination(blueBallNum, blueSBall.min)

	ps.ticketNum = redComb * blueComb

	if ps.ticketNum != ps.numLottery.BetNum {
		return false
	}

	var nowBetMoney int = ps.ticketNum * 2 * ps.numLottery.Multiple

	if nowBetMoney != ps.numLottery.Money {
		return false
	}
	return true

}

//检查胆拖投注号码的正确性
//01,02|03,04,05,06-01
//前驱胆拖  胆【1-5】 拖码 【2-31】
//后驱【1-16】
func (ps *PlaySsq) checkBallDanTuo() bool {
	if strings.Count(ps.numLottery.LotNum, "-") != 1 {
		return false
	}

	var ballArr []string = strings.Split(ps.numLottery.LotNum, "-")

	//判断前驱的胆拖
	if strings.Count(ballArr[0], "|") != 1 {
		return false
	}

	if !strings.Contains(ballArr[0], ",") {
		return false
	}

	//检查红球投注是否已重复的号码及其投注号码格式的正确性

	var redBallDanTuoArr []string = strings.Split(ballArr[0], "|")

	//红球胆
	var rDanArr []string = strings.Split(redBallDanTuoArr[0], ",") //红球胆码
	if ok := ps.redBall(rDanArr); !ok {
		return false
	}
	//红球拖
	var rTuoArr []string = strings.Split(redBallDanTuoArr[1], ",") //红球拖码
	if ok := ps.redBall(rTuoArr); !ok {
		return false
	}
	var rDanNum int = len(rDanArr)
	var rTuoNum int = len(rTuoArr)
	//红球胆
	if rDanNum < redSBallDan.min || rDanNum > redSBallDan.max {
		return false
	}

	//红球托
	if rTuoNum < redSBallTuo.min || rTuoNum > redSBallTuo.max {
		return false
	}

	//总共红球的个数不能小于6个
	if (rDanNum + rTuoNum) < 7 {
		return false
	}

	var redZhu int = utils.Combination(rTuoNum, (redSBall.min - rDanNum))

	//验证号码的正确性和是否有重复
	var blueBallArr []string = strings.Split(ballArr[1], ",")

	if oks := ps.blueBall(blueBallArr); !oks {
		return false
	}

	var blueNums int = len(blueBallArr)

	if blueNums < blueSBall.min || blueNums > blueSBall.max {
		return false
	}

	//计算注数

	ps.ticketNum = redZhu * blueNums

	return true
}

//检查红球是否在规定的里面
func (ps *PlaySsq) redBall(redBallArr []string) bool {
	return ps.Ball(redBallArr, redSBall.names)
}

//PlaySsq 检测篮球是否在规定的里面
func (ps *PlaySsq) blueBall(blueBallArr []string) bool {
	return ps.Ball(blueBallArr, blueSBall.names)
}

func NewPlaySsq() LotteryInterface {
	return &PlaySsq{}
}

func init() {
	Register("ssq", NewPlaySsq)
}
