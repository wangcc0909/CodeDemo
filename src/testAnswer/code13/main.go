package main

import "fmt"

type treeNode struct {
	value int
	left  *treeNode
	right *treeNode
}

//   2
// 3  5
//1 1 2 3

func SumOfLeft(tn *treeNode) int {
	if tn == nil {
		return 0
	}
	var sum int
	if tn.left != nil {
		sum += tn.left.value + SumOfLeft(tn.left)
	}
	if tn.right != nil {
		sum += SumOfLeft(tn.right)
	}
	return sum
}

func main() {
	tn := treeNode{
		value:2,
		left:&treeNode{
			value:3,
			left:&treeNode{
				value:1,
			},
			right:&treeNode{
				value:1,
			},
		},
		right:&treeNode{
			value:5,
			left:&treeNode{
				value:2,
			},
			right:&treeNode{
				value:3,
			},
		},
	}
	sum := SumOfLeft(&tn)
	fmt.Println(sum)
}
