package verifygo

import (
	"fmt"
	"log"
	"testing"
)

func Test_main(t *testing.T) {
	fmt.Println("--------开始调试------")
	lottery, err := NewLottery("qlc")
	t.Log(err)
	if err != nil {
		log.Println("不存在")
	}

	t.Error(err)

	t.Log("-------------")

	lottery.CreateTicket(1, 3, "09,14,24|04,08,16,23,26", 10, 5, 1)

	_, err = lottery.Verification()
	t.Error(err)

}
