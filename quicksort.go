// 3-way concurrent quicksort. 3-way quicksort should be more efficient than classic quicksort if
// there are a lot of duplicates in sorted data.
// TODO: find out why it is slower than default sort.Sort
package quicksort

import (
    "sort" //we will use sort.Interface for our sort
    "math/rand"
    "sync"
   // "fmt"
)

// This is the abstraction for sorting job. Since each sort is independent of others, I've
// decided to try using goroutines and channels instead of recursion.
type Srt struct {
    Data sort.Interface
    Start int
    End int
}

var wg sync.WaitGroup

// Exported package interface
func QuickSort (data sort.Interface, goroutines int) {
    shuffle(data)
    sorts := make(chan Srt, 100)
    done := make(chan int, 100)
    go sortCloser(sorts, done, data.Len())
    for i := 1; i <= goroutines; i++ {
        wg.Add(1)
        go quickSort(sorts, done)
    }
    sorts <- Srt{data, 0, data.Len() - 1}
    wg.Wait()
}

// Recursive concurrent 3-way quicksort
// Goroutine takes sorting job from sorts channel, sorts it if the size is less then 12 or
// partitions and puts two subsorts into sorts channel.
func quickSort (sorts chan Srt, done chan int) {

    for s := range sorts {

        data := s.Data
        lo := s.Start
        hi := s.End

        // Cutting off to insertion sort as it is more efficient for small arrays
        if hi - lo < 12 {
            insertionSort(data, lo, hi)
            done <- hi - lo + 1
            continue
        }

        // 3-way quicksort logic
        lt, gt := lo, hi
        i := lo
        for i <= gt {
            if data.Less(i, lt) {
                data.Swap(lt, i)
                i++
                lt++
            } else if data.Less(lt, i) {
                data.Swap(i, gt)
                gt--
            } else {
                i++
            }
        }

        done <- gt - lt + 1

        if lo <= lt - 1 {go func () {sorts <- Srt{data, lo, lt - 1}}()}
        if gt + 1 <= hi {go func () {sorts <- Srt{data, gt + 1, hi}}()}
    }

    wg.Done()
}

// For small subarrays
func insertionSort (data sort.Interface, lo, hi int) {
    for i := lo; i <= hi; i++ {
        for j := i; j > lo; j-- {
            if data.Less(j, j-1) {
                data.Swap(j, j-1)
            } else {
                break
            }
        }
    }
}

// Shuffle data randomly before quicksort
func shuffle (data sort.Interface) {
    N := data.Len()
    for i := 0; i < N; i++ {
        r := rand.Intn(i + 1)
        data.Swap(r, i)
    }
}

// This goroutine counts number of elements that are already sorted and
// when it is equal to the length of entire slice, closes the sorts channel.
func sortCloser (sorts chan Srt, done chan int, length int) {
    var sum int
    for s := range done {
        sum += s
        if sum == length {
            close(sorts)
            close(done)
        }
    }
}