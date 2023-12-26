package main

var sizes = [9]string{
	"micro",
	"dmicro",
	"small",
	"dsmall",
	"medium",
	"dmedium",
	"large",
	"dlarge",
	"xlarge",
}

func validateBoxSize(s string) bool {
	for _, v := range sizes {
		if v == s {
			return true
		}
	}
	return false
}
