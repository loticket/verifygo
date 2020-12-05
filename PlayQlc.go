package verifygo

import (
	"errors"
	"github.com/loticket/verifygo/utils"
	"sort"
	"strings"
)

//七乐彩投注
type PlayQlc struct {
	BasePlay //投注基本信息
}

var (
	qlcBall = ball{7, 30, map[string]uint{"01": 1, "02": 2, "03": 3, "04": 4, "05": 5, "06": 6, "07": 7, "08": 8, "09": 9,
		"10": 10, "11": 11, "12": 12, "13": 13, "14": 14, "15": 15, "16": 16, "17": 17, "18": 18, "19": 19, "20": 20, "21": 21,
		"22": 22, "23": 23, "24": 24, "25": 25, "26": 26, "27": 27, "28": 28, "29": 29, "30": 30}}

	qlcBallDan = ball{1, 6, nil}  //胆码
	qlcBallTuo = ball{2, 29, nil} //拖码
)

//验证投注的号码是否正确
func (plc *PlayQlc) Verification() (bool, error) {
	var flags bool = plc.CheckPlaytype()
	if !flags {
		return false, errors.New("子玩法或者投注方式错误")
	}

	flags = plc.PlayCheck()

	if !flags {
		return false, errors.New("投注号码错误")
	}

	if plc.ticketNum != plc.numLottery.BetNum {
		return false, errors.New("注数计算错误")
	}

	//单注不能超过两万
	if (plc.numLottery.Money / plc.numLottery.Multiple) > 20000 {
		return false, errors.New("单注金额不能超过2万")
	}

	if plc.ticketNum*2*plc.numLottery.Multiple != plc.numLottery.Money {
		return false, errors.New("金额计算错误")
	}

	return true, nil
}

//拆票
func (plc *PlayQlc) GetSpliteTicket() []NumLottery {
	var multi []int = plc.spliteMultiple()
	var newTicket []NumLottery = make([]NumLottery, 0)
	var ticket NumLottery = plc.numLottery
	for _, val := range multi {
		ticket.Multiple = val
		ticket.Money = val * ticket.BetNum * 2
		newTicket = append(newTicket, ticket)
	}
	return newTicket
}

//验证子玩法和投注方式
func (plc *PlayQlc) CheckPlaytype() bool {
	if plc.numLottery.PlayType != 1 {
		return false
	}

	if plc.numLottery.LotType != 1 && plc.numLottery.LotType != 2 && plc.numLottery.LotType != 3 {
		return false
	}

	return true
}

//验证投注的号码是否正确
func (plc *PlayQlc) PlayCheck() bool {
	var checkRes bool = false
	if plc.numLottery.LotType == 3 {
		checkRes = plc.checkBallDanTuo()
	} else if plc.numLottery.LotType == 2 {
		checkRes = plc.checkBallNormal()
	} else {
		checkRes = plc.checkBallNormal()
	}

	return checkRes
}

//验证号码是否正确
//01,02,03,04,05,06,07
func (plc *PlayQlc) checkBallNormal() bool {
	if strings.Count(plc.numLottery.LotNum, ",") < 6 {
		return false
	}

	var redBallNum int = 0

	var ballArr []string = strings.Split(plc.numLottery.LotNum, ",")

	redBallNum = len(ballArr)

	if redBallNum < qlcBall.min || redBallNum > qlcBall.max {
		return false
	}

	if plc.numLottery.LotType == 1 && redBallNum != 7 {
		return false
	}

	if plc.numLottery.LotType == 2 && redBallNum < 8 {
		return false
	}

	//红球的投注号码

	if ok := plc.Ball(ballArr, qlcBall.names); !ok {
		return false
	}

	//验证注数
	plc.ticketNum = utils.Combination(redBallNum, qlcBall.min)

	return true

}

//检查胆拖投注号码的正确性
//02,03,04,05|01,06,07,09,17,21,26,28,30
//前驱胆拖  胆【1-6】
//后驱【2-28】
func (plc *PlayQlc) checkBallDanTuo() bool {
	if strings.Count(plc.numLottery.LotNum, "|") != 1 {
		return false
	}

	var ballArr []string = strings.Split(plc.numLottery.LotNum, "|")

	//判断胆码[1-6]

	if !strings.Contains(ballArr[1], ",") {
		return false
	}

	//检查红球投注是否已重复的号码及其投注号码格式的正确性

	var ballAll string = strings.Replace(plc.numLottery.LotNum, "|", ",", -1)

	var ballAllArr []string = strings.Split(ballAll, ",")
	sort.Sort(utils.StringNumArr(ballAllArr))

	if ok := plc.Ball(ballAllArr, qlcBall.names); !ok {
		return false
	}

	//红球胆
	var rDanArr []string = strings.Split(ballArr[0], ",") //红球胆码
	if ok := plc.Ball(rDanArr, qlcBall.names); !ok {
		return false
	}

	var rTuoArr []string = strings.Split(ballArr[1], ",") //红球拖码
	if ok := plc.Ball(rTuoArr, qlcBall.names); !ok {
		return false
	}

	var rDanNum int = len(rDanArr)
	var rTuoNum int = len(rTuoArr)
	//红球胆
	if rDanNum < qlcBallDan.min || rDanNum > qlcBallDan.max {
		return false
	}

	//红球托
	if rTuoNum < qlcBallTuo.min || rTuoNum > qlcBallTuo.max {
		return false
	}

	//总共红球的个数不能小于8个
	if (rDanNum + rTuoNum) < 8 {
		return false
	}

	//计算注数

	plc.ticketNum = utils.Combination(rTuoNum, (qlcBall.min - rDanNum))

	return true
}

func NewPlayQlc() LotteryInterface {
	return &PlayQlc{}
}

func init() {
	Register("qlc", NewPlayQlc)
}
