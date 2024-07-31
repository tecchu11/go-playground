package env

import (
	"fmt"
	"os"
)

// ApplyString look up environment variable with given key and lookup func.
// Default lookup func is [os.LookupEnv].
// Please check given err after applying.
func ApplyString(err *error, key string, lookup func(string) (string, bool)) string {
	if lookup == nil {
		lookup = os.LookupEnv
	}
	v, ok := lookup(key)
	if !ok {
		*err = fmt.Errorf("env %s is not found", key)
		return ""
	}
	return v
}
