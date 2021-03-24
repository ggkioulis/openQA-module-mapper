package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const webui_url = "https://openqa.suse.de"

func Parse() {
	var parent string
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
	s := document.Find("div.px-2.build-label.text-nowrap").First()
	s.Find("a").Each(func(k int, slc *goquery.Selection) {
		buildNumber := strings.TrimSpace(slc.Text())
		href, _ := slc.Attr("href")
		link := webui_url + href
		fmt.Println("build:", buildNumber, "/link:", link)
	})
}

func ParseJobs() {
	document := ParseAndGetDocument("https://openqa.suse.de/tests/overview?distri=sle&version=15-SP3&build=163.1&groupid=110")
	document.Find("tr").Each(func(i int, rows *goquery.Selection) {
		rows.Find("a").First().Each(func(i int, s *goquery.Selection) {
			name, status := s.Attr("data-title")
			if status == true {
				fmt.Println(name)
			}
		})
		rows.Find("td").Each(func(i int, s *goquery.Selection) {
			str, status := s.Attr("name")
			if status == true {
				res := strings.Split(str, "_")
				fmt.Println("to res einai", res[2])
			}
		})
	})
}

func ParseModules() {
	url := "https://openqa.suse.de/tests/5701825"
	autoinst_log := url + "/file/autoinst-log.txt"
	resp, err := http.Get(autoinst_log)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	}
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
		log.Fatal(err)
		return nil
	}

	return document
}
