package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/LotusLabsSoftworks/simulacra-server/src/handlers"
	"github.com/LotusLabsSoftworks/simulacra-server/src/helpers"
	"io"
	"log"
	"net"
	"os"
)

const ()

var (
	markerNotFound = errors.New("Marker Wasnt Found")
)

func main() {
	workingDir, err := os.Getwd()

	if err != nil {
		fmt.Println(err)
	}

	portFlagPtr := flag.Int("port", 9191, "port for simulacra server")
	folderPathPtr := flag.String("dir", workingDir, "folder with simulacra manifest and assets")

	flag.Parse()

	port := ":" + fmt.Sprintf("%v", *portFlagPtr)
	directory := *folderPathPtr

	l, err := net.Listen("tcp4", port)
	if err != nil {
		log.Fatal(err)
	}

	defer l.Close()

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err)
		}

		go HandleConnection(c, directory)
	}
}

func HandleConnection(c net.Conn, directory string) {
	fmt.Println("Serving ", c.RemoteAddr().String())
	packet := make([]byte, 4096)
	packetSegment := make([]byte, 4096)
	response := make([]byte, 4096)
	defer c.Close()

	for {
		_, err := c.Read(packetSegment)
		if err != nil {
			if err != io.EOF {
				fmt.Println("read error: ", err)
			}
			break
		}

		packet = bytes.Trim(
			append(
				bytes.Trim(packet, "\x00"),
				packetSegment...,
			),
			"\x00",
		)

		fmt.Println(string(packet))

		isHttpRequest, err := helpers.IsHttpPacket(packet)
		fmt.Println("is http:", isHttpRequest)

		if isHttpRequest {
			res, err := handlers.HandleHttpRequest(packet, directory)
			if err != nil {
				fmt.Println(err)

			}
			response = res
			break
		}

		// fmt.Println(string(packet))
		isSimulacraRequest, err := helpers.IsSimulacraPacket(packet)
		// fmt.Println("is simulacra:", isSimulacraRequest)

		if isSimulacraRequest {
			res, err := handlers.HandleSimulacraRequest(packet, directory)
			if err != nil {
				fmt.Println(err)
			}

			response = res
			break
		}

	}
	c.Write(response)
}
