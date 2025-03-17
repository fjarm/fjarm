package interceptor

type DelayDuration int

const (
	// DelayDuration_15000ms represents a delay of 15000 milliseconds or 15 seconds. 15 seconds is the request timeout
	// for most clients.
	DelayDuration_15000ms DelayDuration = 15000
)
