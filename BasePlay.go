package verifygo

import (
	"sort"
	"strconv"
	"strings"
)

//彩种投注号码定义
type ball struct {
	min, max int             //号码的个数
	names    map[string]uint //投注的内容
}

//数字彩投注内容
type NumLottery struct {
	PlayType int    `json:"playtype" binding:"required"` //子玩法
	LotType  int    `json:"lottype" binding:"required"`  //选号方式
	LotNum   string `json:"lotnum" binding:"required"`   //投注号码
	Money    int    `json:"money" binding:"required"`    //投注金额
	BetNum   int    `json:"betnum" binding:"required"`   //注数
	Multiple int    `json:"multiple" binding:"required"` //倍数
}

type BasePlay struct {
	ticketNum  int        //计算出来的注数
	numLottery NumLottery //一张票
}

//设置投注票信息
func (bp *BasePlay) CreateTicket(playtype int, lottype int, lotnum string, money int, betnum int, multiple int) {
	bp.numLottery = NumLottery{
		PlayType: playtype,
		LotType:  lottype,
		LotNum:   lotnum,
		Money:    money,
		BetNum:   betnum,
		Multiple: multiple,
	}
}

//获取计算出来的注数
func (bp *BasePlay) GetTicketNum() int {
	return bp.ticketNum
}

//检查红球是否在规定的里面
//[01 04 08 07 15 28 30]
func (bp *BasePlay) Ball(balls []string, ballReal map[string]uint) bool {
	var flag bool = true
	var ballsNum int = len(balls)
	var ballBak []string = make([]string, ballsNum)
	copy(ballBak, balls)
	sort.Strings(balls)
	//验证号码是否正确
	for k, v := range balls {
		_, ok := ballReal[v]
		if !ok || !strings.EqualFold(v, ballBak[k]) {
			flag = false
			break
		}
		if k > 0 {
			if strings.EqualFold(balls[k], balls[k-1]) {
				flag = false
				break
			}
		}

	}
	return flag
}

//重新排序各种球
func (bp *BasePlay) BallAsc(balls []string, ballReal map[string]uint) bool {
	var flag bool = true
	var ballsNum int = len(balls)
	var ballBak []string = make([]string, ballsNum)
	copy(ballBak, balls)
	sort.Sort(sort.Reverse(sort.StringSlice(balls)))
	//验证号码是否正确
	for k, v := range balls {
		_, ok := ballReal[v]
		if !ok || !strings.EqualFold(v, ballBak[k]) {
			flag = false
			break
		}
		if k > 0 {
			if strings.EqualFold(balls[k], balls[k-1]) {
				flag = false
				break
			}
		}

	}
	return flag
}

//检查红球是否在规定的里面
func (bp *BasePlay) BallInt(balls []string, ballReal map[string]uint) bool {
	var flag bool = true
	var ballsNum int = len(balls)
	var ballIntNum []int = make([]int, ballsNum)
	for i := 0; i < ballsNum; i++ {
		tint, err := strconv.Atoi(balls[i])
		if err != nil {
			return false
			break
		}
		ballIntNum[i] = tint
	}
	var ballBak []int = make([]int, ballsNum)
	copy(ballBak, ballIntNum)
	sort.Ints(ballIntNum)
	//验证号码是否正确
	for k, v := range balls {
		_, ok := ballReal[v]
		if !ok || ballIntNum[k] != ballBak[k] {
			flag = false
			break
		}
		if k > 0 {
			if ballIntNum[k] == ballIntNum[k-1] {
				flag = false
				break
			}
		}

	}
	return flag
}

//检查投注数字是否在规定的里面
func (bp *BasePlay) BallSame(balls []string, ballReal map[string]uint) bool {
	var flag bool = true
	var ballsNum int = len(balls)
	var ballBak []string = make([]string, ballsNum)
	copy(ballBak, balls)
	//验证号码是否正确
	for k, v := range balls {
		_, ok := ballReal[v]
		if !ok || !strings.EqualFold(v, ballBak[k]) {
			flag = false
			break
		}

	}
	return flag
}

func (bp *BasePlay) BallSameSort(balls []string, ballReal map[string]uint) bool {
	var flag bool = true
	var ballsNum int = len(balls)
	var ballBak []string = make([]string, ballsNum)
	copy(ballBak, balls)
	sort.Strings(balls)
	//验证号码是否正确
	for k, v := range balls {
		_, ok := ballReal[v]
		if !ok || !strings.EqualFold(v, ballBak[k]) {
			flag = false
			break
		}

	}
	return flag
}

//检查红球是否在规定的里面--整数
func (bp *BasePlay) BallSameInt(balls []string, ballReal map[string]uint) bool {
	var flag bool = true
	var ballsNum int = len(balls)
	var ballIntNum []int = make([]int, ballsNum)
	for i := 0; i < ballsNum; i++ {
		tint, err := strconv.Atoi(balls[i])
		if err != nil {
			return false
			break
		}
		ballIntNum[i] = tint
	}
	var ballBak []int = make([]int, ballsNum)
	copy(ballBak, ballIntNum)
	sort.Ints(ballIntNum)
	//验证号码是否正确
	for k, v := range balls {
		_, ok := ballReal[v]
		if !ok || ballIntNum[k] != ballBak[k] {
			flag = false
			break
		}

	}
	return flag
}
