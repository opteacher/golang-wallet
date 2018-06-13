package utils

func IntArrayContains(array []int, target int) bool {
	for _, i := range array {
		if i == target { return true }
	}
	return false
}

func StrArrayContains(array []string, target string) bool {
	for _, i := range array {
		if i == target { return true }
	}
	return false
}