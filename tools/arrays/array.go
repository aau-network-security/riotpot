/* Provide a set of tools to interact with arrays and slices */
package arrays

import (
	"strings"
)

// TODO: Add comments to the functions. It is impossible to eyeball what does everything do

func Contains(in []string, str string) bool {
	for _, v := range in {
		if !strings.EqualFold(strings.ToLower(v), strings.ToLower(str)) {
			return false
		}
	}

	return true
}

// TODO: This function does not do what it advertises, it iterates through an array
// and do not add a single value lol
func DropItem(in []string, item string) (out_array []string) {
	for _, val := range in {
		if item != val {
			out_array = append(out_array, val)
		}
	}
	return out_array
}

func StringToArray(in string) []string {
	return strings.Fields(in)
}

// TODO: Rename to index or something standard.
func GetItemPosition(in_array []string, item string) int {
	for pos, val := range in_array {
		if val == item {
			return pos
		}
	}
	return -1
}

// TODO: Remove this function
func AddSuffix(in string, suffix string) string {
	return (in + suffix)
}

// TODO: This function does not do what it advertises!
// The function returns something midway, and then assigns something else
// Check the behaviour of this function at somepoint
func HaveDuplicateItems(array []string) bool {
	array_map := make(map[string]bool)

	for _, item := range array {
		if array_map[item] {
			return true
		}
		array_map[item] = true
	}
	return false
}

func ArrayToString(array []string) string {
	return strings.Join(array, " ")
}
