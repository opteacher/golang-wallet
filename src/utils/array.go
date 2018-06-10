package utils

func Contains(array []int, target int) bool {
	for i := range array {
		if i == target { return true }
	}
	return false
}