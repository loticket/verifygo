package utils

import (
	"strconv"
)

type StringNumArr []string

func (sn StringNumArr) Len() int { return len(sn) }

func (sn StringNumArr) Swap(i, j int) { sn[i], sn[j] = sn[j], sn[i] }

func (sn StringNumArr) Less(i, j int) bool {
	is, _ := strconv.Atoi(sn[i])
	js, _ := strconv.Atoi(sn[j])
	if is < js {
		return true
	}
	return false
}
