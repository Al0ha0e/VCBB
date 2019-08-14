package master_side

type Dependency struct {
	dependencyJobMeta *JobMeta
	keys              []string
}

/*
type Dependency struct {
	DependencyJob     *Job
	DependencyJobMeta *JobMeta `json:"meta"`
}*/
