package status

type Status int

const (
	Exists Status = 1 + iota
	Created
	Deleted
	DoesNotExist
	InProgress
	Error
)
