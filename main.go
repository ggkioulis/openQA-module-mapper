package main

import (
	"github.com/ggkioulis/openQA-module-mapper/utils"
)

func main() {
	osd := utils.Webui{
		Name: "OSD",
		Url:  "https://openqa.suse.de",
	}

	osd.Scrape()
	// utils.ParseBuilds()
	//utils.ParseJobs()
	//utils.ParseModules()
}
