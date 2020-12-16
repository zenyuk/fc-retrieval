package tcpcomms

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"io"
	"time"

	"github.com/ConsenSys/fc-retrieval-gateway/pkg/messages"
)

// ReadTCPMessage read the tcp message from a given reader
func ReadTCPMessage(reader *bufio.Reader) (byte, []byte, error) {
	afterChan := time.After(30 * time.Second) // TODO: The timeout default value is 30, it should be configurable.
	triggerChan := make(chan bool)
	var request []byte
	var err error
	go func(reader *bufio.Reader, triggerChan chan bool) {
		// TODO: Need to check. Here assumes each tcp message starts with
		// (length four bytes big endian, msg_type one byte)
		length := make([]byte, 4)
		_, err = io.ReadFull(reader, length)
		if err != nil {
			triggerChan <- true
		}

		// Get request
		request = make([]byte, binary.BigEndian.Uint32(length))
		_, err = io.ReadFull(reader, request)
		if err != nil {
			triggerChan <- true
		}
		triggerChan <- true
	}(reader, triggerChan)
	select {
	case <-afterChan:
		return 0, nil, &TimeoutError{}
	case <-triggerChan:
		if err != nil {
			return 0, nil, err
		}
		return request[0], request[1:], nil
	}
}

// SendTCPMessage sends a tcp message to a given writer
func SendTCPMessage(writer *bufio.Writer, msgType byte, data []byte) error {
	// TODO: Here assumes each tcp message starts with
	// (length four bytes big endian, msg_type one byte)
	length := make([]byte, 4)
	binary.BigEndian.PutUint32(length, uint32(1+len(data)))
	_, err := writer.Write(append(append(length, msgType), data...))
	if err != nil {
		return err
	}
	return writer.Flush()
}

// SendProtocolMismatch sends a protocol mistmatch message to a given writer
func SendProtocolMismatch(writer *bufio.Writer) error {
	data, _ := json.Marshal(messages.ProtocolMismatchResponse{MessageType: messages.ProtocolMismatch})
	return SendTCPMessage(writer, messages.ProtocolMismatch, data)
}

// SendInvalidMessage sends an invalid message to a given writer
func SendInvalidMessage(writer *bufio.Writer) error {
	data, _ := json.Marshal(messages.InvalidMessageResponse{MessageType: messages.InvalidMessage})
	return SendTCPMessage(writer, messages.InvalidMessage, data)
}
