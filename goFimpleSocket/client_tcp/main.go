package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	d := net.Dialer{Timeout: 2 * time.Second, KeepAlive: 30 * time.Second}
	conn, err := d.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("connected to :9000")
	defer conn.Close()

	msg := []byte("hola tcp con framing")
	fmt.Println("Sending message :9000")
	if err := writeFrame(conn, msg); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Reading message :9000")
	echo, err := readFrame(conn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Received message :9000")
	fmt.Println("echo:", string(echo))
}

func writeFrame(conn net.Conn, b []byte) error {
	lenBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(b)))
	if _, err := conn.Write(lenBuf); err != nil {
		return err
	}
	_, err := conn.Write(b)
	return err
}

func readFrame(conn net.Conn) ([]byte, error) {
	r := bufio.NewReader(conn)
	lenBuf := make([]byte, 4)
	if _, err := io.ReadFull(r, lenBuf); err != nil {
		return nil, err
	}
	n := binary.BigEndian.Uint32(lenBuf)
	data := make([]byte, n)
	if _, err := io.ReadFull(r, data); err != nil {
		return nil, err
	}
	return data, nil
}
