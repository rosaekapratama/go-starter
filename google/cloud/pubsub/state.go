package pubsub

type State string

const (
	StateStart   State = "start"
	StateSuccess State = "success"
	StateFailed  State = "failed"
	StateAny     State = "*"
)
