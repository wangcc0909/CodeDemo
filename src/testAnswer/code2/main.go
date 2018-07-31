package main

import (
	"math/rand"
	"time"
	"fmt"
)

func GetAwardUserName(users map[string]int64) (name string) {
	sizeOfUsers := len(users)
	randSize := rand.Intn(sizeOfUsers)

	var index int
	for uName := range users {
		if index == randSize {
			name = uName
			return
		}
		index++
	}
	return
}

func GetAwardGenerator(users map[string]int64) (generator func() string) {
	var sum_num int64
	name_arr := make([]string, len(users))
	offset_arr := make([]int64, len(users))
	var index int
	for u_name, num := range users {
		name_arr[index] = u_name
		offset_arr[index] = sum_num
		sum_num += num
		index += 1
	}

	generator = func() string {
		award_num := rand.Int63n(sum_num)
		return name_arr[binary_search(offset_arr, award_num)]
	}

	return
}

func binary_search(nums []int64, target int64) int {
	start, end := 0, len(nums)-1

	for start <= end {
		mid := (start + end) / 2
		if nums[mid] > target {
			end = mid - 1
		} else if nums[mid] < target {
			if mid+1 == len(nums) {
				return mid
			}
			if nums[mid+1] > target {
				return mid
			}
			start = mid + 1
		} else {
			return mid
		}
	}
	return -1
}

func main() {
	var users = map[string]int64{
		"a": 10,
		"b": 6,
		"c": 3,
		"d": 2,
		"e": 1,
	}

	rand.Seed(time.Now().Unix())
	var awards = make(map[string]int)
	generator := GetAwardGenerator(users)
	for i := 0; i < 100000; i++ {
		name := generator()
		if count, ok := awards[name]; ok {
			awards[name] = count + 1
		} else {
			awards[name] = 1
		}
	}

	for name, count := range awards {
		fmt.Printf("name :%s,count : %d\n", name, count)
	}
}
