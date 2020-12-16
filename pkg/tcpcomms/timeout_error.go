package tcpcomms

// TimeoutError is used when the other party does not respond within a given time,
// it should be ignored.
type TimeoutError struct{}

func (t *TimeoutError) Error() string {
	return "Communication timeout."
}
