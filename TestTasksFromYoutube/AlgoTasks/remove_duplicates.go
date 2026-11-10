package AlgoTasks

func RemoveDuplicates(nums []int) int {
	startLen := len(nums)
	if startLen < 2 {
		return startLen
	}

	el := nums[0]
	l := 1
	// 1 1 2 2 3 3 4 4
	//         |
	//     |
	// 1 2
	// c 3		p 2
	for i := 1; i < startLen; i++ {
		if nums[i] != el {
			el = nums[i]
			nums[l] = nums[i]
			l++
		}

	}

	return l
}