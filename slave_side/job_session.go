package slave_side

type jobRunTimeError struct {
	id  string
	err error
}

type Job struct {
	id       string
	baseTest string
	hardware string
	sch      *Scheduler
}

func NewJob(id, baseTest, hardware string, sch *Scheduler) *Job {
	return &Job{
		id:       id,
		baseTest: baseTest,
		hardware: hardware,
		sch:      sch,
	}
}

func (this *Job) Init() {

}

func (this *Job) StartSession() {

}
