package protocol

type Job struct {
	ID         int
	Priority   int
	Delay      int
	TimeToRun  int
	TotalBytes int
	Body       string
	Tube       string
}
