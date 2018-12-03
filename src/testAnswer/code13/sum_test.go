package main

import (
	"testing"
	"fmt"
)

func TestSumOfLeft(t *testing.T) {
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
