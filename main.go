package main

import (
	"github.com/ggkioulis/openQA-module-mapper/utils"
)

var o3 = utils.Webui{
	Name: "O3",
	Url:  "https://openqa.opensuse.org",
}

var osd = utils.Webui{
	Name: "OSD",
	Url:  "https://openqa.suse.de",
}

func main() {
	osd.Scrape()
}
