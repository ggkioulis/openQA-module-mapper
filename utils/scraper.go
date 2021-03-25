package utils

import (
	"bufio"
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
			// JOB NAME
			name, status := s.Attr("data-title")
			if status == true {
				fmt.Println(name)
			}
		})
		rows.Find("td").Each(func(i int, cell *goquery.Selection) {
			str, status := cell.Attr("name")
			if status == true {
				// We are inside the specific job's cell
				var result string
				var failed_modules []string

				// JOB ID
				job_id := strings.Split(str, "_")
				title, status1 := cell.Find("i").Attr("title")
				if status1 == true {
					state := strings.Split(title, ":")
					if state[0] == "Done" {
						// RESULT
						result = strings.TrimSpace(state[1])
						fmt.Println("The job", job_id[2], "is", state[0], "with result", result)
					} else {
						fmt.Println("The job", job_id[2], "is", state[0])
					}
				}
				if result == "failed" {
					cell.Find("span").Each(func(i int, n *goquery.Selection) {
						failing_module_multiline, exists := n.Attr("title")
						if exists == true {
							failing_module_string := strings.ReplaceAll(failing_module_multiline, "\n", "")
							modules := strings.Split(failing_module_string, "- ")
							if len(modules) > 1 {
								modules = modules[1:]
							}
							// FAILED MODULES
							failed_modules = append(failed_modules, modules...)
							fmt.Printf("To job %s exei ta eksis extra failed modules %s \n", job_id[2], failed_modules)
						}
					})
				}
			}
		})
	})
}

func ParseJson(job_id string) string {
	vars_json := "https://openqa.suse.de/tests/" + job_id + "/file/vars.json"
	resp, err := http.Get(vars_json)
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

		scanner := bufio.NewScanner(strings.NewReader(bodyString))
		for scanner.Scan() {
			line := scanner.Text()

			if strings.Contains(line, `"ARCH" :`) {
				fmt.Println(scanner.Text())
				break
			}
		}
	}
	return "geiaaaaaaaaa"
}

func ParseModules() {
	var modules []string

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
		// fmt.Println(bodyString)

		reached_scheduling := false
		scanner := bufio.NewScanner(strings.NewReader(bodyString))
		for scanner.Scan() {
			line := scanner.Text()

			if strings.Contains(line, "[debug] scheduling") {
				reached_scheduling = true
				testline := strings.Split(line, " ")
				modules = append(modules, strings.Split(testline[len(testline)-1], "tests/")[1])
			} else if reached_scheduling {
				// if there are no more scheduling lines, break
				break
			}
		}
		fmt.Println(modules)
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
