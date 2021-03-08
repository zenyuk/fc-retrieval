package network

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

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/tatsushid/go-fastping"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
)

// Ping pings a remote server via ICMP
func Ping(pingserver string) bool {

	if len(pingserver) == 0 {
		logging.ErrorAndPanic("Error: Cannot ping empty servername")
		os.Exit(1)
	} else {
		logging.Info("Attempting to ping %s", pingserver)
	}

	p := fastping.NewPinger()
	p.Network("udp")
	ra, err := net.ResolveIPAddr("ip4:icmp", pingserver)
	if err != nil {
		logging.Error1(err)
		return false
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		logging.Info("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
	}
	p.OnIdle = func() {
		logging.Info("finish")
	}
	err = p.Run()
	if err != nil {
		logging.Error1(err)
		return false
	}
	return true
}

// SendMessage sends a message to a given url and obtain a response
func SendMessage(url string, message *fcrmessages.FCRMessage) (*fcrmessages.FCRMessage, error) {
	mJSON, _ := json.Marshal(message)
	logging.Info("Client Manageer sending JSON: %v to url: %v", string(mJSON), url)
	contentReader := bytes.NewReader(mJSON)
	req, _ := http.NewRequest("POST", "http://"+url+"/v1", contentReader)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		logging.Fatal("Error: %+v", err)
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	var data fcrmessages.FCRMessage
	json.NewDecoder(res.Body).Decode(&data)
	return &data, nil
}
