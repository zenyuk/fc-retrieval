package fcrp2pserver

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
	"errors"
	"net"
	"time"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrpaymentmgr"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrregistermgr"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
)

// FCRP2PServer represents a server handling p2p connection using tcp.
type FCRP2PServer struct {
	start   bool
	name    string
	timeout time.Duration

	// Payment manager
	paymentMgr  *fcrpaymentmgr.FCRPaymentMgr
	registerMgr *fcrregistermgr.FCRRegisterMgr

	// handlers for different message type
	handlers   map[int32]func(reader *FCRServerReader, writer *FCRServerWriter, request *fcrmessages.FCRMessage) error
	requesters map[int32]func(reader *FCRServerReader, writer *FCRServerWriter, args ...interface{}) error
}

// NewFCRP2PServer creates an empty FCRP2PServer.
func NewFCRP2PServer(
	name string,
	paymentMgr *fcrpaymentmgr.FCRPaymentMgr,
	registerMgr *fcrregistermgr.FCRRegisterMgr,
	timeout time.Duration) *FCRP2PServer {
	return &FCRP2PServer{
		start:       false,
		name:        name,
		handlers:    make(map[int32]func(reader *FCRServerReader, writer *FCRServerWriter, request *fcrmessages.FCRMessage) error),
		requesters:  make(map[int32]func(reader *FCRServerReader, writer *FCRServerWriter, args ...interface{}) error),
		paymentMgr:  paymentMgr,
		registerMgr: registerMgr,
	}
}

// AddHandler is used to add a handler to the server for a given type.
func (s *FCRP2PServer) AddHandler(msgType int32, handler func(reader *FCRServerReader, writer *FCRServerWriter, request *fcrmessages.FCRMessage) error) *FCRP2PServer {
	if s.start {
		return s
	}
	s.handlers[msgType] = handler
	return s
}

// AddRequester is used to add a requester to the server for a given type.
func (s *FCRP2PServer) AddRequester(msgType int32, requester func(reader *FCRServerReader, writer *FCRServerWriter, args ...interface{}) error) *FCRP2PServer {
	if s.start {
		return s
	}
	s.requesters[msgType] = requester
	return s
}

// Start is used to start the server.
func (s *FCRP2PServer) Start(listenAddr string) error {
	// Start server
	if s.start {
		return errors.New("Server already started")
	}
	ln, err := net.Listen("tcp", ":"+listenAddr)
	if err != nil {
		return err
	}
	go func(ln net.Listener) {
		for {
			conn, err := ln.Accept()
			if err != nil {
				logging.Error(err.Error())
				continue
			}
			logging.Info("P2P server %s has incoming connection from :%s", s.name, conn.RemoteAddr())
			go s.handleIncomingConnection(conn)
		}
	}(ln)
	logging.Info("P2P server %s starts listening on %s for connections.", s.name, listenAddr)
	return nil
}

// handleIncomingConnection handles incomming connection using given handlers.
func (s *FCRP2PServer) handleIncomingConnection(conn net.Conn) {
	// Close connection on exit.
	defer conn.Close()

	// Loop until error occurs and connection is dropped.
	for {
		message, err := readTCPMessage(conn, s.timeout)
		if err != nil {
			// Error in tcp communication, drop the connection.
			logging.Error("P2P Server has error reading message from %s: %s", conn.RemoteAddr(), err.Error())
			return
		}
		handler := s.handlers[message.GetMessageType()]
		if handler != nil {
			// Call handler to handle the request
			writer := &FCRServerWriter{conn: conn}
			reader := &FCRServerReader{conn: conn}
			err = handler(reader, writer, message)
			if err != nil {
				// Error that couldn't ignore, drop the connection.
				logging.Error("P2P Server has error handling message from %s: %s", conn.RemoteAddr(), err.Error())
				return
			}
		} else {
			// Message is invalid.
			err = sendInvalidMessage(conn, s.timeout)
			if err != nil {
				// Error in tcp communication, drop the connection.
				logging.Error("P2P Server has error responding to %s: %s", conn.RemoteAddr(), err.Error())
				return
			}
		}
	}
}

// RequestGatewayFromGateway uses a given requester to send a request to a given gateway from gateway.
func (s *FCRP2PServer) RequestGatewayFromGateway(id *nodeid.NodeID, msgType int, args ...interface{}) error {
	return errors.New("Not yet implemented")
}

// RequestGatewayFromProvider uses a given requester to send a request to a given gateway from provider.
func (s *FCRP2PServer) RequestGatewayFromProvider(id *nodeid.NodeID, msgType int, args ...interface{}) error {
	return errors.New("Not yet implemented")
}

// RequestProvider uses a given requester to send a request to a given provider. (Only possible from gateway)
func (s *FCRP2PServer) RequestProvider(id *nodeid.NodeID, msgType int, args ...interface{}) error {
	return errors.New("Not yet implemented")
}
