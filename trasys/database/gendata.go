package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func rands(min, max float32) float64 {
	max = max - min
	rand.Seed(time.Now().UnixNano()) //设置随机种子，使每次结果不一样
	res := Round2(float64(min+max*rand.Float32()), 2)
	return res
}

func Round2(f float64, n int) float64 {
	floatStr := fmt.Sprintf("%."+strconv.Itoa(n)+"f", f)
	inst, _ := strconv.ParseFloat(floatStr, 64)
	return inst
}

func main() {
	var num1 []int
	var num2 []int
	var num3 []float64

	et := 600000
	for ; et <= 600099; et++ {
		ans := rands(float32(10), float32(500))
		fmt.Printf("('1', '" + fmt.Sprint(et) + "', " + fmt.Sprint(ans) + "),\n")

		num1 = append(num1, 1)
		num2 = append(num2, et)
		num3 = append(num3, ans)
	}

	ct := 1
	for ; ct <= 100; ct++ {
		ans := rands(float32(10), float32(500))
		fmt.Printf("('2', '%06d', %.2f),\n", ct, ans)

		num1 = append(num1, 2)
		num2 = append(num2, ct)
		num3 = append(num3, ans)
	}

	user := 1
	l := len(num1)
	for ; user <= 100; user++ {
		//每个用户持仓1-100
		rand.Seed(time.Now().UnixNano())
		t := random(1, 10)
		for i := 0; i < t; i++ {
			//从股票中随机选择一个作为值
			idx := random(0, l-1)

			num := random(1, 100)

			sum := num3[idx] * float64(num)

			fmt.Printf("('%012d', '%d', '%06d', %.2f, %d, %.2f),\n", user, num1[idx], num2[idx], num3[idx], num, sum)
		}
	}
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}
