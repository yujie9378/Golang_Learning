package main

import (
	"fmt"
	"sync"
	"testing"
)

var wg sync.WaitGroup

func TestDay01(t *testing.T) {
	a := [][]int{
		{1, 2, 3, 4, 5, 6},
		{7, 8, 9, 10, 11, 12},
		{13, 14, 15, 16, 17, 18},
		{19, 20, 21, 22, 23, 24}}
	sum := 0
	wg.Add(len(a))
	for i := 0; i < len(a); i++ {
		go func(i int) {

			sum += addSlice1(a[i])
			wg.Done()
		}(i)

	}
	wg.Wait()
	fmt.Println(sum)
	//addSlice(a)
}
func TestDay02(t *testing.T) {
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
			sum += addSlice1(array)
			wg.Done()
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
func TestDay03(t *testing.T) {
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
				sum += addSlice1(array)
			}
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println(sum)
}
