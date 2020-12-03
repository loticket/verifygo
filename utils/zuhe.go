package utils

type Zuhe struct {
	Nums []interface{} //组合的数组
	In   int           //组合数量
	Out  int           //取出的数量
}

//先计算注数

func (z *Zuhe) MathZhu() int {
	return z.Factorial(z.In) / (z.Factorial(z.In-z.Out) * z.Factorial(z.Out))
}

//计算组合
func (z *Zuhe) ZuheResult() [][]interface{} {
	if z.In < 1 || z.Out > z.In {
		return [][]interface{}{}
	}

	//保存最终结果的数组，总数直接通过数学公式计算
	result := make([][]interface{}, 0, z.MathZhu())

	//保存每一个组合的索引的数组，1表示选中，0表示未选中
	indexs := make([]interface{}, z.In)
	for i := 0; i < z.In; i++ {
		if i < z.Out {
			indexs[i] = 1
		} else {
			indexs[i] = 0
		}
	}

	//第一个结果
	result = z.addTo(result, indexs)

	for {
		find := false
		//每次循环将第一次出现的 1 0 改为 0 1，同时将左侧的1移动到最左侧
		for i := 0; i < z.In-1; i++ {
			if indexs[i] == 1 && indexs[i+1] == 0 {
				find = true

				indexs[i], indexs[i+1] = 0, 1
				if i > 1 {
					z.moveOneToLeft(indexs[:i])
				}
				result = z.addTo(result, indexs)

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

//将ele复制后添加到arr中，返回新的数组
func (z *Zuhe) addTo(arr [][]interface{}, ele []interface{}) [][]interface{} {
	newEle := make([]interface{}, len(ele))
	copy(newEle, ele)
	arr = append(arr, newEle)
	return arr
}

//移动一位1
func (z *Zuhe) moveOneToLeft(leftNums []interface{}) {
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
func (z *Zuhe) FindNumsByIndexs() [][]interface{} {
	indexs := z.ZuheResult()
	if len(indexs) == 0 {
		return [][]interface{}{}
	}

	result := make([][]interface{}, len(indexs))

	for i, v := range indexs {
		line := make([]interface{}, 0)
		for j, v2 := range v {
			if v2 == 1 {
				line = append(line, z.Nums[j])
			}
		}
		result[i] = line
	}

	return result
}

//阶乘
func (z *Zuhe) Factorial(n int) int {
	result := 1
	for i := 2; i <= n; i++ {
		result *= i
	}

	return result
}

//实例化组合
func NewZuhe(nums []interface{}, out int) Zuhe {
	var in int = len(nums)
	return Zuhe{
		Nums: nums,
		In:   in,
		Out:  out,
	}
}
