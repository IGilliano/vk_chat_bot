package events

type Fetcher interface {
	Fetch() ([]Event, error)
}

type Processor interface {
	Process(e Event) error
}

type Event struct {
	Type string
	Text string
	Meta interface{}
}
