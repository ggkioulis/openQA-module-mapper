package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

func Parse() {
	var parent = ""
	var currentparent = &parent
	document := ParseAndGetDocument("https://openqa.suse.de")
	fmt.Println("Inside Parse", document)
	document.Find("a.dropdown-item").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if strings.Contains(href, "parent") {
			*currentparent = s.Text()
			fmt.Println("=================\n", s.Text(), "\n --------------")
		} else {
			fmt.Println(s.Text(), "with current parent: ", parent, "\n --------------")
		}
	})
}

func ParseBuilds() {
	document := ParseAndGetDocument("https://openqa.suse.de/group_overview/110")
	fmt.Println("Inside Builds", document)
	s := document.Find("div.px-2.build-label.text-nowrap").First()
	//href, _ := s.Attr("href")
	str1 := s.Text()
	strTrim := strings.TrimSpace(str1)
	fmt.Println("==>" + strTrim + "<==")
	res1 := strings.Split(strTrim, "(")
	fmt.Println("To split einai: ", res1[0])
	//fmt.Println("value: ", href, "to text einai: ", s.Text())
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
