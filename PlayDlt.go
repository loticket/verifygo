package verifygo

import (
	"github.com/loticket/verifygo/utils"
	"sort"
	"strings"
)

//大乐透投注信息
type PlayDlt struct {
	BasePlay //投注基本信息
}

//前驱胆拖  胆【1-4】 拖码 【2-33】
//后驱胆拖  胆【0-1】 拖码 【2-10】

var (
	redBall = ball{5, 35, map[string]uint{"01": 1, "02": 2, "03": 3, "04": 4, "05": 5, "06": 6, "07": 7, "08": 8, "09": 9,
		"10": 10, "11": 11, "12": 12, "13": 13, "14": 14, "15": 15, "16": 16, "17": 17, "18": 18, "19": 19, "20": 20, "21": 21,
		"22": 22, "23": 23, "24": 24, "25": 25, "26": 26, "27": 27, "28": 28, "29": 29, "30": 30, "31": 31, "32": 32, "33": 33, "34": 34, "35": 35,
	}}

	redBallDan = ball{1, 4, nil}  //前驱胆码
	redBallTuo = ball{2, 33, nil} //前驱胆码

	blueBall    = ball{2, 12, map[string]uint{"01": 1, "02": 2, "03": 3, "04": 4, "05": 5, "06": 6, "07": 7, "08": 8, "09": 9, "10": 10, "11": 11, "12": 12}}
	blueBallDan = ball{0, 1, nil}  //后驱胆码
	blueBallTuo = ball{2, 10, nil} //后驱拖码
)

//验证子玩法和投注方式
func (pd *PlayDlt) CheckPlaytype() bool {
	if pd.numLottery.PlayType != 1 && pd.numLottery.PlayType != 2 {
		return false
	}

	if pd.numLottery.LotType != 1 && pd.numLottery.LotType != 2 && pd.numLottery.LotType != 3 {
		return false
	}

	return true
}

//验证投注的号码是否正确
func (pd *PlayDlt) PlayCheck() bool {
	var checkRes bool = true
	if pd.numLottery.LotType == 3 {
		checkRes = pd.checkBallDanTuo()
	} else if pd.numLottery.LotType == 2 {
		checkRes = pd.checkBallNormal()
	} else if pd.numLottery.LotType == 1 {
		checkRes = pd.checkBallSingle()
	}

	return checkRes
}

//验证号码是否正确--复式
//01,02,03,04,05-01,02
func (pd *PlayDlt) checkBallNormal() bool {
	if strings.Count(pd.numLottery.LotNum, "-") != 1 {
		return false
	}
	//判断是单式或者复试
	if pd.numLottery.BetNum == 1 {
		return false
	}

	var redBallNum int = 0
	var blueBallNum int = 0

	var ballArr []string = strings.Split(pd.numLottery.LotNum, "-")

	var redBallArr = strings.Split(ballArr[0], ",")
	var blueBallArr = strings.Split(ballArr[1], ",")

	redBallNum = len(redBallArr)
	blueBallNum = len(blueBallArr)

	if redBallNum < redBall.min || redBallNum > redBall.max {
		return false
	}

	if blueBallNum < blueBall.min || blueBallNum > blueBall.max {
		return false
	}

	//红球的投注号码

	if ok := pd.redBall(redBallArr); !ok {
		return false
	}

	//篮球的号码
	if ok := pd.blueBall(blueBallArr); !ok {
		return false
	}

	//验证注数
	var redComb int = utils.Combination(redBallNum, redBall.min)
	var blueComb int = utils.Combination(blueBallNum, blueBall.min)
	pd.ticketNum = redComb * blueComb

	return true

}

//单式
func (pd *PlayDlt) checkBallSingle() bool {
	if strings.Count(pd.numLottery.LotNum, "-") != 1 {
		return false
	}
	//判断是单式或者复试
	if pd.numLottery.BetNum != 1 {
		return false
	}

	var redBallNum int = 0
	var blueBallNum int = 0

	var ballArr []string = strings.Split(pd.numLottery.LotNum, "-")

	var redBallArr = strings.Split(ballArr[0], ",")
	var blueBallArr = strings.Split(ballArr[1], ",")

	redBallNum = len(redBallArr)
	blueBallNum = len(blueBallArr)

	if redBallNum != redBall.min {
		return false
	}

	if blueBallNum != blueBall.min {
		return false
	}

	//红球的投注号码

	if ok := pd.redBall(redBallArr); !ok {
		return false
	}

	//篮球的号码
	if ok := pd.blueBall(blueBallArr); !ok {
		return false
	}

	pd.ticketNum = 1

	return true
}

//检查胆拖投注号码的正确性
//01,02|03,04,05,06-01|02,03
//前驱胆拖  胆【1-4】 拖码 【2-33】
//后驱胆拖  胆【0-1】 拖码 【2-10】
func (pd *PlayDlt) checkBallDanTuo() bool {
	if strings.Count(pd.numLottery.LotNum, "-") != 1 {
		return false
	}

	var ballArr []string = strings.Split(pd.numLottery.LotNum, "-")
	var redDanTuoFlag bool = true
	var blueDanTuoFlag bool = true
	//判断前驱的胆拖
	if strings.Count(ballArr[0], "|") != 1 {
		redDanTuoFlag = false
	}

	if strings.Count(ballArr[1], "|") != 1 {
		blueDanTuoFlag = false
	}

	if !redDanTuoFlag && !blueDanTuoFlag {
		return false
	}

	if !strings.Contains(ballArr[0], ",") {
		return false
	}

	//检查红球投注是否已重复的号码及其投注号码格式的正确性

	var redBallAll string = strings.Replace(ballArr[0], "|", ",", -1)

	var redBallAllArr []string = strings.Split(redBallAll, ",")

	sort.Strings(redBallAllArr)

	if ok := pd.redBall(redBallAllArr); !ok {
		return false
	}

	var rDanArr []string
	var rTuoArr []string
	//判断红球是否有胆拖
	if redDanTuoFlag {
		var redBallDanTuoArr []string = strings.Split(ballArr[0], "|")
		//红球胆
		rDanArr = strings.Split(redBallDanTuoArr[0], ",") //红球胆码
		rTuoArr = strings.Split(redBallDanTuoArr[1], ",") //红球拖码
	} else {
		rTuoArr = strings.Split(ballArr[0], ",")
		rDanArr = []string{}
	}

	var rDanNum int = len(rDanArr)
	var rTuoNum int = len(rTuoArr)
	//红球胆
	if (rDanNum < redBallDan.min || rDanNum > redBallDan.max) && redDanTuoFlag {
		return false
	}

	//红球托
	if (rTuoNum < redBallTuo.min || rTuoNum > redBallTuo.max) && redDanTuoFlag {
		return false
	}

	//总共红球的个数不能小于6个
	if (rDanNum + rTuoNum) < 6 {
		return false
	}

	//验证篮球拖码
	if !strings.Contains(ballArr[1], ",") {
		return false
	}

	var redZhu int = utils.Combination(rTuoNum, (redBall.min - rDanNum))

	//篮球是否有胆码
	var blueDanOr bool = blueDanTuoFlag
	var blueDanNum int = 0
	var blueTuoNum int = 0
	var blueDanTuoStr string = ballArr[1]
	if blueDanOr {
		blueDanTuoStr = strings.Replace(blueDanTuoStr, "|", ",", -1)
	}

	//验证号码的正确性和是否有重复
	var blueDanTuoStrArr []string = strings.Split(blueDanTuoStr, ",")
	sort.Strings(blueDanTuoStrArr)

	if oks := pd.blueBall(blueDanTuoStrArr); !oks {
		return false
	}

	//胆码的个数及其拖码的个数
	var blueBallDanArr []string
	var blueBallTuoArr []string
	if blueDanOr {
		blueBallDanTuoAll := strings.Split(ballArr[1], "|")
		blueBallDanArr = strings.Split(blueBallDanTuoAll[0], ",")
		blueBallTuoArr = strings.Split(blueBallDanTuoAll[1], ",")
	} else {
		blueBallTuoArr = strings.Split(ballArr[1], ",")

	}

	blueDanNum = len(blueBallDanArr)
	blueTuoNum = len(blueBallTuoArr)

	if blueDanNum > blueBallDan.max {
		return false
	}

	if (blueTuoNum < blueBallTuo.min || blueTuoNum > blueBallTuo.max) && blueDanOr {
		return false
	}

	//计算注数
	var blueZhu int = utils.Combination(blueTuoNum, (blueBall.min - blueDanNum))

	pd.ticketNum = redZhu * blueZhu

	return true
}

//检查红球是否在规定的里面
func (pd *PlayDlt) redBall(redBallArr []string) bool {
	return pd.Ball(redBallArr, redBall.names)
}

//检测篮球是否在规定的里面
func (pd *PlayDlt) blueBall(blueBallArr []string) bool {
	return pd.Ball(blueBallArr, blueBall.names)
}

func NewPlayDlt() LotteryInterface {
	return &PlayDlt{}
}

func init() {
	Register("dlt", NewPlayDlt)
}
