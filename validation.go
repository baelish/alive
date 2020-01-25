package main

var sizes = []string{"micro", "dmicro", "small", "dsmall", "medium", "dmedium", "large", "dlarge", "xlarge", "status"}

func validateBoxSize(s string) bool {
	for _, v := range sizes {
		if v == s {
			return true
		}
	}
	return false
}
