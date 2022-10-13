package main

import (
	"fmt"
	"time"
)

/*
主要是用golang
[] 给定一个二维的int slice 计算所有元素的合
 1. 按照单个slice并发计算出所有的元素的和
 2. 按照单个slice并发计算出所有的元素的和，但是最多只允许5个并发进行计算

golang的并发控制基本操作
var data [][]int
*/
func main() {
	a := [][]int{{1, 2, 3, 4, 5, 6}, {7, 8, 9, 10, 11, 12}, {13, 14, 15, 16, 17, 18}, {19, 20, 21, 22, 23, 24}}
	addslice1(a)
}

func addslice(slice [][]int) int {
	n := len(slice)
	m := len(slice[0])
	sum := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {

			go func(i, j int) {
				fmt.Println(i, j)
				sum += slice[i][j]
			}(i, j)
		}
	}
	time.Sleep(10 * time.Second)
	return sum
}
func addslice1(slice [][]int) {
	N := 5
	n := len(slice)
	m := len(slice[0])
	a := n * m
	k := a / N
	fmt.Println("k", k)
	y := a % N
	//sum := 0
	var ch [5]chan int
	for i := 0; i < N; i++ {
		ch[i] = make(chan int, k+y)
		for j := 0; j < k; j++ {
			r := (i*k + j) / m
			c := (i*k + j) % m
			fmt.Println("slice", i, slice[r][c])
			ch[i] <- slice[r][c]
		}
		if i == N-1 {
			for j := 0; j < y; j++ {
				r := (k*N + j) / m
				c := (k*N + j) % m
				fmt.Println("slice", slice[r][c])
				ch[N-1] <- slice[r][c]
			}
		}

	}
	sum := 0
	for i := 0; i < N; i++ {
		//fmt.Println( <-ch[i])
		worker(ch[i])

	}
	fmt.Println(sum)

}
func worker(ch chan int) int {
	sum := 0
	for i := range ch {
		fmt.Println(i)
		sum += i
	}
	return sum
}
