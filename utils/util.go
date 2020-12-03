package utils

import (
	"sort"
)

//组合计算公式
func Combination(n int, r int) int {
	return int(Factorial(int64(n)) / (Factorial(int64(n-r)) * Factorial(int64(r))))
}

//阶乘计算
func Factorial(n int64) int64 {
	if n == 1 {
		return 1
	} else if n == 2 {
		return 2
	}

	var result int64 = 1
	var i int64
	for i = 2; i <= n; i++ {
		result *= i
	}
	return result
}

//二分查找算法
func binarySearch(sortedList []int, lookingFor int) int {
	var lo int = 0
	var hi int = len(sortedList) - 1
	sort.Ints(sortedList)
	for lo <= hi {
		var mid int = lo + (hi-lo)/2
		var midValue int = sortedList[mid]
		if midValue == lookingFor {
			return midValue
		} else if midValue > lookingFor {
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	return -1
}
