package main

import (
	"fmt"
	"sync"
	"testing"
)

var wg sync.WaitGroup
var lock sync.Mutex

func TestDay01_5(t *testing.T) {
	a := [][]int{
		{1, 2, 3, 4, 5, 6},
		{7, 8, 9, 10, 11, 12},
		{13, 14, 15, 16, 17, 18},
		{19, 20, 21, 22, 23, 24}}
	sum := 0
	wg.Add(len(a))
	for i := 0; i < len(a); i++ {
		go func(i int) {
			lock.Lock()
			sum += addSlice1(a[i])
			defer lock.Unlock()
			defer wg.Done()
		}(i)

	}
	wg.Wait()
	fmt.Println(sum)
	//addSlice(a)
}
func TestDay01_2(t *testing.T) {
	a := [][]int{
		{1, 2, 3, 4, 5, 6},
		{7, 8, 9, 10, 11, 12},
		{13, 14, 15, 16, 17, 18},
		{19, 20, 21, 22, 23, 24},
		{1, 2, 3, 4, 5, 6},
		{7, 8, 9, 10, 11, 12},
		{13, 14, 15, 16, 17, 18},
		{19, 20, 21, 22, 23, 24},
		{1, 2, 3, 4, 5, 6},
		{7, 8, 9, 10, 11, 12},
		{13, 14, 15, 16, 17, 18},
		{19, 20, 21, 22, 23, 24},
	}
	sum := 0
	n := len(a) / 5
	if len(a)%5 != 0 {
		n += 1
	}
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(i int) {
			array := make([]int, 0)
			for j := 0; j < n && (i*n+j) < len(a); j++ {

				array = append(array, a[i*n+j]...)
			}
			lock.Lock()
			sum += addSlice1(array)
			defer lock.Unlock()
			defer wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println(sum)
}

func addSlice1(a []int) int {

	var result int
	for _, v := range a {

		result = result + v
	}
	return result

}
func TestDay01_3(t *testing.T) {
	a := [][]int{
		{1, 2, 3, 4, 5, 6},
		{7, 8, 9, 10, 11, 12},
		{13, 14, 15, 16, 17, 18},
		{19, 20, 21, 22, 23, 24},
		{1, 2, 3, 4, 5, 6},
		{7, 8, 9, 10, 11, 12},
		{13, 14, 15, 16, 17, 18},
		{19, 20, 21, 22, 23, 24},
		{1, 2, 3, 4, 5, 6},
		{7, 8, 9, 10, 11, 12},
		{13, 14, 15, 16, 17, 18},
		{19, 20, 21, 22, 23, 24},
		{1, 2, 3, 4, 5, 6},
		{7, 8, 9, 10, 11, 12},
		{13, 14, 15, 16, 17, 18},
		{19, 20, 21, 22, 23, 24},
		{1, 2, 3, 4, 5, 6},
		{7, 8, 9, 10, 11, 12},
		{13, 14, 15, 16, 17, 18},
		{19, 20, 21, 22, 23, 24},
	}
	sum := 0
	n := len(a) / 5
	wg.Add(5)
	for i := 0; i < 5; i++ {
		go func(i int) {
			array := make([]int, 0)
			for j := 0; j < n; j++ {
				array = append(array, a[i*n+j]...)
			}
			sum += addSlice1(array)
			if i == 4 && len(a)%5 != 0 {
				k := i % 5
				for j := 0; j < k; j++ {

					array = append(array, a[n*5+j]...)
				}
				lock.Lock()
				sum += addSlice1(array)
				defer lock.Unlock()
			}
			defer wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println(sum)
}
func TestDay01_4(t *testing.T) {
	a := [][]int{
		{1, 2, 3, 4, 5, 6},
		{7, 8, 9, 10, 11, 12},
		{13, 14, 15, 16, 17, 18},
		{19, 20, 21, 22, 23, 24},
		{1, 2, 3, 4, 5, 6},
		{7, 8, 9, 10, 11, 12},
		{13, 14, 15, 16, 17, 18},
		{19, 20, 21, 22, 23, 24},
		{1, 2, 3, 4, 5, 6},
		{7, 8, 9, 10, 11, 12},
		{13, 14, 15, 16, 17, 18},
		{19, 20, 21, 22, 23, 24},
	}
	sum := 0
	limit := make(chan int, 5)
	wg.Add(len(a))
	for index := range a {
		limit <- 1
		go func(idx int) {
			lock.Lock()
			sum += addSlice(a[idx])
			defer lock.Unlock()
			defer wg.Done()
			<-limit
		}(index)
	}
	wg.Wait()
	fmt.Println("sum", sum)

}
func addSlice(a []int) int {
	var result int
	for _, v := range a {
		result = result + v
	}
	return result
}
