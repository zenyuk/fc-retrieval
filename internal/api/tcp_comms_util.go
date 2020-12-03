package api

import (
	"bufio"
	"encoding/binary"
	"io"
)

func readTCPMessage(reader *bufio.Reader) (byte, []byte, error) {
	// TODO: Here assumes each tcp message starts with
	// (length four bytes big endian, msg_type one byte)
	length := make([]byte, 4)
	_, err := io.ReadFull(reader, length)
	if err != nil {
		return 0, nil, err
	}

	// Get request
	request := make([]byte, binary.BigEndian.Uint32(length))
	_, err = io.ReadFull(reader, request)
	if err != nil {
		return 0, nil, err
	}
	return request[0], request[1:], nil
}

func sendTCPMessage(writer *bufio.Writer, msgType byte, data []byte) error {
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
