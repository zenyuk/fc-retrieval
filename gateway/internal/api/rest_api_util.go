package api
// Copyright (C) 2020 ConsenSys Software Inc

import "regexp"

func isB64OrSimpleAscii(s string) bool {
	regex := regexp.MustCompile("^*([A-Za-z0-9\\-_])$");
	return regex.MatchString(s)
}
