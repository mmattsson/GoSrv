package main

import (
	"flag"
	"fmt"
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

func main() {
	fmt.Println("Running Go Serve...")
	parseArgs()

}
