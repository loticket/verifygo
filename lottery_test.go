package verifygo

import (
	"fmt"
	"log"
	"testing"
)

func Test_main(t *testing.T) {
	fmt.Println("--------开始调试------")
	lottery, err := NewLottery("sfc")
	t.Log(err)
	if err != nil {
		log.Println("不存在")
	}

	t.Error(err)

	t.Log("-------------")

	lottery.CreateTicket(1, 2, "3,3 1 0,3,3,3,3,3,3,3,3,3,3,3,3", 6, 3, 1)

	_, err = lottery.Verification()
	fmt.Println(lottery.GetSpliteTicket())
	t.Error(err)

}
