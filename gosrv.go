package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var port uint
var echo bool
var data uint
var size uint
var count uint
var interval uint

func usage() {

}

func errorMsg(msg string, err error) {
	fmt.Printf(
		"Error: %s ('%v') exiting program.\n\n",
		msg, err)
	os.Exit(0)
}

func parseArgs() {
	flag.UintVar(&port, "port", 12345, "The port number to listen on.")
	flag.BoolVar(&echo, "echo", true, "True if the server should echo the incoming data.")
	flag.UintVar(&data, "data", 0xFE, "The data to send if sending data that is not echo'd.")
	flag.UintVar(&size, "size", 0, "The number of bytes of data to send.")
	flag.UintVar(&count, "c", 0, "Number of times to send data.")
	flag.UintVar(&interval, "i", 1, "Interval in seconds between sending data.")
	flag.Parse()

	fmt.Println("\nUsing following parameters:")
	fmt.Printf(" Listen port:  %d\n", port)
	fmt.Printf(" Echo Data:    %t\n", echo)
	fmt.Printf(" Data to Send: 0x%X\n", data)
	fmt.Printf(" Size of data: %d\n", size)
	fmt.Printf(" Data Count:   %d\n", count)
	fmt.Printf(" Interval:     %d\n", interval)
	fmt.Println("")
}

func writeToPeer(conn net.Conn, buff []byte) {
	peer := conn.RemoteAddr().String()
	n, err := conn.Write(buff)
	if err != nil {
		log.Printf("Failed to write to peer '%s' (%s), terminating connection\n", peer, err)
		conn.Close()
		return
	}
	log.Printf("Wrote %d bytes to peer '%s'.\n", n, peer)
}

func handleConnection(conn net.Conn) {
	buff := make([]byte, 2000)
	payload := bytes.Repeat([]byte{byte(data)}, int(size))
	peer := conn.RemoteAddr().String()
	cycle := time.Second * time.Duration(interval)

	fmt.Printf("Accepting connection from peer '%s'.\n", peer)
	defer conn.Close()

	timeout := time.Now().Add(cycle)
	for {
		conn.SetReadDeadline(timeout)
		n, err := conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Connection to peer '%s' closed.\n", peer)
				return
			} else if os.IsTimeout(err) {
				if count > 0 {
					writeToPeer(conn, payload)
					count--
				}
				timeout = time.Now().Add(cycle)
			} else {
				log.Printf("Failed to read from peer '%s' (%s), terminating connection\n", peer, err)
				return
			}
		} else {
			log.Printf("Read %d bytes from peer '%s'.\n", n, peer)
			writeToPeer(conn, buff[:n])
		}
	}
}

func main() {
	fmt.Println("Running Go Serve...")
	parseArgs()

	addr := ":" + strconv.FormatUint(uint64(port), 10)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		errorMsg("Failed to open listening socket", err)
	}
	for {
		conn, err := ln.Accept()
		tcpcon := conn.(*net.TCPConn)
		tcpcon.SetKeepAlive(false)
		if err != nil {
			log.Printf("Failed to accept incoming connection (%s)\n", err)
		}
		go handleConnection(conn)
	}

}
