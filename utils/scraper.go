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

const separator = " > "

// TODO add in data module
type Webui struct {
	Name string
	Url  string
}

type JobGroup struct {
	Path string
	Url  string
}

type Build struct {
	Path string
	Url  string
}

type Job struct {
	Path          string
	Url           string
	Name          string
	ID            string
	result        string
	failedModules []string
}

func (webui *Webui) Scrape() {
	jobGroups := webui.ParseJobGroups()

	for _, jobGroup := range jobGroups {
		build := webui.ParseBuilds(jobGroup)
		jobs := webui.ParseJobs(build)

		for _, job := range jobs {
			modules := webui.ParseModules(job.Url)
			fmt.Println(modules)
			break
		}
		break
	}
}

func (webui *Webui) ParseJobGroups() []JobGroup {
	var parent string
	pathPrefix := webui.Name + separator + "Job Groups"
	var jobGroups []JobGroup

	document := ParseAndGetDocument(webui.Url)
	document.Find("a.dropdown-item").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if strings.Contains(href, "parent") {
			parent = s.Text()
		} else {
			// We found a job group, append it to the end of the path
			path := pathPrefix + separator + parent + separator + s.Text()
			jobGroup := JobGroup{
				Path: path,
				Url:  webui.Url + href,
			}

			jobGroups = append(jobGroups, jobGroup)
		}
	})

	return jobGroups
}

func (webui *Webui) ParseBuilds(jobGroup JobGroup) Build {
	document := ParseAndGetDocument(jobGroup.Url)
	var build Build
	s := document.Find("div.px-2.build-label.text-nowrap").First()
	s.Find("a").Each(func(k int, slc *goquery.Selection) {
		// We found a build, append the build number to the end of the path
		path := jobGroup.Path + separator + strings.TrimSpace(slc.Text())
		href, _ := slc.Attr("href")
		build = Build{
			Path: path,
			Url:  webui.Url + href,
		}
	})
	return build
}

func (webui *Webui) ParseJobs(build Build) []Job {
	var jobs []Job

	document := ParseAndGetDocument("https://openqa.suse.de/tests/overview?distri=sle&version=15-SP3&build=163.1&groupid=110")
	document.Find("tr").Each(func(i int, rows *goquery.Selection) {
		var jobName string
		rows.Find("a").First().Each(func(i int, s *goquery.Selection) {
			name, status := s.Attr("data-title")
			if status == true {
				jobName = name
			}
		})
		rows.Find("td").Each(func(i int, cell *goquery.Selection) {
			description, status := cell.Attr("name")
			if status == true {
				// We are inside the specific job's cell
				var jobId string
				var result string
				var failedModules []string

				job_description_slice := strings.Split(description, "_")
				title, exists := cell.Find("i").Attr("title")
				if exists == true {
					state := strings.Split(title, ":")
					if state[0] == "Done" {
						result = strings.TrimSpace(state[1])
						if result != "skipped" {
							// If job is Done and Not Skipped, get the job data
							jobId = job_description_slice[2]

							if result == "failed" {
								cell.Find("span").Each(func(i int, n *goquery.Selection) {
									failing_module_multiline, exists := n.Attr("title")
									if exists == true {
										failing_module_string := strings.ReplaceAll(failing_module_multiline, "\n", "")
										modules := strings.Split(failing_module_string, "- ")
										if len(modules) > 1 {
											modules = modules[1:]
										}
										failedModules = append(failedModules, modules...)
									}
								})
							}

							arch, err := getArchFromJson(jobId)
							if err != nil {
								log.Fatal(err)
							}

							job := Job{
								Path:          build.Path + separator + jobName + separator + arch,
								Url:           webui.Url + "/tests/" + jobId,
								Name:          jobName,
								ID:            jobId,
								result:        result,
								failedModules: failedModules,
							}

							jobs = append(jobs, job)
						}
					}
				}
			}
		})
	})
	return jobs
}

func getArchFromJson(job_id string) (string, error) {
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
				arch := strings.Split(line, `"`)[3]
				return arch, nil
			}
		}
	}
	return "", fmt.Errorf("could not parse json file")
}

func (webui *Webui) ParseModules(url string) []string {
	var modules []string

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
	}
	return modules
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
