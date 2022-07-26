package main

import (
	"fmt"
	//"regexp"
)

func main() {
	var x = (Food{"rice", 2} == Food{"rice", 3})
	fmt.Println(x)
}

type Food struct {
	Name string
	Rating int
}
