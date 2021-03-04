/* Provide a set of tools to interact with arrays and slices */
package arrays

func Contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
