package util

func GetStringAtIndex(arr []string, index int) string {
	if len(arr) >= index+1 {
		return arr[index]
	}

	return ""
}
