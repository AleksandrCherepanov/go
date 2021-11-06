package utils

import "errors"

func MapInvert(inputMap map[string]int) (map[int]string, error) {
	result := make(map[int]string, len(inputMap))

	for key, value := range inputMap {
		if _, ok := result[value]; ok {
			return nil, errors.New("duplicate key exists")
		}
		result[value] = key
	}

	return result, nil
}
