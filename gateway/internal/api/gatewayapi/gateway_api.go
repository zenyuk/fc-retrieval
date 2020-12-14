package gatewayapi

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"time"

	"github.com/ConsenSys/fc-retrieval-gateway/internal/gateway"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
	"github.com/ConsenSys/fc-retrieval-gateway/pkg/tcpcomms"
)

// StartGatewayAPI starts the TCP API as a separate go routine.
func StartGatewayAPI(settings util.AppSettings, g *gateway.Gateway) error {
	// Start server
	ln, err := net.Listen("tcp", settings.BindGatewayAPI)
	if err != nil {
		return err
	}
	go func(ln net.Listener) {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Println(err.Error())
				continue
			}
			log.Printf("Incoming connection from gateway at :%s\n", conn.RemoteAddr())
			go handleGatewayCommunication(conn, g)
		}
	}(ln)
	log.Printf("Listening on %s for connections from Gateways\n", settings.BindGatewayAPI)
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
				// Connection has something wrong, exit the routine
				log.Println(err.Error())
				return
			}
			if msgType == messages.GatewayDHTDiscoverRequestType {
				request := messages.GatewayDHTDiscoverRequest{}
				if json.Unmarshal(data, &request) != nil {
					log.Printf("Message from gateway: %s can not be parsed\n", conn.RemoteAddr())
					err = tcpcomms.SendInvalidMessage(writer)
					if err != nil {
						// Connection has something wrong, exit the routine
						log.Println(err.Error())
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
							log.Println(err.Error())
							return
						}
						// Deregister upon exiting the routine
						defer gateway.DeregisterProviderCommunication(gComms.NodeID)
					}
					err = handleGatewayDHTDiscoverRequest(reader, writer, &request)
					if err != nil {
						// Connection has something wrong, exit the routine
						log.Println(err.Error())
						return
					}
				}
			} else {
				log.Printf("Message from gateway: %s is of wrong type\n", conn.RemoteAddr())
				err = tcpcomms.SendInvalidMessage(writer)
				if err != nil {
					// Connection has something wrong, exit the routine
					log.Println(err.Error())
					return
				}
			}
		case request := <-gComms.CommsRequestChan:
			// Do something about the internal requeest
			log.Printf("Internal request: %s\n", request)
			// Send the response to the internal requester
			response := []byte{1, 2, 3}
			gComms.CommsResponseChan <- response
		}
	}
}

// This starts a go-routine that checks if reader's buffer has any content every 100ms
// Returns a channel of boolean for trigger and control
// Send true to control channel to toggle between starting detecting and pausing detecting.
// By default, detecting is paused. Each time
// Send the channel with a true to toggle between starting detecting and pausing detecting, by default, the detecting is paused.
// Send the channel with a false to shutdown the go routine.
// The channel will gets a true signal if there is content detected.
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
