package quicksort

import (
    "testing"
    "sort"
    "math"
    "math/rand"
)

const N = 1000000

// Test cases from "sort" standard package
var ints = []int{74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586}
var float64s = []float64{74.3, 59.0, math.Inf(1), 238.2, -784.0, 2.3, math.NaN(), math.NaN(), math.Inf(-1), 9845.768, -959.7485, 905, 7.8, 7.8}
var strings = []string{"", "Hello", "foo", "bar", "foo", "f00", "%*&^*&^&", "***", "Кирилиця", "ще одна кирилиця"}

func TestInts (t *testing.T) {
    data := sort.IntSlice(ints)
    QuickSort(data, 5)
    if !sort.IsSorted(data) {
        t.Errorf("sorted %v", ints)
        t.Errorf("   got %v", data)
    }
}

func TestFloat64s (t *testing.T) {
    data := sort.Float64Slice(float64s)
    QuickSort(data, 20)
    if !sort.IsSorted(data) {
        t.Errorf("sorted %v", float64s)
        t.Errorf("   got %v", data)
    }
}

func TestStrings (t *testing.T) {
    data := sort.StringSlice(strings)
    QuickSort(data, 1)
    if !sort.IsSorted(data) {
        t.Errorf("sorted %v", strings)
        t.Errorf("   got %v", data)
    }
}

func TestLargeInput (t *testing.T) {
    data := sort.IntSlice(rand.Perm(N)) //million
    QuickSort(data, 100)
    if !sort.IsSorted(data) {
        t.Error("Million integers")
    }
}

func TestEmpty (t *testing.T) {
    empty := make([]int, 0)
    data := sort.IntSlice(empty)
    if !sort.IsSorted(data) {
        t.Errorf("sorted %v", empty)
        t.Errorf("   got %v", data)
    }
}

func BenchmarkInt1M (b *testing.B) {
    b.StopTimer()
    for i := 0; i < b.N; i++ {
        data := sort.IntSlice(rand.Perm(N))
        b.StartTimer()
        QuickSort(data, 3)
        b.StopTimer()
    }
}

func BenchmarkInt1MSorted (b *testing.B) {
    b.StopTimer()
    for i := 0; i < b.N; i++ {
        s := make([]int, N)
        for j := 0; j < N; j++ {
            s[j] = j
        }
        data := sort.IntSlice(s)
        b.StartTimer()
        QuickSort(data, 4)
        b.StopTimer()
    }
}

func BenchmarkInt1MReverse (b *testing.B) {
    b.StopTimer()
    for i := 0; i < b.N; i++ {
        s := make([]int, N)
        for j := 0; j < N; j++ {
            s[j] = N - j
        }
        data := sort.IntSlice(s)
        b.StartTimer()
        QuickSort(data, 5)
        b.StopTimer()
    }
}

func BenchmarkInt1KDuplicates (b *testing.B) {
    b.StopTimer()
    for i := 0; i < b.N; i++ {
        s := make([]int, N)
        for j := 0; j < N; j++ {
            s[j] = rand.Intn(100)
        }
        data := sort.IntSlice(s)
        b.StartTimer()
        QuickSort(data, 6)
        b.StopTimer()
    }
}