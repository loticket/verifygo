package verifygo

import (
	"fmt"
	"log"
	"testing"
)

func Test_main(t *testing.T) {
	fmt.Println("--------开始调试------")
	lot, err := NewLottery("qxc")
	t.Log(err)
	if err != nil {
		log.Println("不存在")
	}

	lot.CheckPlaytype()

	lotterys := NewPlayStrategy().Lottery(lot)

	lotterys.SetTicket(1, 2, "4;3;3,5;3,7;4;6-4,6,14", 2400, 12, 100)

	yd, yderr := lotterys.Verification()

	t.Error(yd)
	t.Error(yderr)

	gg := lotterys.GetSpliteTicket()

	t.Log(gg)
	//	t.Log("-------------")

	t.Log("-------------")

	//lottery.

	//lottery.SetTicket(1, 2, "3,3 1 0,3,3,3,3,3,3,3,3,3,3,3,3", 6, 3, 1)

	//_, err = lottery.Verification()
	//fmt.Println(lottery.GetSpliteTicket())
	t.Error(err)

}
