package util
// Copyright (C) 2020 ConsenSys Software Inc

const ver = "1.0"
const build = "unknown"

// VersionInfo holds the version information for the applicaiton.
type VersionInfo struct {
	Version   string
	BuildDate string
}

// GetVersion returns the static build information.
func GetVersion() VersionInfo {
	return VersionInfo{ver, build}
}


