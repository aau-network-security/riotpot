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

func DropItem(s []string, item string) (out_list []string) {
	for _, val := range s {
		if item !=val {
			out_list = append(out_list, val)
		}
	}
	return out_list
}
