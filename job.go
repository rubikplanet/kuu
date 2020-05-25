package kuu

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"os"
	"sync"
)

// DefaultCron (set option 5 cron to convet 6 cron)
var DefaultCron = cron.New(cron.WithSeconds())

var (
	runningJobs   = make(map[cron.EntryID]bool)
	runningJobsMu sync.RWMutex

	jobs   = make(map[cron.EntryID]*Job)
	jobsMu sync.RWMutex
)

// Job
type Job struct {
	Spec        string              `json:"spec" valid:"required"`
	Cmd         func(c *JobContext) `json:"-,omitempty"`
	Name        string              `json:"name" valid:"required"`
	RunAfterAdd bool                `json:"runAfterAdd"`
	EntryID     cron.EntryID        `json:"entryID,omitempty"`
	cmd         func()
}

// JobContext
type JobContext struct {
	name string
	errs []error
	l    *sync.RWMutex
}

func (j *Job) NewJobContext() *JobContext {
	return &JobContext{
		name: j.Name,
		l:    new(sync.RWMutex),
	}
}

func (c *JobContext) Error(err error) {
	c.l.Lock()
	defer c.l.Unlock()

	c.errs = append(c.errs, err)
}

// AddJobEntry
func AddJobEntry(j *Job) error {
	jobsMu.Lock()
	defer jobsMu.Unlock()

	if os.Getenv("KUU_JOB") == "" || j.Cmd == nil {
		return nil
	}

	if _, err := govalidator.ValidateStruct(j); err != nil {
		return err
	}

	cmd := func() {
		runningJobsMu.Lock()
		defer runningJobsMu.Unlock()

		if runningJobs[j.EntryID] {
			return
		}
		runningJobs[j.EntryID] = true
		INFO("----------- Job '%s' start -----------", j.Name)

		c := j.NewJobContext()
		j.Cmd(c)
		if len(c.errs) > 0 {
			for i, err := range c.errs {
				c.errs[i] = errors.Wrap(err, fmt.Sprintf("Job '%s' execute error", j.Name))
			}
			ERROR(c.errs)
		}
		INFO("----------- Job '%s' finish -----------", j.Name)
		runningJobs[j.EntryID] = false
	}
	v, err := DefaultCron.AddFunc(j.Spec, cmd)
	if err == nil {
		j.EntryID = v
		j.cmd = cmd
		jobs[j.EntryID] = j
	}
	return err
}

func runAllRunAfterJobs() {
	jobsMu.RLock()
	defer jobsMu.RUnlock()

	for _, job := range jobs {
		if job.RunAfterAdd {
			job.cmd()
		}
	}
}

// AddJob
func AddJob(spec string, name string, cmd func(c *JobContext)) (cron.EntryID, error) {
	job := Job{
		Spec: spec,
		Name: name,
		Cmd:  cmd,
	}
	INFO(fmt.Sprintf("Add job: %s %s", job.Name, job.Spec))
	err := AddJobEntry(&job)
	return job.EntryID, err
}
