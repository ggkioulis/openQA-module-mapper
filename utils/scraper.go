package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

func Parse() {
	var parent string
	webui_url := "https://openqa.suse.de"
	document := ParseAndGetDocument(webui_url)
	fmt.Println("Inside Parse", document)
	document.Find("a.dropdown-item").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if strings.Contains(href, "parent") {
			parent = s.Text()
		} else {
			build_url := webui_url + href
			fmt.Println(webui_url, " > Job Groups > ", parent, " > ", s.Text(), "\nlink to Builds: ", build_url, "\n---")
		}
	})
}

func ParseBuilds() {
	document := ParseAndGetDocument("https://openqa.suse.de/group_overview/110")
	fmt.Println("Inside Builds", document)
	s := document.Find("div.px-2.build-label.text-nowrap").First()
	str := s.Text()
	res := strings.Split(str, "(")
	respace := res[0]
	build_latest := strings.TrimSpace(respace)
	fmt.Println("The last build is: ", build_latest)
}

func ParseJobs() {
	document := ParseAndGetDocument("https://openqa.suse.de/tests/overview?distri=sle&version=15-SP3&build=163.1&groupid=110")
	fmt.Println("Inside Parse", document)
	document.Find("a").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
	})
}

func ParseModules() {
	document := ParseAndGetDocument("https://openqa.opensuse.org/tests/1677136")
	fmt.Println("Inside document", document)
	document.Find("td.component").Each(func(i int, s *goquery.Selection) {
		fmt.Println("ta modules einai:", s.Text())
	})
}

func ParseAndGetDocument(uri string) *goquery.Document {
	// Make HTTP request
	response, err := http.Get(uri)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Errorf("Error loading HTTP response body (%v)", err)
		return nil
	}

	return document
}
