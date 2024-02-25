package util

import "fmt"

func GetStringAtIndex(arr []string, index int) string {
	if len(arr) >= index+1 {
		return arr[index]
	}

	return ""
}

func IsStringEmpty(items ...string) error {
	for _, i := range items {
		if i == "" {
			return fmt.Errorf("%s is empty when it shouldn't be", i)
		}
	}

	return nil
}
