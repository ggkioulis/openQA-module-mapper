package main

import (
	"github.com/ggkioulis/openQA-module-mapper/utils"
)

func main() {
	utils.Parse()
	utils.ParseBuilds()
	utils.ParseJobs()
	utils.ParseModules()
}
