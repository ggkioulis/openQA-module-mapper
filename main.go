package main

import (
	"github.com/ggkioulis/openQA-module-mapper/utils"
)

func main() {
	osd := utils.Webui{
		Name: "O00",
		Url:  "https://openqa.opensuse.org",
	}

	osd.Scrape()
	// utils.ParseBuilds()
	//utils.ParseJobs()
	//utils.ParseModules()
}
