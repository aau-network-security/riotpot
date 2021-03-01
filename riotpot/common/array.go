/* Provide a set of tools to interact with arrays and slices */
package array

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}