package utils

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/ggkioulis/openQA-module-mapper/data"
)

const separator = " > "

var PageNotFoundError = errors.New("Page not found")
var JsonParseError = errors.New("Json could not be parsed")

// var jobGroupBarrier sync.WaitGroup

// var jobBarrier sync.WaitGroup
// var jobChan chan bool

type Webui struct {
	Name string
	Url  string
}

func (webui *Webui) Scrape() {
	defer TimeTrack(time.Now(), webui.Name)
	jobGroups := webui.ParseJobGroups()

	// number of jobs being parsed in parallel
	// jobChan = make(chan bool, 1)

	for _, jobGroup := range jobGroups {
		// jobBarrier.Add(1)

		// Do not run ParallelizeJobs concurrently, results in cancelled requests
		webui.ParallelizeJobs(jobGroup)
	}
	// jobBarrier.Wait()
}

func (webui *Webui) ParallelizeJobs(jobGroup data.JobGroup) {
	// defer jobBarrier.Done()
	// jobChan <- true

	build := webui.ParseBuilds(jobGroup)

	// Skip empty Job Groups
	if build.Url != "" {
		jobs := webui.ParseJobs(build)

		for _, job := range jobs {
			// jobBarrier.Add(1)
			webui.ParseModules(job)
		}
	}

	// <-jobChan
	// jobBarrier.Wait()
}

func (webui *Webui) ParseJobGroups() []data.JobGroup {
	var parent string
	pathPrefix := webui.Name + separator + "Job Groups"
	var jobGroups []data.JobGroup

	// TODO retries
	fmt.Println("ParseJobGroups: ", pathPrefix, " url:", webui.Url)
	document := ParseAndGetDocument(webui.Url)

	document.Find("a.dropdown-item").Each(func(i int, s *goquery.Selection) {
		href, _ := s.Attr("href")
		if strings.Contains(href, "parent") {
			parent = s.Text()
		} else {
			// If the parent job group should not be skipped
			if !data.JobGroupsToSkip[parent] {
				// We found a job group, append it to the end of the path
				path := pathPrefix + separator + parent + separator + s.Text()
				jobGroup := data.JobGroup{
					Path: path,
					Url:  webui.Url + href,
				}

				jobGroups = append(jobGroups, jobGroup)
			}
		}
	})

	return jobGroups
}

func (webui *Webui) ParseBuilds(jobGroup data.JobGroup) data.Build {
	// TODO retries
	fmt.Println("ParseBuilds: ", jobGroup.Path, " url:", jobGroup.Url)
	document := ParseAndGetDocument(jobGroup.Url)

	var build data.Build
	s := document.Find("div.px-2.build-label.text-nowrap").First()
	s.Find("a").Each(func(k int, slc *goquery.Selection) {
		// We found a build, append the build number to the end of the path
		path := jobGroup.Path + separator + strings.TrimSpace(slc.Text())
		href, _ := slc.Attr("href")
		build = data.Build{
			Path: path,
			Url:  webui.Url + href,
		}
	})
	return build
}

func (webui *Webui) ParseJobs(build data.Build) []data.Job {
	var jobs []data.Job

	// TODO retries
	fmt.Println("ParseJobs: ", build.Path, " url:", build.Url)
	document := ParseAndGetDocument(build.Url)

	document.Find("tr").Each(func(i int, rows *goquery.Selection) {
		var jobName string
		rows.Find("span").First().Each(func(i int, s *goquery.Selection) {
			name, status := s.Attr("title")
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
				var failedModuleAliases []string

				job_description_slice := strings.Split(description, "_")
				title, exists := cell.Find("i").Attr("title")
				if exists == true {
					state := strings.Split(title, ":")
					if state[0] == "Done" {
						result = strings.TrimSpace(state[1])
						if result != "skipped" && result != "incomplete" && result != "parallel_failed" && result != "timeout_exceeded" {
							// If job is Done and Not Skipped or Incomplete, or Parallel Failed, get the job data
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
										failedModuleAliases = append(failedModuleAliases, modules...)
									}
								})
							}

							arch, machine, yaml_schedule, err := webui.getArchFromJson(jobId)
							if err != nil {
								if errors.Is(err, JsonParseError) {
									log.Fatal("Error getting Arch from Json for ", jobId, err)
								} else if errors.Is(err, PageNotFoundError) {
									// Do nothing for now, just skip the job
								}
							} else {

								job := data.Job{
									Path:                build.Path + separator + jobName + separator + arch,
									Url:                 webui.Url + "/tests/" + jobId,
									Name:                jobName,
									ID:                  jobId,
									Machine:             machine,
									Yaml_schedule:       yaml_schedule,
									Result:              result,
									FailedModuleAliases: failedModuleAliases,
								}

								jobs = append(jobs, job)
							}
						}
					}
				}
			}
		})
	})
	return jobs
}

func (webui *Webui) getArchFromJson(job_id string) (string, string, string, error) {
	vars_json := webui.Url + "/tests/" + job_id + "/file/vars.json"
	resp, err := http.Get(vars_json)
	if err != nil {
		log.Fatal("Unable to get vars.json for job", job_id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Unable to read json for job", job_id, err)
		}

		bodyString := string(bodyBytes)

		scanner := bufio.NewScanner(strings.NewReader(bodyString))
		var arch string
		var machine string
		var yaml_schedule string

		for scanner.Scan() {
			line := scanner.Text()

			if strings.Contains(line, `"ARCH" :`) {
				arch = strings.Split(line, `"`)[3]
				// return arch, nil
			}
			if strings.Contains(line, `"MACHINE" :`) {
				machine = strings.Split(line, `"`)[3]
			}
			if strings.Contains(line, `"YAML_SCHEDULE" :`) {
				yaml_schedule = strings.Split(line, `"`)[3]
			}
		}

		return arch, machine, yaml_schedule, nil
	} else {
		// Scan for page not found
		document := ParseAndGetDocument(vars_json)
		name, _ := document.Find("h1").Html()
		if name == "Page not found" {
			return "", "", "", PageNotFoundError
		}
	}

	return "", "", "", JsonParseError
}

func (webui *Webui) ParseModules(job data.Job) {
	// defer jobBarrier.Done()
	// jobChan <- true

	job.Schedule = ""
	job.ModuleMap = make(map[string]bool)

	autoinst_log := job.Url + "/file/autoinst-log.txt"

	resp, err := http.Get(autoinst_log)
	if err != nil {
		log.Fatal("Unable to get autoinst-log.txt from", job.Url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("Unabe to read autoinst-log.txt from", job.Url, err)
		}
		bodyString := string(bodyBytes)

		scanner := bufio.NewScanner(strings.NewReader(bodyString))
		for scanner.Scan() {
			line := scanner.Text()

			if strings.Contains(line, "[debug] scheduling") {
				testline := strings.Split(line, " tests/")

				// if module is in tests, add it
				// we ignore lib modules that are being run, like sle-15-SP3-Online-aarch64-Build163.1-lynis_gnome
				if len(testline) > 1 {
					moduleName := testline[1]
					moduleAlias := strings.Split(testline[0], "scheduling ")[1]

					if _, ok := job.ModuleMap[moduleName]; !ok {
						// module not yet registered
						job.ModuleMap[moduleName] = true
						strippedName := strings.Split(moduleName, ".pm")[0]
						job.Schedule += "tests/" + strippedName + ","

						for _, failedModule := range job.FailedModuleAliases {
							// failedModules contains the modules by their aliases
							if failedModule == moduleAlias {
								// mark this module as failed
								job.ModuleMap[moduleName] = false
							}
						}
					}
				}
			}
		}
		// If we can parse the autoinst-log.txt, create an entry
		reportJobResults(job)
	}
	// <-jobChan
}

func ParseAndGetDocument(uri string) *goquery.Document {
	// Make HTTP request
	response, err := http.Get(uri)
	if err != nil {
		log.Fatal("Unable to get document from", uri, err)
	}
	defer response.Body.Close()

	// Create a goquery document from the HTTP response
	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Unable to create a goquery document from", uri, err)
		return nil
	}

	return document
}

func reportJobResults(job data.Job) {
	fmt.Println("Parsed job:", job.Url, "| with path:", job.Path)
	fmt.Println("Machine: ", job.Machine)
	fmt.Println("YAML_SCHEDULE: ", job.Yaml_schedule)
	fmt.Println("Schedule: ", job.Schedule)
	// fmt.Println("ModuleMap: ", job.ModuleMap)
	fmt.Printf("Failed Modules: ")
	for key, val := range job.ModuleMap {
		if !val {
			fmt.Printf("%s ", key)
		}
	}
	fmt.Println("\n-----------------------------------")
}
