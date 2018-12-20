package common

func ArrContains(arr []string, value string) bool{
	for _, element := range arr {
		if (element == value) {
			return true
		}
	}
	return false
}
