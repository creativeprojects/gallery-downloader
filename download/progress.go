package download

type Event int

const (
	EventStart Event = iota
	EventProgress
	EventFinished
	EventNotSaving
	EventError
)

type Progress struct {
	FileID     int
	TotalFiles int
	URL        string
	Event      Event
	Err        error
	Downloaded int64
	Wait       int
}
