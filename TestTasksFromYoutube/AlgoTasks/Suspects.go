package AlgoTasks

import "fmt"

func FilterInnocents() {
	suspects := []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	//suspects := []int{1, 2, 3}
	innocents := []int{1, 2, 3, 5}
	result := filter2(suspects, innocents)
	fmt.Println(result)
}

func filter(suspects, innocents []int) []int {
	result := make([]int, 0, len(suspects))
	si, ii := 0, 0
	for si < len(suspects) && ii < len(innocents) {
		if innocents[ii] > suspects[si] {
			result = append(result, suspects[si])
			si++
		} else {
			if innocents[ii] == suspects[si] {
				si++
			}
			ii++
		}
	}
	if si < len(suspects) {
		result = append(result, suspects[si:]...)
	}
	return result
}

func filter2(suspects, innocents []int) []int {
	kv := make(map[int]struct{}, len(innocents))
	result := make([]int, 0, len(suspects))

	for _, v := range innocents {
		kv[v] = struct{}{}
	}

	for _, s := range suspects {
		if _, ok := kv[s]; !ok {
			result = append(result, s)
		}
	}

	return result
}
