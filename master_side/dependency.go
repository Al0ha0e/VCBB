package master_side

type Dependency struct {
	DependencyJob     *Job
	DependencyJobMeta *JobMeta `json:"meta"`
}
