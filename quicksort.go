// 3-way concurrent quicksort. 3-way quicksort should be more efficient than classic quicksort if
// there are a lot of duplicates in sorted data.
// TODO: find out why it is slower than default sort.Sort
package quicksort

import (
    "sort" //we will use sort.Interface for our sort
    "math/rand"
    "sync"
    //"fmt"
)

type Srt struct {
    Data sort.Interface
    Start int
    End int
}

var wg sync.WaitGroup
var doneSoFar int
var closed bool

// Exported package interface
func QuickSort (data sort.Interface, goroutines int) {
    shuffle(data)
    sorts := make(chan Srt, 100)
    for i := 1; i <= goroutines; i++ {
        wg.Add(1)
        go quickSort(sorts)
    }
    sorts <- Srt{data, 0, data.Len() - 1}
    wg.Wait()
}

// Recursive concurrent 3-way quicksort
func quickSort (sorts chan Srt) {

    for s := range sorts {

        data := s.Data
        lo := s.Start
        hi := s.End

        // Cutting off to insertion sort as it is more efficient for small arrays
        if hi - lo < 7 {
            insertionSort(data, lo, hi)
            doneSoFar += hi - lo + 1
            //fmt.Println("Done so far: ", doneSoFar)
            if doneSoFar == data.Len() && !closed {
                close(sorts)
                closed = true
            }
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

        doneSoFar += gt - lt + 1
        //fmt.Println("Done so far: ", doneSoFar)
        if doneSoFar == data.Len() && !closed {
            close(sorts)
            closed = true
            continue
        }

        go func () {
            if !closed {
                sorts <- Srt{data, lo, lt - 1}
                sorts <- Srt{data, gt + 1, hi}
            }
        }()
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