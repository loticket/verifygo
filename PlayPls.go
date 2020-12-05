package verifygo

import (
	"errors"
	"github.com/loticket/verifygo/utils"
	"sort"
	"strconv"
	"strings"
)

//排列3投注
type PlayPls struct {
	BasePlay //投注基本信息
}

var (
	zxSum     = [28]int{1, 3, 6, 10, 15, 21, 28, 36, 45, 55, 63, 69, 73, 75, 75, 73, 69, 63, 55, 45, 36, 28, 21, 15, 10, 6, 3, 1} //直选和值
	z3Sum     = [27]int{0, 1, 2, 1, 3, 3, 3, 4, 5, 4, 5, 5, 4, 5, 5, 4, 5, 5, 4, 5, 4, 3, 3, 3, 1, 2, 1}
	z6Sum     = [25]int{0, 0, 0, 1, 1, 2, 3, 4, 5, 7, 8, 9, 10, 10, 10, 10, 9, 8, 7, 5, 4, 3, 2, 1, 1}
	plsBall   = ball{1, 10, map[string]uint{"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9}}
	zxSumBall = ball{1, 28, map[string]uint{"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
		"10": 10, "11": 11, "12": 12, "13": 13, "14": 14, "15": 15, "16": 16, "17": 17, "18": 18, "19": 19, "20": 20, "21": 21,
		"22": 22, "23": 23, "24": 24, "25": 25, "26": 26, "27": 27}}
	z3SumBall = ball{1, 26, map[string]uint{"1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
		"10": 10, "11": 11, "12": 12, "13": 13, "14": 14, "15": 15, "16": 16, "17": 17, "18": 18, "19": 19, "20": 20, "21": 21,
		"22": 22, "23": 23, "24": 24, "25": 25, "26": 26}}
	z6SumBall = ball{1, 24, map[string]uint{"3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9,
		"10": 10, "11": 11, "12": 12, "13": 13, "14": 14, "15": 15, "16": 16, "17": 17, "18": 18, "19": 19, "20": 20, "21": 21,
		"22": 22, "23": 23, "24": 24}}
	z3Ball = ball{2, 10, nil}
	z6ball = ball{3, 10, nil}
)

//验证投注的号码是否正确
func (pls *PlayPls) Verification() (bool, error) {
	var flags bool = pls.CheckPlaytype()
	if !flags {
		return false, errors.New("子玩法或者投注方式错误")
	}

	flags = pls.PlayCheck()

	if !flags {
		return false, errors.New("投注号码错误")
	}

	if pls.ticketNum != pls.numLottery.BetNum {
		return false, errors.New("注数计算错误")
	}

	//单注不能超过两万
	if (pls.numLottery.Money / pls.numLottery.Multiple) > 20000 {
		return false, errors.New("单注金额不能超过2万")
	}

	var oneMoney int = 2
	if pls.numLottery.PlayType == 2 {
		oneMoney = 3
	}

	if pls.ticketNum*oneMoney*pls.numLottery.Multiple != pls.numLottery.Money {
		return false, errors.New("金额计算错误")
	}

	return true, nil
}

//验证投注的号码是否正确
func (pls *PlayPls) PlayCheck() bool {

	var checkRes bool = false
	var oneLot NumLottery = pls.numLottery
	if oneLot.LotType == 1 { //直选单式
		checkRes = pls.zhixuanSingle()
	} else if oneLot.LotType == 2 { //直选复式
		checkRes = pls.zhixuanComplex()
	} else if oneLot.LotType == 3 { //直选和值
		checkRes = pls.zhixuanSum()
	} else if oneLot.LotType == 4 { //组三单式
		checkRes = pls.z3Single()
	} else if oneLot.LotType == 5 { //组三复式
		checkRes = pls.z3Complex()
	} else if oneLot.LotType == 6 { //组六和值
		checkRes = pls.z3Sum()
	} else if oneLot.LotType == 7 { //组六单式
		checkRes = pls.z6Single()
	} else if oneLot.LotType == 8 { //组六复式
		checkRes = pls.z6Complex()
	} else if oneLot.LotType == 9 { //组六和值
		checkRes = pls.z6Sum()
	} else if oneLot.LotType == 10 { //组三组六和值
		checkRes = pls.z6z3Sum()
	} else if oneLot.LotType == 11 { //组三胆拖
		checkRes = pls.zu3DanTuo()
	} else if oneLot.LotType == 12 { //组六胆拖
		checkRes = pls.zu6DanTuo()
	} else {
		checkRes = false
	}

	return checkRes
}

//检查子玩法和投注方式是否正确

func (pls *PlayPls) CheckPlaytype() bool {
	if pls.numLottery.PlayType != 1 && pls.numLottery.PlayType != 2 && pls.numLottery.PlayType != 3 && pls.numLottery.PlayType != 4 {
		return false
	}

	if pls.numLottery.LotType != 1 && pls.numLottery.LotType != 2 && pls.numLottery.LotType != 3 && pls.numLottery.LotType != 4 && pls.numLottery.LotType != 5 && pls.numLottery.LotType != 6 && pls.numLottery.LotType != 7 && pls.numLottery.LotType != 8 && pls.numLottery.LotType != 9 && pls.numLottery.LotType != 10 {
		return false
	}

	return true

}

//拆票
func (pls *PlayPls) GetSpliteTicket() []NumLottery {
	var multi []int = pls.spliteMultiple()
	var oneMoney int = 2
	var newTicket []NumLottery = make([]NumLottery, 0)
	var ticket NumLottery = pls.numLottery
	for _, val := range multi {
		ticket.Multiple = val
		ticket.Money = val * ticket.BetNum * oneMoney
		newTicket = append(newTicket, ticket)
	}
	return newTicket
}

//直选单式
//格式为 1,2,3
func (pls *PlayPls) zhixuanSingle() bool {
	if strings.Count(pls.numLottery.LotNum, ";") != 2 {
		return false
	}
	var lotNumsArr []string = strings.Split(pls.numLottery.LotNum, ";")
	var ballFlag bool = pls.BallSame(lotNumsArr, plsBall.names)

	if !ballFlag {
		return false
	}

	if pls.numLottery.BetNum != 1 {
		return false
	}

	pls.ticketNum = 1

	return true
}

//直选复试
//格式为：1 2,2 3,3 4
func (pls *PlayPls) zhixuanComplex() bool {
	if strings.Count(pls.numLottery.LotNum, ";") != 2 {
		return false
	}

	//总共3位数 要判断是不是复试
	var lotNumsArr []string = strings.Split(pls.numLottery.LotNum, ";")
	var lotNums [3]int
	var flag bool = true
	for i := 0; i < 3; i++ {
		var tempPosition []string = strings.Split(lotNumsArr[i], ",")
		lotNums[i] = len(tempPosition)
		flag = pls.Ball(tempPosition, qBall.names)
		if !flag {
			break
		}
	}

	if !flag {
		return false
	}

	//计算注数
	var allZhu int = 1
	for j := 0; j < 3; j++ {
		allZhu *= lotNums[j]
	}

	pls.ticketNum = allZhu

	return flag
}

//直选和值
//格式为：05
func (pls *PlayPls) zhixuanSum() bool {
	var lotNumsArr []string = strings.Split(pls.numLottery.LotNum, ",")
	var ballNums = len(lotNumsArr)
	if ballNums < zxSumBall.min || ballNums > zxSumBall.max {
		return false
	}

	var ballFlag bool = pls.BallInt(lotNumsArr, zxSumBall.names)
	if !ballFlag {
		return false
	}

	//计算注数
	var betZhu int = 0
	for i := 0; i < ballNums; i++ {
		index, _ := strconv.Atoi(lotNumsArr[i])
		betZhu += zxSum[index]
	}

	pls.ticketNum = betZhu

	return true
}

//组三单式
//格式为：2,2,3
func (pls *PlayPls) z3Single() bool {
	if strings.Count(pls.numLottery.LotNum, ",") != 2 {
		return false
	}

	if 2*pls.numLottery.Multiple != pls.numLottery.Money {
		return false
	}

	var flag bool = true
	var lotNumsArr []string = strings.Split(pls.numLottery.LotNum, ",")

	var ballsNum int = len(lotNumsArr)

	if ballsNum != 3 {
		return false
	}

	if ok := pls.BallSameSort(lotNumsArr, plsBall.names); !ok {
		return false
	}

	if !flag {
		return false
	}

	if !strings.EqualFold(lotNumsArr[0], lotNumsArr[1]) && !strings.EqualFold(lotNumsArr[1], lotNumsArr[2]) {
		return false
	}

	pls.ticketNum = 1

	return true
}

//组三复试
//格式为：1,2,3
func (pls *PlayPls) z3Complex() bool {
	var lotNumsArr []string = strings.Split(pls.numLottery.LotNum, ",")
	var ballNums = len(lotNumsArr)
	if ballNums < 2 || ballNums > 10 {
		return false
	}

	var ballFlag bool = pls.BallInt(lotNumsArr, plsBall.names)

	if !ballFlag {
		return false
	}
	//计算注数
	var betZhu int = utils.Combination(ballNums, 2)
	pls.ticketNum = betZhu * 2

	return true

}

//组三和值
//02,03
func (pls *PlayPls) z3Sum() bool {
	var lotNumsArr []string = strings.Split(pls.numLottery.LotNum, ",")
	var ballNums = len(lotNumsArr)
	if ballNums < z3SumBall.min || ballNums > z3SumBall.max {
		return false
	}

	//有没有重复的数字
	var ballFlag bool = pls.BallInt(lotNumsArr, z3SumBall.names)
	if !ballFlag {
		return false
	}

	//计算注数
	var betZhu int = 0
	for i := 0; i < ballNums; i++ {
		index, _ := strconv.Atoi(lotNumsArr[i])
		betZhu += z3Sum[index]
	}
	pls.ticketNum = betZhu

	return true
}

//组六单式 必须为3个号码 并且三个号码都不相同
func (pls *PlayPls) z6Single() bool {

	if 2*pls.numLottery.Multiple != pls.numLottery.Money {
		return false
	}
	if strings.Count(pls.numLottery.LotNum, ",") != 2 {
		return false
	}
	var lotNumsArr []string = strings.Split(pls.numLottery.LotNum, ",")

	var ballFlag bool = pls.BallInt(lotNumsArr, plsBall.names)
	if !ballFlag {
		return false
	}

	if pls.numLottery.BetNum != 1 {
		return false
	}

	if lotNumsArr[0] == lotNumsArr[1] || lotNumsArr[1] == lotNumsArr[2] {
		return false
	}

	pls.ticketNum = 1

	return true
}

//组六复试 -- 号码都不相同
func (pls *PlayPls) z6Complex() bool {
	if strings.Count(pls.numLottery.LotNum, ",") < 2 {
		return false
	}
	//投注数量
	var lotNumsArr []string = strings.Split(pls.numLottery.LotNum, ",")
	var ballNums = len(lotNumsArr)
	if ballNums < z6ball.min || ballNums > z6ball.max {
		return false
	}
	var ballFlag bool = pls.BallInt(lotNumsArr, plsBall.names)
	if !ballFlag {
		return false
	}

	//计算注数
	pls.ticketNum = utils.Combination(ballNums, z6ball.min)

	return true
}

//组六和值
func (pls *PlayPls) z6Sum() bool {
	var lotNumsArr []string = strings.Split(pls.numLottery.LotNum, ",")
	var ballNums = len(lotNumsArr)
	if ballNums < z6SumBall.min || ballNums > z6SumBall.max {
		return false
	}

	//有没有重复的数字
	var ballFlag bool = pls.BallInt(lotNumsArr, z6SumBall.names)
	if !ballFlag {
		return false
	}

	//计算注数
	var betZhu int = 0
	for i := 0; i < ballNums; i++ {
		index, _ := strconv.Atoi(lotNumsArr[i])
		betZhu += z6Sum[index]
	}

	pls.ticketNum = betZhu

	return true
}

//组三组六和值
func (pls *PlayPls) z6z3Sum() bool {
	var lotNumsArr []string = strings.Split(pls.numLottery.LotNum, ",")
	var ballNums int = len(lotNumsArr)
	if ballNums == 0 {
		return false
	}

	//有没有重复的数字
	var ballFlag bool = pls.BallInt(lotNumsArr, z3SumBall.names)
	if !ballFlag {
		return false
	}

	var zuSanBetAll int = 0
	var zuLiuBetAll int = 0

	for i := 0; i < ballNums; i++ {
		index, _ := strconv.Atoi(lotNumsArr[i])
		zuSanBetAll += z3Sum[index]
		if index > 2 && index < 25 {
			zuLiuBetAll += z6Sum[index]
		}

	}

	pls.ticketNum = zuSanBetAll + zuLiuBetAll

	return true
}

//组三胆拖
//2|3,4,5
//前驱胆拖  胆【1】
//拖【2,9】
func (pls *PlayPls) zu3DanTuo() bool {
	var ballRange []int = []int{1, 1, 2, 9}
	flag, nums := pls.zuDanTuo(ballRange)
	if !flag {
		return flag
	}

	//计算注数

	pls.ticketNum = nums[1] * 2

	return true
}

//组六胆拖验证
func (pls *PlayPls) zu6DanTuo() bool {
	var ballRange []int = []int{1, 2, 2, 9}
	flag, nums := pls.zuDanTuo(ballRange)
	if !flag {
		return flag
	}

	//计算注数

	pls.ticketNum = utils.Combination(nums[1], (3 - nums[0]))

	return true
}

//组三 组六胆拖验证
func (pls *PlayPls) zuDanTuo(dantuo []int) (bool, []int) {
	if strings.Count(pls.numLottery.LotNum, "|") != 1 {
		return false, []int{}
	}

	var ballArr []string = strings.Split(pls.numLottery.LotNum, "|")

	//判断胆码[1-6]

	if !strings.Contains(ballArr[1], ",") {
		return false, []int{}
	}

	//检查红球投注是否已重复的号码及其投注号码格式的正确性

	var ballAll string = strings.Replace(pls.numLottery.LotNum, "|", ",", -1)

	var ballAllArr []string = strings.Split(ballAll, ",")
	sort.Sort(utils.StringNumArr(ballAllArr))

	if ok := pls.Ball(ballAllArr, plsBall.names); !ok {
		return false, []int{}
	}

	//红球胆zu6DanTuo
	var rDanArr []string = strings.Split(ballArr[0], ",") //红球胆码
	if ok := pls.Ball(rDanArr, plsBall.names); !ok {
		return false, []int{}
	}

	var rTuoArr []string = strings.Split(ballArr[1], ",") //红球拖码
	if ok := pls.Ball(rTuoArr, plsBall.names); !ok {
		return false, []int{}
	}

	var rDanNum int = len(rDanArr)
	var rTuoNum int = len(rTuoArr)

	//胆
	if rDanNum < dantuo[0] || rDanNum > dantuo[1] {
		return false, []int{}
	}

	//红球托
	if rTuoNum < dantuo[2] || rTuoNum > dantuo[3] {
		return false, []int{}
	}

	return true, []int{rDanNum, rTuoNum}
}

//排列组三组六和值 主要拆单
func (pls *PlayPls) GetNumLottery() []NumLottery {
	var betAll []NumLottery = make([]NumLottery, 0)

	if pls.numLottery.LotType == 10 { //组三组六复试
		tempLotConn := pls.z6z3SumSplite()
		betAll = append(betAll, tempLotConn...)
	} else {
		betAll = append(betAll, pls.numLottery)
	}

	return betAll
}

//组三组六复试拆票
func (pls *PlayPls) z6z3SumSplite() []NumLottery {
	var betAll []NumLottery = make([]NumLottery, 0)
	var lotNumsArr []string = strings.Split(pls.numLottery.LotNum, ",")
	var ballNums int = len(lotNumsArr)
	var zuSan []string = make([]string, 0)
	var zuLiu []string = make([]string, 0)
	var zuSanBetAll int = 0
	var zuLiuBetAll int = 0
	for i := 0; i < ballNums; i++ {
		var ball string = lotNumsArr[i]
		index, _ := strconv.Atoi(ball)
		zuSanBetAll += z3Sum[index]
		zuSan = append(zuSan, ball)

		if index > 2 && index < 25 {
			zuLiu = append(zuLiu, ball)
			zuLiuBetAll += z6Sum[index]
		}

	}

	if len(zuSan) > 0 { //组三和值
		var tempNumLottery NumLottery = NumLottery{}
		tempNumLottery.Multiple = pls.numLottery.Multiple
		tempNumLottery.BetNum = zuSanBetAll
		tempNumLottery.Money = zuSanBetAll * pls.numLottery.Multiple * 2
		tempNumLottery.PlayType = 2
		tempNumLottery.LotType = 8
		tempNumLottery.LotNum = strings.Join(zuSan, ",")
		betAll = append(betAll, tempNumLottery)
	}

	if len(zuLiu) > 0 { //组六和值
		var tempNumLottery NumLottery = NumLottery{}
		tempNumLottery.Multiple = pls.numLottery.Multiple
		tempNumLottery.BetNum = zuLiuBetAll
		tempNumLottery.Money = zuLiuBetAll * pls.numLottery.Multiple * 2
		tempNumLottery.PlayType = 2
		tempNumLottery.LotType = 9
		tempNumLottery.LotNum = strings.Join(zuLiu, ",")
		betAll = append(betAll, tempNumLottery)
	}

	return betAll

}

func NewPlayPls() LotteryInterface {
	return &PlayPls{}
}

func init() {
	Register("pls", NewPlayPls)
}
