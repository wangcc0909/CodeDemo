package main

import (
	"fmt"
)

type treeNode struct {
	value int
	left  *treeNode
	right *treeNode
}

//        10
//   8         12
// 6   9     7   13
//1 7 5  11 4 9 11  14

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
		value:10,
		left:&treeNode{
			value:8,
			left:&treeNode{
				value:6,
				left:&treeNode{
					value:1,
				},
				right:&treeNode{
					value:7,
				},
			},
			right:&treeNode{
				value:9,
				left:&treeNode{
					value:5,
				},
				right:&treeNode{
					value:11,
				},
			},
		},
		right:&treeNode{
			value:12,
			left:&treeNode{
				value:7,
				left:&treeNode{
					value:4,
				},
				right:&treeNode{
					value:9,
				},
			},
			right:&treeNode{
				value:13,
				left:&treeNode{
					value:11,
				},
				right:&treeNode{
					value:14,
				},
			},
		},
	}
	sum := SumOfLeft(&tn)
	fmt.Println(sum)
}
