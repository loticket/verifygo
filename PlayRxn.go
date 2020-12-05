package verifygo

import (
	"errors"
	"github.com/loticket/verifygo/utils"
	"regexp"
	"sort"
	"strings"
)

//任选九
type PlayRxn struct {
	BasePlay     //投注基本信息
	spliteTicket []NumLottery
}

//验证投注的号码是否正确
func (play *PlayRxn) Verification() (bool, error) {
	var flags bool = play.CheckPlaytype()
	if !flags {
		return false, errors.New("子玩法或者投注方式错误")
	}

	flags = play.PlayCheck()

	if !flags {
		return false, errors.New("投注号码错误")
	}

	if play.ticketNum != play.numLottery.BetNum {
		return false, errors.New("注数计算错误")
	}

	//单注不能超过两万
	if (play.numLottery.Money / play.numLottery.Multiple) > 20000 {
		return false, errors.New("单注金额不能超过2万")
	}

	if play.ticketNum*2*play.numLottery.Multiple != play.numLottery.Money {
		return false, errors.New("金额计算错误")
	}

	return true, nil
}

//验证子玩法和投注方式
func (play *PlayRxn) CheckPlaytype() bool {
	if play.numLottery.PlayType != 1 {
		return false
	}

	if play.numLottery.LotType != 1 && play.numLottery.LotType != 2 && play.numLottery.LotType != 3 {
		return false
	}

	return true
}

//检查每一注的基本信息
func (play *PlayRxn) PlayCheck() bool {
	var checkRes bool = false
	if play.numLottery.LotType == 1 { //单式
		checkRes = play.checkBallSingle()
	} else if play.numLottery.LotType == 2 { //复试
		checkRes = play.checkBallReexamine()
	} else if play.numLottery.LotType == 3 { //胆拖
		checkRes = play.checkBallDanTuo()
	}
	return checkRes
}

//拆票
func (play *PlayRxn) GetSpliteTicket() []NumLottery {
	var multi []int = play.spliteMultiple()
	var oneMoney int = 2
	var newTicket []NumLottery = make([]NumLottery, 0)
	for i := 0; i < len(play.spliteTicket); i++ {
		var tempTikcket NumLottery = play.spliteTicket[i]
		for _, val := range multi {
			tempTikcket.Multiple = val
			tempTikcket.Money = val * tempTikcket.BetNum * oneMoney
			newTicket = append(newTicket, tempTikcket)
		}
	}

	return newTicket
}

//单式检查
func (play *PlayRxn) checkBallSingle() bool {

	if strings.Count(play.numLottery.LotNum, ",") != 13 {
		return false
	}

	if play.numLottery.BetNum != 1 {
		return false
	}

	if (2 * play.numLottery.Multiple) != play.numLottery.Money {
		return false
	}

	var lotNumsArr []string = strings.Split(play.numLottery.LotNum, ",")
	var betNum int = 0
	for i := 0; i < 14; i++ {
		if lotNumsArr[i] != "3" && lotNumsArr[i] != "1" && lotNumsArr[i] != "0" && lotNumsArr[i] != "#" {
			return false
			break
		}

		if lotNumsArr[i] == "3" || lotNumsArr[i] == "1" || lotNumsArr[i] == "0" {
			betNum++
		}
	}

	if betNum != 9 {
		return false
	}

	play.ticketNum = 1
	play.spliteTicket = []NumLottery{play.numLottery}

	return true
}

//复试检查 限制只能选9场比赛
func (play *PlayRxn) checkBallReexamine() bool {
	if strings.Count(play.numLottery.LotNum, ",") != 13 {
		return false
	}

	var lotNumsArr []string = strings.Split(play.numLottery.LotNum, ",")
	if len(lotNumsArr) != 14 {
		return false
	}

	var ballFlag bool = true
	var athers int = 0
	var CalZhuAll []int = make([]int, 0)
	for i := 0; i < 14; i++ {
		if lotNumsArr[i] == "#" {
			athers++
		} else {
			var tempPosition []string = strings.Split(lotNumsArr[i], " ")
			CalZhuAll = append(CalZhuAll, len(tempPosition))
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

	}

	//拆成单式计算注数

	if athers > 5 {
		return false
	}

	if !ballFlag {
		return false
	}
	play.ticketNum = 1
	if athers == 5 {
		for i := 0; i < 9; i++ {
			play.ticketNum *= CalZhuAll[i]
		}
	} else {
		play.ticketNum = play.ballReexamineSplite(lotNumsArr)
	}

	return true
}

//复试检查 限制智能
func (play *PlayRxn) checkBallDanTuo() bool {
	if strings.Count(play.numLottery.LotNum, ",") != 13 {
		return false
	}
	var lotNumsArr []string = strings.Split(play.numLottery.LotNum, ",")
	if len(lotNumsArr) != 14 {
		return false
	}

	var freeNum int = strings.Count(play.numLottery.LotNum, "#")

	if freeNum > 4 {
		return false
	}

	var freeNumAll int = 0
	var danKey []int = make([]int, 0)
	var danNum map[int]int = make(map[int]int)
	var ballFlag bool = true
	var CalZhu []int = make([]int, 0)
	var CalZhuAll map[int]int = make(map[int]int)
	for k, bet := range lotNumsArr {
		if bet == "#" {
			freeNumAll++
			continue
		}

		var tempPosition []string
		reg, _ := regexp.Compile(`\(.*\)`)
		flag := reg.Match([]byte(bet))
		if flag {
			rs := []rune(bet)
			var ends int = len(rs) - 1
			tempPosition = strings.Split(string(rs[1:ends]), " ")
			danKey = append(danKey, k)
			danNum[k] = len(tempPosition)
		} else {
			tempPosition = strings.Split(bet, " ")
			CalZhu = append(CalZhu, k)
		}

		CalZhuAll[k] = len(tempPosition)

		for _, ball := range tempPosition {
			if ball != "3" && ball != "1" && ball != "0" {
				ballFlag = false
				break
			}
		}
	}

	if !ballFlag {
		return false
	}

	if freeNumAll > 4 {
		return false
	}

	//拆票并计算注数
	play.ballDanTuoSplite()

	return true
}

//老足彩任选九复试拆票
func (play *PlayRxn) ballReexamineSplite(lotNumsArr []string) int {
	var betAll []NumLottery = make([]NumLottery, 0)

	//拆票处理
	var danNums []int = make([]int, 0)
	var tuoNums []int = make([]int, 0)
	var allNums []int = make([]int, 0)
	var betNums []int = make([]int, 14)

	for i := 0; i < 14; i++ {
		reg, _ := regexp.Compile(`\(.*\)`)
		var flag bool = reg.MatchString(lotNumsArr[i])
		if flag {
			danNums = append(danNums, i)
		} else if lotNumsArr[i] != "#" {
			tuoNums = append(tuoNums, i)
		}

		if lotNumsArr[i] != "#" {
			allNums = append(allNums, i)
		}

		if lotNumsArr[i] == "#" {
			betNums[i] = 0
		} else {
			var betNum []string = strings.Split(lotNumsArr[i], " ")
			betNums[i] = len(betNum)
		}

	}

	//组合票处理
	var index int = 9 - len(danNums)
	zhus := utils.ZuheInt{
		Nums: tuoNums,
		Nnum: len(tuoNums),
		Mnum: index,
	}

	var ticketInts [][]int = zhus.ZuheResults()
	var allZhus int = 0
	for j := 0; j < len(ticketInts); j++ {
		temp := ticketInts[j]
		if len(danNums) > 0 {
			temp = append(temp, danNums...)
		}

		sort.Ints(temp)
		//计算奖金
		var tempTiclet []string = make([]string, 14)
		var indexCode int = 0
		var indexZhu int = 1
		for t := 0; t < 14; t++ {
			if temp[indexCode] == t {
				if betNums[temp[indexCode]] != 0 {
					indexZhu = indexZhu * betNums[temp[indexCode]]
				}
				tempTiclet[t] = lotNumsArr[temp[indexCode]]
				indexCode++

				if indexCode == 9 {
					indexCode = 8
				}
			} else {
				tempTiclet[t] = "#"
			}

		}
		allZhus += indexZhu
		var tempNumLottery NumLottery = NumLottery{}
		tempNumLottery.Multiple = play.numLottery.Multiple
		tempNumLottery.BetNum = indexZhu
		tempNumLottery.Money = indexZhu * play.numLottery.Multiple * 2

		tempNumLottery.PlayType = 1

		if indexZhu == 1 {
			tempNumLottery.LotType = 1
		} else {
			tempNumLottery.LotType = 2
		}

		tempNumLottery.LotNum = strings.Join(tempTiclet, ",")
		betAll = append(betAll, tempNumLottery)

	}
	play.spliteTicket = betAll
	return allZhus
}

//设胆拆票
func (play *PlayRxn) ballDanTuoSplite() ([]NumLottery, int) {
	var betAll []NumLottery = make([]NumLottery, 0)
	var lotNumsArr []string = strings.Split(play.numLottery.LotNum, ",")
	var freeNumAll int = 0
	var danKey []int = make([]int, 0)
	var danNum map[int]int = make(map[int]int)
	var CalZhu []int = make([]int, 0)
	var oneLotMap map[int]string = make(map[int]string)
	var CalZhuAll map[int]int = make(map[int]int)
	for k, bet := range lotNumsArr {
		if bet == "#" {
			freeNumAll++
			continue
		}

		var tempPosition []string
		reg, _ := regexp.Compile(`\(.*\)`)
		flag := reg.Match([]byte(bet))
		if flag {
			rs := []rune(bet)
			var ends int = len(rs) - 1
			oneLotMap[k] = string(rs[1:ends])
			tempPosition = strings.Split(string(rs[1:ends]), " ")
			danKey = append(danKey, k)
			danNum[k] = len(tempPosition)
		} else {
			oneLotMap[k] = string(bet)
			tempPosition = strings.Split(bet, " ")
			CalZhu = append(CalZhu, k)
		}

		CalZhuAll[k] = len(tempPosition)

	}

	var tuoNum int = 9 - len(danKey)

	zhus := utils.ZuheInt{
		Nums: CalZhu,
		Nnum: len(CalZhu),
		Mnum: tuoNum,
	}

	zhuConns := zhus.ZuheResults()
	var allZhu int = 0
	for i := 0; i < len(zhuConns); i++ {
		temps := append(zhuConns[i], danKey...)
		sort.Ints(temps)
		var tempZhu int = 1
		for _, vs := range temps {
			tempZhu = tempZhu * CalZhuAll[vs]
		}

		allZhu += tempZhu
		var tempNumLottery NumLottery = NumLottery{}
		tempNumLottery.Multiple = play.numLottery.Multiple
		tempNumLottery.BetNum = tempZhu
		tempNumLottery.Money = tempZhu * play.numLottery.Multiple * 2

		tempNumLottery.PlayType = 1

		if tempZhu == 1 {
			tempNumLottery.LotType = 1
		} else {
			tempNumLottery.LotType = 2
		}

		//3,3,#,1,3,1,3,#,#,#,#,#,#,#
		//0 1 3 4 5 6 2 7 11

		var tempLotConn []string = make([]string, 14)
		var js int = 0
		for i := 0; i < 14; i++ {
			if temps[js] == i {
				tempLotConn[i] = oneLotMap[temps[js]]
				if js < 8 {
					js++
				}
			} else {
				tempLotConn[i] = "#"
			}

		}
		tempNumLottery.LotNum = strings.Join(tempLotConn, ",")
		betAll = append(betAll, tempNumLottery)
	}

	play.spliteTicket = betAll
	play.ticketNum = allZhu
	return betAll, allZhu

}

func NewPlayRxn() LotteryInterface {
	return &PlayRxn{}
}

func init() {
	Register("rxn", NewPlayRxn)
}
