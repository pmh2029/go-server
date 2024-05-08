package utils

func CheckSliceDiff(inputList, currentList []int) []int {
	exists := make(map[int]bool)
	var output []int
	for _, value := range currentList {
		exists[value] = true
	}

	for _, value := range inputList {
		if !exists[value] {
			output = append(output, value)
		}
	}

	return output
}

func CalcTotalPage(total int64, limit int) int {
	t := int(total) / limit
	if (int(total) % limit) != 0 {
		t += 1
	}
	return t
}
