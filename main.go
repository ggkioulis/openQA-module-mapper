package main

import (
	"fmt"

	"github.com/ggkioulis/openQA-module-mapper/utils"
)

func main() {
	defer fmt.Println("za")
	fmt.Println("Inside main")
	utils.Parse()
}
