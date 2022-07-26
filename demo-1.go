package main

import (
	"fmt"
	"regexp"
)

func main() {
	stri := "mateen (304)"
	regex, _ := regexp.Compile("\\([0-9]+\\)")
	res := regex.FindString(stri)
	res2 := regexp.MatchString(stri)
	fmt.Println(res, res2)
}
