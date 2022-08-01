package gsort

import (
	"math/rand"
	"sort"
	"testing"
)

func GLess(a, b int) bool {
	return a < b
}

type ISl []int

func (i ISl) Len() int {
	return len(i)
}
func (i ISl) Less(x, y int) bool {
	return i[x] < i[y]
}
func (i ISl) Swap(x, y int) {
	i[x], i[y] = i[y], i[x]
}

var slice = rand.Perm(100000)
var islice = ISl(slice)

func BenchmarkGenericSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		GSort(slice, GLess)
	}
}

func BenchmarkInterfaceSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sort.Sort(islice)
	}
}

func BenchmarkSliceSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sort.Slice(slice, func(i, j int) bool { return slice[i] < slice[j] })
	}
}

func BenchmarkConcreteSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ISort(slice, GLess)
	}
}
