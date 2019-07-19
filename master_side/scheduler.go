package master_side

import (
	"vcbb/blockchain"
)

type Scheduler struct {
	jobs      []*Job
	bcHandler blockchain.BlockChainHandler
}

func NewScheduler(jobs []*Job) *Scheduler {
	return &Scheduler{jobs: jobs}
}

func (this *Scheduler) AppendJob(job *Job) {
	this.jobs = append(this.jobs, job)
}

func (this *Scheduler) Dispatch() {
	for _, job := range this.jobs {
		job.StartSession(this)
	}
}
