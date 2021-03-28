package data

type JobGroup struct {
	Path string
	Url  string
}

type Build struct {
	Path string
	Url  string
}

type Job struct {
	Path                string
	Url                 string
	Name                string
	ID                  string
	Result              string
	FailedModuleAliases []string
	ModuleMap           map[string]bool
	Schedule            string
}
