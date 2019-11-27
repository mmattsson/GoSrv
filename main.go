package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strconv"
)

var port uint
var data uint
var echo bool
var count uint

func usage() {

}

func parseArgs() {
	flag.UintVar(&port, "port", 12345, "The port number to listen on.")
	flag.BoolVar(&echo, "echo", true, "True if the server should echo the incoming data.")
	flag.UintVar(&data, "data", 0xFE, "The data to send if sending data that is not echo'd.")
	flag.UintVar(&count, "c", 0, "Number of times to send data")
	flag.Parse()

	fmt.Println("\nUsing following parameters:")
	fmt.Printf(" Listen port:  %d\n", port)
	fmt.Printf(" Echo Data:    %t\n", echo)
	fmt.Printf(" Data to Send: 0x%X\n", data)
	fmt.Printf(" Data Count:   %d\n", count)
	fmt.Println("")
}

func handleConnection(conn net.Conn) {
	buff := make([]byte, 2000)

	fmt.Printf("Accepting connection from peer '%s'.\n", conn.RemoteAddr().String())
	defer conn.Close()

	for {
		n, err := conn.Read(buff)
		if err != nil {
			log.Printf("Failed to read from peer '%s', terminating connection\n", conn.RemoteAddr().String())
			return
		}
		log.Printf("Read %d bytes from peer '%s'.\n", n, conn.RemoteAddr().String())
		n, err = conn.Write(buff[:n])
		if err != nil {
			log.Printf("Failed to write to peer '%s', terminating connection\n", conn.RemoteAddr().String())
			return
		}
		log.Printf("Wrote %d bytes to peer '%s'.\n", n, conn.RemoteAddr().String())
	}
}

func main() {
	fmt.Println("Running Go Serve...")
	parseArgs()

	addr := ":" + strconv.FormatUint(uint64(port), 10)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Panicf("Failed to open listening socket (%s)\n", err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection (%s)\n", err)
		}
		go handleConnection(conn)
	}

}
