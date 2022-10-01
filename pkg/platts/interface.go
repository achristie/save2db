package platts

// Allow pages to be fetched concurrently
type Concurrentable interface {
	GetTotalPages() int
}

// Allow records to be stored in a service
type Writeable interface {
	GetResults() []interface{}
	GetTotalPages() int
}

type Result[T Concurrentable] struct {
	Message T
	Err     error
}
