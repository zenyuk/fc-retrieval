package gatewayapi

import (
	"bufio"
	"encoding/json"
	"net"
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
	// Init struct, Set register status to false
	gComms := gateway.CommunicationChannels{
		CommsRequestChan:  make(chan []byte),
		CommsResponseChan: make(chan []byte),
	}
	registered := false
	// Initialise a reader and a writer
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	defer conn.Close()
	// Start a go routine checks if any new message from the other gateway
	triggerChan, controlChan := startDetectingRequest(reader)
	// switch off control chan on exit to close routine.
	defer func(controlChan chan bool) {
		controlChan <- false
	}(controlChan)
	// Start detecting
	controlChan <- true
	// Start loop
	for {
		select {
		case <-triggerChan:
			msgType, data, err := tcpcomms.ReadTCPMessage(reader)
			if err != nil {
				if _, ok := err.(*tcpcomms.TimeoutError); ok {
					// A regular timeout, continue
					controlChan <- true
					continue
				}
				// Connection has something wrong, exit the routine
				logging.Error1(err)
				return
			}
			if msgType == messages.GatewayDHTDiscoverRequestType {
				request := messages.GatewayDHTDiscoverRequest{}
				if json.Unmarshal(data, &request) != nil {
					logging.Error("Message from gateway: %s can not be parsed\n", conn.RemoteAddr())
					err = tcpcomms.SendInvalidMessage(writer)
					if err != nil {
						// Connection has something wrong, exit the routine
						logging.Error1(err)
						return
					}
				} else {
					// Request is legal and correct
					if !registered {
						// Haven't registered yet.
						// Register.
						gComms.NodeID = request.GatewayID
						err = gateway.RegisterGatewayCommunication(gComms.NodeID, &gComms)
						if err != nil {
							// This gateway can not be registered, exit the routine
							logging.Error1(err)
							return
						}
						// Deregister upon exiting the routine
						defer gateway.DeregisterProviderCommunication(gComms.NodeID)
					}
					err = handleGatewayDHTDiscoverRequest(reader, writer, &request)
					if err != nil {
						// Connection has something wrong, exit the routine
						logging.Error1(err)
						return
					}
				}
			} else {
				logging.Warn("Message from gateway: %s is of wrong type", conn.RemoteAddr())
				err = tcpcomms.SendInvalidMessage(writer)
				if err != nil {
					// Connection has something wrong, exit the routine
					logging.Error1(err)
					return
				}
			}
			// Resume detecting routine.
			controlChan <- true
		case request := <-gComms.CommsRequestChan:
			// Pause the detecting routine
			controlChan <- true
			// Assume the internal request is already a formatted request.
			err := tcpcomms.SendTCPMessage(writer, request[0], request[1:])
			if err != nil {
				// Connection has something wrong, exit the routine
				// Respond with a nil
				gComms.CommsResponseChan <- nil
				logging.Error1(err)
				return
			}
			// Getting a response
			msgType, data, err := tcpcomms.ReadTCPMessage(reader)
			if err != nil {
				if _, ok := err.(*tcpcomms.TimeoutError); ok {
					// A regular timeout, continue
					gComms.CommsResponseChan <- nil
					controlChan <- true
					continue
				}
				// Connection has something wrong, exit the routine
				gComms.CommsRequestChan <- nil
				logging.Error1(err)
				return
			}
			// Respond to internal requester
			gComms.CommsResponseChan <- append([]byte{msgType}, data...)
			// Resume detecting routine.
			controlChan <- true
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
