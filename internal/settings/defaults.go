package settings

const (
	// DefaultLogLevel is the default amount of logging to show.
	defaultLogLevel = "trace"

	// DefaultLogTarget is the default output location of log output.
	defaultLogTarget = "STDOUT"

	// DefaultLogServiceName is the default service name of logging.
	defaultLogServiceName = "provider-admin"

	// DefaultRegisterURL is the default location of the Register service.
	// register:9020 is the value that will work for the integration test system.
	defaultRegisterURL = "http://register:9020"
)
