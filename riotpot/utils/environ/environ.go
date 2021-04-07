/*
Package environ provides functions used to interact with the environment
*/
package environ

import "os"

/*
	Get the value of the variable set in the Environment if it exists,
	otherwise it returns the fallback value
	@key: the string name of the variable set in the environment
	@fallback: the default value in case @key does not exists
*/
func Getenv(key string, fallback string) string {

	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
