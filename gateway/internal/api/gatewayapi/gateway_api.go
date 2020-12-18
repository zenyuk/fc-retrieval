package gatewayapi

import (
	"bufio"
	"encoding/json"
	"net"
	"sync"
	"time"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/tcpcomms"
)

// StartGatewayAPI starts the TCP API as a separate go routine.
func StartGatewayAPI(settings settings.AppSettings, g *gateway.Gateway) error {
	// Start server
	ln, err := net.Listen("tcp", settings.BindGatewayAPI)
	if err != nil {
		return err
	}
	go func(ln net.Listener) {
		for {
			conn, err := ln.Accept()
			if err != nil {
				logging.Error1(err)
				continue
			}
			logging.Info("Incoming connection from gateway at :%s", conn.RemoteAddr())
			go handleGatewayCommunication(conn, g)
		}
	}(ln)
	logging.Info("Listening on %s for connections from Gateways", settings.BindGatewayAPI)
	return nil
}

func handleGatewayCommunication(conn net.Conn, g *gateway.Gateway) {
	// Initialise a reader and a writer
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	// Init struct, Set register status to false
	gComms := gateway.CommunicationChannels{
		CommsLock:                 sync.RWMutex{},
		InterruptRequestChan:      make(chan bool),
		InterruptResponseChan:     make(chan bool),
		CommsRequestChan:          make(chan []byte),
		CommsResponseChan:         make(chan []byte),
		CommsResponseError:        make(chan error),
		CommsResponseErrorIgnored: make(chan bool),
	}
	registered := false
	defer conn.Close()
	// Start a go routine checks if any new message from the other gateway
	triggerChan, controlChan := startDetectingRequest(reader)
	// switch off control chan on exit to close routine.
	defer func(controlChan chan bool) {
		controlChan <- false
	}(controlChan)
	// Start detecting
	controlChan <- true
	// There are two states, listening mode and interrupt mode.
	// By default, the thread is in listening mode.
	listening := true
	// Start loop
	for {
		switch listening {
		case true:
			// This is in listening mode.
			select {
			case <-triggerChan:
				// There is a message received.
				msgType, data, err := tcpcomms.ReadTCPMessage(reader)
				if err != nil {
					// Connection error can not be ignored here. exit the routine.
					logging.Error1(err)
					return
				}
				// No error in reading the message.
				if msgType == messages.GatewayDHTDiscoverRequestType {
					request := messages.GatewayDHTDiscoverRequest{}
					if json.Unmarshal(data, &request) != nil {
						logging.Warn("Message from gateway: %s can not be parsed.\n", conn.RemoteAddr())
						err = tcpcomms.SendInvalidMessage(writer)
						if err != nil {
							// This error can not be ignored. exit the routine.
							logging.Error1(err)
							return
						}
						// Unable to parse json error can be ignored. continue
						controlChan <- true
						continue
					}
					// Request is legal and correct.
					if !registered {
						// Haven’t registered yet.
						err = gateway.RegisterGatewayCommunication(&request.GatewayID, &gComms)
						if err != nil {
							// This gateway can not be registered, exit the routine
							logging.Error1(err)
							return
						}
						// Deregister upon exiting the routine
						defer gateway.DeregisterGatewayCommunication(&request.GatewayID)
					}
					// Handle request
					err = handleGatewayDHTDiscoverRequest(reader, writer, &request)
					if err != nil {
						// Connection error can not be ignored here. exit the routine.
						logging.Error1(err)
						return
					}
				} else {
					logging.Warn("Message from gateway: %s is of wrong type", conn.RemoteAddr())
					err = tcpcomms.SendInvalidMessage(writer)
					if err != nil {
						// Connection error can not be ignored here. exit the routine.
						logging.Error1(err)
						return
					}
				}
				// The request has been handled properly. Resume the detecting routine.
				controlChan <- true
			case interrupt := <-gComms.InterruptRequestChan:
				if interrupt {
					// Pause the detecting routine.
					controlChan <- true
					// Switch the state to interrupt mode.
					listening = false
					// Respond with a true to indicate readiness
					gComms.InterruptResponseChan <- true
				}
			}
		case false:
			// This is in interrupt mode.
			select {
			case interrupt := <-gComms.InterruptRequestChan:
				if !interrupt {
					// Resume the detecting routine.
					controlChan <- true
					// Switch the state to listening mode.
					listening = true
				}
			case request := <-gComms.CommsRequestChan:
				// Assume the internal request is always a formatted request
				err := tcpcomms.SendTCPMessage(writer, request[0], request[1:])
				if err != nil {
					// Error here can not be ignored.
					gComms.CommsResponseError <- err
					logging.Error1(err)
					return
				}
				// Waiting for a response
				msgType, data, err := tcpcomms.ReadTCPMessage(reader)
				if err != nil {
					if _, ok := err.(*tcpcomms.TimeoutError); ok {
						// Timeout can be ignored, continue
						gComms.CommsResponseChan <- nil
						continue
					}
					// Connection error can not be ignored, exit the routine
					gComms.CommsResponseError <- err
					logging.Error1(err)
					return
				}
				// Respond to internal requester
				gComms.CommsResponseChan <- append([]byte{msgType}, data...)
			}
		}
	}
}

// This starts a go-routine that checks if reader's buffer has any content every 100ms
// Returns a channel of boolean for trigger and control
// Send true to control channel to toggle between starting detecting and pausing detecting.
// By default, detecting is paused.
// Each time the go-routine detects any incoming messages, it will send a true signal to trigger channel and the detecting is paused.
func startDetectingRequest(reader *bufio.Reader) (chan bool, chan bool) {
	triggerChan := make(chan bool)
	controlChan := make(chan bool)
	go func(reader *bufio.Reader, triggerChan, controlChan chan bool) {
		pause := true
		for {
			switch pause {
			case true:
				// Pausing state
				control := <-controlChan
				if control {
					// Resume
					pause = false
				} else {
					// End routine.
					return
				}
			case false:
				afterChan := time.After(100 * time.Microsecond)
				select {
				case control := <-controlChan:
					if control {
						// Pause
						pause = true
					} else {
						// End routine
						return
					}
				case <-afterChan:
					// Time to check if reader's buffer has any content.
					if reader.Buffered() > 0 {
						triggerChan <- true
						// Change to pause state
						pause = true
					}
				}
			}
		}
	}(reader, triggerChan, controlChan)
	return triggerChan, controlChan
}
