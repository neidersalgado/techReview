package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

func main() {
	// create a socket server (listener) open a listener on TCP 9000 port.
	// with host empty listen to all interfaces aviable on the red (0.0.0.0)
	//ln is net.listener, accept entry connections (ln.accept)
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Listening on :9000")
	fmt.Println("Press Ctrl-C to exit")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("accept:", err)
			continue
		}
		// in each go rutina handle the close for conn inside handle
		go handle(conn)
	}

}

func handle(conn net.Conn) {
	defer conn.Close()
	// avoid reading byte by byte from the socket(expensive, syscalls)
	//rad large chunks int ram and server from there (faster)
	r := bufio.NewReader(conn)
	// does not write byte directly to the socket.
	// store data in ram and flushes in bulk(when full or flush)
	w := bufio.NewWriter(conn)

	for {
		// 1. read just the length of the message

		// create a slice of byte, 4 of length and then 32 bits.
		lenBuf := make([]byte, 4)
		//read full read exactly the length buf, block until get these 4 bytes, or until an error happens.
		//r.Read() can return less than requested, but read full guarantees fulfill the buffer or fails.
		if _, err := io.ReadFull(r, lenBuf); err != nil {
			if err != io.EOF {
				log.Println(" read len:", err)
			}
			return
		}

		// cast the 4 byte in a unit 32 int, BigEnding first byte in order the most significative at beginning [00 00 01 F4] → 500
		n := binary.BigEndian.Uint32(lenBuf)
		// this n number from the bug long, says the length of message, and then do a bit shifting to calculate 10MB (10*2**20)
		if n > 10<<20 {
			log.Println("frame too large")
			return
		}

		// 2. read payload

		// set and reserve the specific space for n (length of a message)
		data := make([]byte, n)
		//read and fulfill data with the space in the byte slice calculated with the header
		if _, err := io.ReadFull(r, data); err != nil {
			if err != io.EOF {
				log.Println("read payload:", err)
			}
		}
		fmt.Println("frame size:", len(data))
		fmt.Println("data:", string(data))
		// 3. echo
		if err := writeFrame(w, data); err != nil {
			log.Println("write frame:", err)
			return
		}
	}
}

func writeFrame(w *bufio.Writer, b []byte) error {
	lenBuf := make([]byte, 4)
	// encode the length of the payload into the header (BigEndian = network byte order)
	binary.BigEndian.PutUint32(lenBuf, uint32(len(b)))
	// write the header first (so the receiver knows how many bytes to expect)
	if _, err := w.Write(lenBuf); err != nil {
		return err
	}

	// write the actual payload
	if _, err := w.Write(b); err != nil {
		return err
	}

	//flush the buffer to ensure data is sent immediately
	return w.Flush()
}
