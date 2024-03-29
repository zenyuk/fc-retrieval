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
