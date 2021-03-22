package data

type Result string

const (
	softfailed = "softfailed"
	passed     = "passed"
	failed     = "failed"
	skipped    = "skipped"
)

type Webui struct {
	name            string
	jobGroupFolders []JobGroupFolder
}

type JobGroupFolder struct {
	name      string
	jobGroups []JobGroup
}

type JobGroup struct {
	name string
	jobs []Job
}

type Job struct {
	name    string
	modules []Module
	result  Result
}

type Module struct {
	name   string
	result Result
}
