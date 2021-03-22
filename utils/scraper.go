package utils

import (
	"fmt"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

func Parse() {
	document := ParseAndGetDocument("https://openqa.suse.de")
	fmt.Println("Inside Parse", document)
	document.Find("a.dropdown-item.dropdown-toggle").Each(func(i int, s *goquery.Selection) {
		fmt.Println(s.Text())
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
