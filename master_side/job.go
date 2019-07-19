package master_side

type JobState uint8

const (
	Preparing JobState = iota
	Running
	Finished
)

type Job struct {
	Dependencies []*Dependency
}
