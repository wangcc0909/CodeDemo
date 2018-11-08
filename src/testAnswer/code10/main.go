package main

import "fmt"

func main() {
	//给定一个数组代表股票每天的价格，请问只能买卖一次的情况下，最大化利润是多少？
	// 日期不重叠的情况下，可以买卖多次呢？输入：{100,80,120,130,70,60,100,125}，
	// 只能买一次：65(60买进，125卖出)；可以买卖多次：115(80买进，130卖出；60买进，125卖出)？
	a := []int{100, 80, 120, 130, 70, 60, 100, 125}
	type Op struct {
		buyDay,saleDay int
		buyPrice,salePrice int
		earning int
	}

	var Ops []Op
	var buyPrice,salePrice = 1 << 32 -1,-1 << 32 -1
	var buyDay,saleDay = -1,-1

	for index,todayPrice := range a {
		//寻找买入点
		if buyDay == -1 {
			if todayPrice < buyPrice {
				buyPrice = todayPrice
				continue
			}
			//买入
			buyDay = index -1
			continue
		}
		//寻找卖出点
		if todayPrice > salePrice {
			 salePrice = todayPrice
			if index < len(a) - 1 {
				 continue
			}
		}
		//卖出
		if index < len(a) - 1 {
			saleDay = index - 1
		}else {
			saleDay = index
		}

		var op = Op{
			buyDay:buyDay,
			saleDay:saleDay,
			buyPrice:buyPrice,
			salePrice:salePrice,
			earning:salePrice - buyPrice,
		}
		Ops = append(Ops,op)
		//进行下一轮
		buyPrice,salePrice = 1 << 32 -1,-1 << 32 -1
		buyDay,saleDay = -1,-1
	}
	fmt.Printf("%v",Ops)
}
