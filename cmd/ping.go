package main

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
	"fmt"
	"net"
	"os"
	"time"

	"github.com/tatsushid/go-fastping"
)

var (
	servername = "example.com"
	pingserver = servername
)

func main() {

	if len(os.Args[1:]) > 0 {
		pingserver = os.Args[1]
	}

	fmt.Println("Attempting to ping " + pingserver)

	p := fastping.NewPinger()
	//ra, err := net.ResolveIPAddr("ip4:icmp", os.Args[1])
	ra, err := net.ResolveIPAddr("ip4:icmp", pingserver)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	p.AddIPAddr(ra)
	p.OnRecv = func(addr *net.IPAddr, rtt time.Duration) {
		fmt.Printf("IP Addr: %s receive, RTT: %v\n", addr.String(), rtt)
	}
	p.OnIdle = func() {
		fmt.Println("finish")
	}
	err = p.Run()
	if err != nil {
		fmt.Println(err)
	}
}
