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
	"sync"
	"time"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrregistermgr"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-common/pkg/nodeid"
)

// FCRP2PServer represents a server handling p2p connection using tcp.
type FCRP2PServer struct {
	start   bool
	name    string
	timeout time.Duration

	// Connection pool
	pool *communicationPool

	// handlers for different message type
	handlers   map[int32]func(reader *FCRServerReader, writer *FCRServerWriter, request *fcrmessages.FCRMessage) error
	requesters map[int32]func(reader *FCRServerReader, writer *FCRServerWriter, args ...interface{}) error
}

// NewFCRP2PServer creates an empty FCRP2PServer.
func NewFCRP2PServer(
	name string,
	registerMgr *fcrregistermgr.FCRRegisterMgr,
	defaultTimeout time.Duration) *FCRP2PServer {
	return &FCRP2PServer{
		start:   false,
		name:    name,
		timeout: defaultTimeout,
		pool: &communicationPool{
			registerMgr:         registerMgr,
			activeGateways:      make(map[string](*communicationChannel)),
			activeGatewaysLock:  sync.RWMutex{},
			activeProviders:     make(map[string](*communicationChannel)),
			activeProvidersLock: sync.RWMutex{},
		},
		handlers:   make(map[int32]func(reader *FCRServerReader, writer *FCRServerWriter, request *fcrmessages.FCRMessage) error),
		requesters: make(map[int32]func(reader *FCRServerReader, writer *FCRServerWriter, args ...interface{}) error),
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
			logging.Error("P2P Server %s has error reading message from %s: %s", s.name, conn.RemoteAddr(), err.Error())
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
				logging.Error("P2P Server %s has error handling message from %s: %s", s.name, conn.RemoteAddr(), err.Error())
				return
			}
		} else {
			// Message is invalid.
			err = sendInvalidMessage(conn, s.timeout)
			if err != nil {
				// Error in tcp communication, drop the connection.
				logging.Error("P2P Server %s has error responding to %s: %s", s.name, conn.RemoteAddr(), err.Error())
				return
			}
		}
	}
}

// RequestGatewayFromGateway uses a given requester to send a request to a given gateway from gateway.
func (s *FCRP2PServer) RequestGatewayFromGateway(id *nodeid.NodeID, msgType int32, args ...interface{}) error {
	if !s.start {
		return errors.New("Server not started")
	}
	requester := s.requesters[msgType]
	if requester == nil {
		return errors.New("No available requester found for given type")
	}
	comm, err := s.pool.getGatewayConn(s.name, id, accessFromGateway)
	if err != nil {
		logging.Error("P2P Server %s has error get gateway connection to %s: %s", s.name, id.ToString(), err.Error())
		return err
	}
	comm.lock.Lock()
	// Call requester to request
	writer := &FCRServerWriter{conn: comm.conn}
	reader := &FCRServerReader{conn: comm.conn}
	err = requester(reader, writer, args)
	comm.lock.Unlock()
	if err != nil {
		// Error that couldn't ignore, remove the connection.
		s.pool.removeActiveGateway(id)
	}
	return err
}

// RequestGatewayFromProvider uses a given requester to send a request to a given gateway from provider.
func (s *FCRP2PServer) RequestGatewayFromProvider(id *nodeid.NodeID, msgType int32, args ...interface{}) error {
	if !s.start {
		return errors.New("Server not started")
	}
	requester := s.requesters[msgType]
	if requester == nil {
		return errors.New("No available requester found for given type")
	}
	comm, err := s.pool.getGatewayConn(s.name, id, accessFromProvider)
	if err != nil {
		logging.Error("P2P Server %s has error get gateway connection to %s: %s", s.name, id.ToString(), err.Error())
		return err
	}
	comm.lock.Lock()
	// Call requester to request
	writer := &FCRServerWriter{conn: comm.conn}
	reader := &FCRServerReader{conn: comm.conn}
	err = requester(reader, writer, args)
	comm.lock.Unlock()
	if err != nil {
		// Error that couldn't ignore, remove the connection.
		s.pool.removeActiveGateway(id)
	}
	return err
}

// RequestProvider uses a given requester to send a request to a given provider. (Only possible from gateway)
func (s *FCRP2PServer) RequestProvider(id *nodeid.NodeID, msgType int32, args ...interface{}) error {
	if !s.start {
		return errors.New("Server not started")
	}
	requester := s.requesters[msgType]
	if requester == nil {
		return errors.New("No available requester found for given type")
	}
	comm, err := s.pool.getProviderConn(s.name, id)
	if err != nil {
		logging.Error("P2P Server %s has error get provider connection to %s: %s", s.name, id.ToString(), err.Error())
		return err
	}
	comm.lock.Lock()
	// Call requester to request
	writer := &FCRServerWriter{conn: comm.conn}
	reader := &FCRServerReader{conn: comm.conn}
	err = requester(reader, writer, args)
	comm.lock.Unlock()
	if err != nil {
		// Error that couldn't ignore, remove the connection.
		s.pool.removeActiveProvider(id)
	}
	return err
}
