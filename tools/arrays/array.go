/* Provide a set of tools to interact with arrays and slices */
package arrays

import (
	"strings"
)
func Contains(in []string, str string) bool {
	for _, v := range in {
		if v == str {
			return true
		}
	}

	return false
}

func DropItem(in []string, item string) (out_array []string) {
	for _, val := range in {
		if item !=val {
			out_array = append(out_array, val)
		}
	}
	return out_array
}

func StringToArray(in string) []string {
	return strings.Fields(in)
}

func GetItemPosition(in_array []string, item string) (int) {
	for pos, val := range in_array {
		if val == item {
			return pos
		}
	}
	return -1
}

func AddSuffix(in string, suffix string) string {
	return (in+suffix)
}

func HasDuplicateItems(array []string) bool {
	array_map := make(map[string]bool)

	for _, item := range array {
		if array_map[item] == true {
			return true
		}	
		array_map[item] = true
	}
	return false
}

