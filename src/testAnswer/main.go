package main

import "fmt"

type S struct {
	v int
}

//一个按升序排列好的数组int[] arry = {-5,-1,0,5,9,11,13,15,22,35,46},输入一个x，int x = 31，
// 在数据中找出和为x的两个数，例如 9 + 22 = 31，要求算法的时间复杂度为O(n);

func main() {
	/*s := []S{{1}, {3}, {5}, {2}}
	// A
	sort.Slice(s, func(i, j int) bool {
		return s[i].v < s[j].v
	})
	fmt.Printf("%#v", s)

	fmt.Println(utf8.RuneCountInString("你好"))

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w,"hello")
	})
	// B
	http.ListenAndServe(":8000",mux)*/
	z := 31
	x,y := add(z)

	for i := range x{
		fmt.Printf("%d = %d + %d\n",z,x[i],y[i])
	}

}

var array = []int{-5,-1,0,5,9,11,13,15,20,22,35,46}

func add(result int) ([]int, []int) {
	var x,y []int

	for i := 0; i < len(array); i++ {

		for j := len(array) - 1; j > i; j-- {
			r := array[i] + array[j]
			if r < result {
				break
			}

			if r == result {
				x = append(x,array[i])
				y = append(y,array[j])
			}
		}

	}

	return x,y
}


