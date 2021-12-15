package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

const (
	CONN_HOST = "127.0.0.1"
	CONN_PORT = "3260"
	CONN_TYPE = "tcp"
)

const (
	BUF_SIZE = 1024
	OPCODE_BITMASK = 0b10111111
	LOGIN_OPCODE = 0x03
	DSL_BITMASK = 0x00FFFFFF
)


type BHS /*Basic Header Segment*/ struct {
	Opcode byte
	DataSegmentLength uint32
}

type ISCSIPacket struct {
	Header BHS
	Data []byte
}

func readISCSIPacket(buf []byte) (Packet ISCSIPacket) {
	Packet.Header = readBHS(buf)
	Packet.Data = buf[48:48+Packet.Header.DataSegmentLength]
	return
}

func readBHS(buf []byte) (Header BHS) {
	Header.Opcode = buf[0] & OPCODE_BITMASK
	Header.DataSegmentLength = binary.BigEndian.Uint32(buf[4:8]) & DSL_BITMASK
	return
}

func createConn()  (net.Conn, error) {
	l, err := net.Listen(CONN_TYPE, CONN_HOST + ":" + CONN_PORT)
	if err != nil {
		return nil, err
	}
	defer l.Close()
	conn, err := l.Accept()
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func handleRequest(conn net.Conn) error {
	defer conn.Close()
	buf := make([]byte, BUF_SIZE)
	reqLen, err := conn.Read(buf)
	if err != nil {
		return err
	}
	fmt.Printf("Bytes readed: %d\n", reqLen)
	packet := readISCSIPacket(buf)
	if packet.Header.Opcode == LOGIN_OPCODE {
		err = parseLoginReq(packet)
	} else {
		fmt.Println("Not Login")
	}
	return nil
}

func parseLoginReq(packet ISCSIPacket) error {
	fmt.Printf("Data length: %d\n", packet.Header.DataSegmentLength)
	args := make(map[string]string)
	for _, i := range strings.Split(string(packet.Data), "\x00") {
		if len(i) != 0 {
			args[strings.Split(i, "=")[0]] = strings.Split(i, "=")[1]
		}
	}
	for i, j := range args {
		fmt.Printf("%s = %s\n", i, j)
	}
	return nil
}

func main() {
	connection, err := createConn()
	if err != nil {
		fmt.Println(err)
	}
	handleRequest(connection)
}
