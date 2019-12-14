package protocol

type Storage interface {
	Save(job *Job) (int, error)
	Delete(job *Job) bool
	Bury(job *Job) bool
	Reserve(tubes []string) (*Job, error)
	UpdateJobs()
}
