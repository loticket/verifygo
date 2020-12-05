package utils

type ZuheInt struct {
	Nums []int //排列的数字
	Nnum int   //总个数
	Mnum int   //取出的数量

}

func (zh *ZuheInt) ZuheResults() [][]int {
	zh.Nnum = len(zh.Nums)
	return zh.findNumsByIndexs()
}

//组合算法(从nums中取出m个数)
func (zh *ZuheInt) zuheResult() [][]int {
	if zh.Mnum < 1 || zh.Mnum > zh.Nnum {
		return [][]int{}
	}

	//保存最终结果的数组，总数直接通过数学公式计算
	result := make([][]int, 0, zh.mathZuhe())
	//保存每一个组合的索引的数组，1表示选中，0表示未选中
	indexs := make([]int, zh.Nnum)
	for i := 0; i < zh.Nnum; i++ {
		if i < zh.Mnum {
			indexs[i] = 1
		} else {
			indexs[i] = 0
		}
	}

	//第一个结果
	result = zh.addTo(result, indexs)
	for {
		find := false
		//每次循环将第一次出现的 1 0 改为 0 1，同时将左侧的1移动到最左侧
		for i := 0; i < zh.Nnum-1; i++ {
			if indexs[i] == 1 && indexs[i+1] == 0 {
				find = true

				indexs[i], indexs[i+1] = 0, 1
				if i > 1 {
					zh.moveOneToLeft(indexs[:i])
				}
				result = zh.addTo(result, indexs)

				break
			}
		}

		//本次循环没有找到 1 0 ，说明已经取到了最后一种情况
		if !find {
			break
		}
	}

	return result
}

//数学方法计算组合数(从n中取m个数)
func (zh *ZuheInt) mathZuhe() int {
	return Combination(zh.Nnum, zh.Mnum)
}

//将ele复制后添加到arr中，返回新的数组
func (zh *ZuheInt) addTo(arr [][]int, ele []int) [][]int {
	newEle := make([]int, len(ele))
	copy(newEle, ele)
	arr = append(arr, newEle)

	return arr
}

func (zh *ZuheInt) moveOneToLeft(leftNums []int) {
	//计算有几个1
	sum := 0
	for i := 0; i < len(leftNums); i++ {
		if leftNums[i] == 1 {
			sum++
		}
	}

	//将前sum个改为1，之后的改为0
	for i := 0; i < len(leftNums); i++ {
		if i < sum {
			leftNums[i] = 1
		} else {
			leftNums[i] = 0
		}
	}
}

//根据索引号数组得到元素数组
func (zh *ZuheInt) findNumsByIndexs() [][]int {
	var indexs [][]int = zh.zuheResult()
	if len(indexs) == 0 {
		return [][]int{}
	}

	result := make([][]int, len(indexs))

	for i, v := range indexs {
		line := make([]int, 0)
		for j, v2 := range v {
			if v2 == 1 {
				line = append(line, zh.Nums[j])
			}
		}
		result[i] = line
	}

	return result
}
