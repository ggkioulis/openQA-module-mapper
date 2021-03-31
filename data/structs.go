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

// Job Groups to skip
var JobGroupsToSkip = map[string]bool{
	"Maintenace: KOTD": true,
	"L3":               true,
	"Development":      true,
	"Released":         true,
	// "SLE 15":                               true,
	// "SLES JeOS":                            true,
	// "Containers":                           true,
	// "SLE Micro":                            true,
	// "Maintenance: Single Incidents":        true,
	// "Maintenance: Test Repo":               true,
	// "Maintenance: Single Incidents SLE-12": true,
	// "Public Cloud":                         true,
}
