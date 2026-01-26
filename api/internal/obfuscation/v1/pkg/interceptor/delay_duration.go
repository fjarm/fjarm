package interceptor

type DelayDuration int

const (
	// DelayDuration_15000ms represents a delay of 15000 milliseconds or 15 seconds. 15 seconds is the request timeout
	// for most clients.
	DelayDuration_15000ms DelayDuration = 15000

	// DelayDuration_100ms represents a delay of 100 milliseconds or 0.1 seconds.
	DelayDuration_100ms DelayDuration = 100

	// DelayDuration_500ms represents a delay of 500 milliseconds or 0.5 seconds.
	DelayDuration_500ms DelayDuration = 500

	// DelayDuration_1000ms represents a delay of 1000 milliseconds or 1 second.
	DelayDuration_1000ms DelayDuration = 1000
)
