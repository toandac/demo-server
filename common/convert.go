package common

import (
	"strconv"
)

func StringToInt64(str string) (int64, error) {
	n, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func RemoveDuplicateValues(strSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}

	for _, entry := range strSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
