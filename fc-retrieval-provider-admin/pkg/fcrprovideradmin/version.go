package fcrprovideradmin

/*
 * Copyright 2020 ConsenSys Software Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on
 * an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations under the License.
 *
 * SPDX-License-Identifier: Apache-2.0
 */

const ver = "1.0"
const build = "unknown"

// VersionInfo holds the version information for the Filecoin Retrieval Client library.
type VersionInfo struct {
	Version   string
	BuildDate string
}

// GetVersion returns the static build information.
func GetVersion() VersionInfo {
	return VersionInfo{ver, build}
}
