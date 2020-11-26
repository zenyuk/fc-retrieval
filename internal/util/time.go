package util
// Copyright (C) 2020 ConsenSys Software Inc

import (
	"time"
)

// GetTimeNowString returns the time now in a standard format.
func GetTimeNowString() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05")
}
