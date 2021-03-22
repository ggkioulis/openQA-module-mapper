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
