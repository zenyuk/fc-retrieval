package settings

// Copyright (C) 2020 ConsenSys Software Inc

// Filecoin Retrieval Gateway Admin Client Settings

import (
	"time"
)

const (
	// DefaultTCPInactivityTimeout is the default TCP timeout
	DefaultTCPInactivityTimeout = 100 * time.Millisecond


	// DefaultEstablishmentTTL is the default Time To Live used with Client - Gateway estalishment messages.
	defaultEstablishmentTTL = int64(100)

	// DefaultLogLevel is the default amount of logging to show.
	defaultLogLevel = "trace"

	// DefaultLogServiceName is the default service name of logging.
	defaultLogServiceName = "gateway-admin"

	// DefaultLogTarget is the default output location of log output.
	defaultLogTarget = "STDOUT"
)
