/* Provide a set of tools to interact with arrays and slices */
package arrays

import (
	"strings"
)

// Function to check whether an array of strings contains a string value
func Contains(in []string, str string) bool {
	for _, v := range in {
		if !strings.EqualFold(strings.ToLower(v), strings.ToLower(str)) {
			return false
		}
	}

	return true
}
