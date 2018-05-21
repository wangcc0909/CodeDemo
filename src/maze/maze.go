package main

import (
	"os"
	"fmt"
	"log"
)

func readMaze(path string) [][]int {
	file,err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var row,col int
	fmt.Fscanf(file,"%d %d",&row,&col)
	fmt.Fscanln(file)//window中读取一行不会换行  加上这个换行

	maze := make([][]int,row)

	for i := range maze{
		maze[i] = make([]int,col)
		for j := range maze[i]{
			fmt.Fscanf(file,"%d",&maze[i][j])
		}
		fmt.Fscanln(file)
	}

	return maze
}

//点的坐标值
type point struct {
	i, j int
}

var dirs = [4]point{
	{-1,0},{0,-1},{0,1},{1,0},
}

func (p point) add(r point) point {
	return point{p.i + r.i,p.j + r.j}
}

func (p point) at(ar [][]int) (int,bool) {
	if p.i < 0 || p.i >= len(ar) {
		return 0, false
	}

	if p.j < 0 || p.j >= len(ar[p.i]) {
		return 0, false
	}

	return ar[p.i][p.j],true
}

func walk(maze [][]int,start,end point) [][]int {

	//创建slience
	steps := make([][]int,len(maze))

	for i := range steps{
		steps[i] = make([]int,len(maze[i]))
	}

	//需要将所有的点存到一个队列
	Q := []point{start}

	for len(Q) > 0 {
		cur := Q[0]
		Q = Q[1:]

		if cur == end {
			break
		}

		for _,dir := range dirs{
			next := cur.add(dir)

			//这里处理的是撞墙的
			val,ok := next.at(maze)
			if !ok || val == 1 {
				continue
			}

			//这里处理越界
			val,ok = next.at(steps)
			if !ok || val != 0 {
				continue
			}

			//和上一个点相同的情况
			if next == start {
				continue
			}

			curSteps,_ := cur.at(steps)
			steps[next.i][next.j] = curSteps + 1

			Q = append(Q,next)
		}

	}
	return steps

}

func main() {
	maze := readMaze("src/maze/maze.in")
	log.Print(len(maze),len(maze[0]))

	steps := walk(maze,point{0,0},point{len(maze),len(maze[0])})

	for _,row := range steps{
		for _,col := range row{
			fmt.Printf("%3d",col)
		}
		fmt.Println()
	}
}
