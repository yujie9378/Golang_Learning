package main

import (
	"fmt"
	"sync"
	"testing"
)

/*
主要是用golang
[] 给定一个二维的int slice 计算所有元素的合
 1. 按照单个slice并发计算出所有的元素的和
 2. 按照单个slice并发计算出所有的元素的和，但是最多只允许5个并发进行计算

golang的并发控制基本操作
var data [][]int
*/
var data = [][]int{
	{1, 2, 3, 4, 5, 6},
	{7, 8, 9, 10, 11, 12},
	{13, 14, 15, 16, 17, 18},
	{19, 20, 21, 22, 23, 24}}

func TestParll(t *testing.T) {

	var result int
	var lock sync.Mutex //互斥锁解决result并发的问题

	var waitFinish sync.WaitGroup //等待所有的协程完成
	waitFinish.Add(len(data))

	for index := range data {
		go func(index int) {
			defer waitFinish.Done()
			tmp := addSlice(data[index])
			lock.Lock()
			defer lock.Unlock()
			result += tmp
		}(index)
	}
	waitFinish.Wait()

	fmt.Printf("result = %d", result)

}

func TestParllWithLimit(t *testing.T) {

	var result int
	var lock sync.Mutex

	var waitFinish sync.WaitGroup

	var limit = make(chan int, 2)
	waitFinish.Add(len(data))
	for index := range data {

		limit <- 1 //等待写入，能写入的话就占个坑，开始并发
		go func(idx int) {
			defer waitFinish.Done()
			defer func() {
				<-limit // 执行完了之后就不占坑了,把机会让给别人
			}()
			tmp := addSlice(data[idx])

			lock.Lock()
			defer lock.Unlock()
			result += tmp

		}(index)
	}
	waitFinish.Wait()
	fmt.Printf("result = %d", result)
}
