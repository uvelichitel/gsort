package gsort

func GSort[T any](data []T, less func(x, y T) bool) {
	n := len(data)
	quickSort(data, 0, n, maxDepth(n), less)
}

func maxDepth(n int) int {
	var depth int
	for i := n; i > 0; i >>= 1 {
		depth++
	}
	return depth * 2
}

func quickSort[T any](data []T, a int, b int, maxDepth int, less func(x, y T) bool) {
	for b-a > 12 { // Use ShellSort for slices <= 12 elements
		if maxDepth == 0 {
			heapSort(data, a, b, less)
			return
		}
		maxDepth--
		mlo, mhi := doPivot(data, a, b, less)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo-a < b-mhi {
			quickSort(data, a, mlo, maxDepth, less)
			a = mhi // i.e., quickSort(data, mhi, b)
		} else {
			quickSort(data, mhi, b, maxDepth, less)
			b = mlo // i.e., quickSort(data, a, mlo)
		}
	}
	if b-a > 1 {
		// Do ShellSort pass with gap 6
		// It could be written in this simplified form cause b-a <= 12
		for i := a + 6; i < b; i++ {
			if less(data[i], data[i-6]) {
				//data.Swap(i, i-6)
				data[i], data[i-6] = data[i-6], data[i]
			}
		}
		insertionSort(data, a, b, less)
	}
}

func medianOfThree[T any](data []T, m1, m0, m2 int, less func(x, y T) bool) {
	// sort 3 elements
	if less(data[m1], data[m0]) {
		data[m1], data[m0] = data[m0], data[m1]
	}
	// data[m0] <= data[m1]
	if less(data[m2], data[m1]) {
		data[m2], data[m1] = data[m1], data[m2]
		// data[m0] <= data[m2] && data[m1] < data[m2]
		if less(data[m1], data[m0]) {
			data[m1], data[m0] = data[m0], data[m1]
		}
	}
	// now data[m0] <= data[m1] <= data[m2]
}

func doPivot[T any](data []T, lo, hi int, less func(x, y T) bool) (midlo, midhi int) {
	m := int(uint(lo+hi) >> 1) // Written like this to avoid integer overflow.
	if hi-lo > 40 {
		// Tukey's ``Ninther,'' median of three medians of three.
		s := (hi - lo) / 8
		medianOfThree(data, lo, lo+s, lo+2*s, less)
		medianOfThree(data, m, m-s, m+s, less)
		medianOfThree(data, hi-1, hi-1-s, hi-1-2*s, less)
	}
	medianOfThree(data, lo, m, hi-1, less)

	// Invariants are:
	//	data[lo] = pivot (set up by ChoosePivot)
	//	data[lo < i < a] < pivot
	//	data[a <= i < b] <= pivot
	//	data[b <= i < c] unexamined
	//	data[c <= i < hi-1] > pivot
	//	data[hi-1] >= pivot
	pivot := lo
	a, c := lo+1, hi-1

	for ; a < c && less(data[a], data[pivot]); a++ {
	}
	b := a
	for {
		for ; b < c && !less(data[pivot], data[b]); b++ { // data[b] <= pivot
		}
		for ; b < c && less(data[pivot], data[c-1]); c-- { // data[c-1] > pivot
		}
		if b >= c {
			break
		}
		// data[b] > pivot; data[c-1] <= pivot
		data[b], data[c-1] = data[c-1], data[b]
		b++
		c--
	}
	// If hi-c<3 then there are duplicates (by property of median of nine).
	// Let's be a bit more conservative, and set border to 5.
	protect := hi-c < 5
	if !protect && hi-c < (hi-lo)/4 {
		// Lets test some points for equality to pivot
		dups := 0
		if !less(data[pivot], data[hi-1]) { // data[hi-1] = pivot
			data[c], data[hi-1] = data[hi-1], data[c]
			c++
			dups++
		}
		if !less(data[b-1], data[pivot]) { // data[b-1] = pivot
			b--
			dups++
		}
		// m-lo = (hi-lo)/2 > 6
		// b-lo > (hi-lo)*3/4-1 > 8
		// ==> m < b ==> data[m] <= pivot
		if !less(data[m], data[pivot]) { // data[m] = pivot
			data[m], data[b-1] = data[b-1], data[m]
			b--
			dups++
		}
		// if at least 2 points are equal to pivot, assume skewed distribution
		protect = dups > 1
	}
	if protect {
		// Protect against a lot of duplicates
		// Add invariant:
		//	data[a <= i < b] unexamined
		//	data[b <= i < c] = pivot
		for {
			for ; a < b && !less(data[b-1], data[pivot]); b-- { // data[b] == pivot
			}
			for ; a < b && less(data[a], data[pivot]); a++ { // data[a] < pivot
			}
			if a >= b {
				break
			}
			// data[a] == pivot; data[b-1] < pivot
			data[a], data[b-1] = data[b-1], data[a]
			a++
			b--
		}
	}
	// Swap pivot into middle
	data[pivot], data[b-1] = data[b-1], data[pivot]
	return b - 1, c
}

func insertionSort[T any](data []T, a int, b int, less func(x, y T) bool) {
	for i := a + 1; i < b; i++ {
		for j := i; j > a && less(data[j], data[j-1]); j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

func heapSort[T any](data []T, a int, b int, less func(x, y T) bool) {
	first := a
	lo := 0
	hi := b - a

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown(data, i, hi, first, less)
	}

	// Pop elements, largest first, into end of data.
	for i := hi - 1; i >= 0; i-- {
		//data.Swap(first, first+i)
		data[first], data[first+i] = data[first+i], data[first]
		siftDown(data, lo, i, first, less)
	}
}

func siftDown[T any](data []T, lo int, hi int, first int, less func(x, y T) bool) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && less(data[first+child], data[first+child+1]) {
			child++
		}
		if !less(data[first+root], data[first+child]) {
			return
		}
		data[first+root], data[first+child] = data[first+child], data[first+root]
		root = child
	}
}
