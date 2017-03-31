// 3-way concurrent quicksort. 3-way quicksort should be more efficient than classic quicksort if
// there are a lot of duplicates in sorted data.
// TODO: find out why it is slower than default sort.Sort
package quicksort

import (
    "sort" //we will use sort.Interface for our sort
    "math/rand"
    "runtime"
)

// Exported package interface
func QuickSort (data sort.Interface) {
    shuffle(data)
    var depth int //this determines the depth of recursion where new goroutines will be created for subsorts
    for i := runtime.NumCPU(); i > 1; i >>= 1 {depth++} //take log2 of number of cores to create appropriate number of goroutines
    if depth == 0 {
        quickSort(data, 0, data.Len() - 1, depth, nil)
    } else {
        done := make(chan bool)
        go quickSort(data, 0, data.Len() - 1, depth, done)
        <-done
    }
}

// Recursive concurrent 3-way quicksort
func quickSort (data sort.Interface, lo , hi int, depth int, done chan bool) {

    // Cutting off to insertion sort as it is more efficient for small arrays
    if hi - lo < 7 {
        insertionSort(data, lo, hi)
        if done != nil {done <- true}
        return
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

    // Deciding whether to run subsorts concurrently in goroutines based on the value of depth argument
    if depth == 0 {
        quickSort(data, lo, lt - 1, 0, nil)
        quickSort(data, gt + 1, hi, 0, nil)
    } else {
        depth--
        childdone := make(chan bool)
        go quickSort(data, lo, lt - 1, depth, childdone)
        go quickSort(data, gt + 1, hi, depth, childdone)
        <-childdone
        <-childdone
    }
    if done != nil {done <- true}
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