package utils

func ArrayContains(array []int, target int) bool {
	for i := range array {
		if i == target { return true }
	}
	return false
}