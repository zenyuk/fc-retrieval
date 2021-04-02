package providerapi

import (
	"net"

	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrmessages"
	"github.com/ConsenSys/fc-retrieval-common/pkg/fcrtcpcomms"
	"github.com/ConsenSys/fc-retrieval-common/pkg/logging"
	"github.com/ConsenSys/fc-retrieval-gateway/internal/util/settings"
)

// StartProviderAPI starts the TCP API as a separate go routine.
func StartProviderAPI(settings settings.AppSettings) error {
	// Start server
	ln, err := net.Listen("tcp", ":"+settings.BindProviderAPI)
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
			logging.Info("Incoming connection from provider at :%s\n", conn.RemoteAddr())
			go handleIncomingProviderConnection(conn, settings)
		}
	}(ln)
	logging.Info("Listening on %s for connections from Providers\n", settings.BindProviderAPI)
	return nil
}

func handleIncomingProviderConnection(conn net.Conn, settings settings.AppSettings) {
	// Close connection on exit.
	defer conn.Close()

	// Loop until error occurs and connection is dropped.
	for {
		message, err := fcrtcpcomms.ReadTCPMessage(conn, settings.TCPInactivityTimeout)

		if err != nil && !fcrtcpcomms.IsTimeoutError(err) {
			// Error in tcp communication, drop the connection.
			logging.Error(err.Error())
			return
		}
		if err == nil {
			logging.Info("Message received: %+v", message)
			if message.GetMessageType() == fcrmessages.ProviderPublishGroupOfferRequestType {
				err = handleProviderPublishGroupCIDRequest(conn, message, settings)
				if err != nil && !fcrtcpcomms.IsTimeoutError(err) {
					// Error in tcp communication, drop the connection
					logging.Error(err.Error())
					return
				}
				continue
			} else if message.GetMessageType() == fcrmessages.ProviderPublishDHTOfferRequestType {
				err = handleProviderDHTPublishGroupCIDRequest(conn, message, settings)
				if err != nil && !fcrtcpcomms.IsTimeoutError(err) {
					// Error in tcp communication, drop the connection
					logging.Error(err.Error())
					return
				}
				continue
			}
			// Message is invalid.
			err = fcrtcpcomms.SendInvalidMessage(conn, settings.TCPInactivityTimeout)
			if err != nil && !fcrtcpcomms.IsTimeoutError(err) {
				// Error in tcp communication, drop the connection.
				logging.Error(err.Error())
				return
			}
		}
	}
}
