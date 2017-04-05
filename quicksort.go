// 3-way concurrent quicksort. 3-way quicksort should be more efficient than classic quicksort if
// there are a lot of duplicates in sorted data.
// TODO: find out why it is slower than default sort.Sort
package quicksort

import (
    "sort" //we will use sort.Interface for our sort
    "math/rand"
    "sync"
)

var wg sync.WaitGroup

// Exported package interface
func QuickSort (data sort.Interface) {
    shuffle(data)
    wg.Add(1)
    quickSort(data, 0, data.Len() - 1)
    wg.Wait()
}

// Recursive concurrent 3-way quicksort
func quickSort (data sort.Interface, lo , hi int) {

    // Cutting off to insertion sort as it is more efficient for small arrays
    if hi - lo < 7 {
        insertionSort(data, lo, hi)
        wg.Done()
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

    wg.Add(1)
    go quickSort(data, lo, lt - 1)
    wg.Add(1)
    go quickSort(data, gt + 1, hi)
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