package main

import (
	"bytes"
	"strings"
	// "encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"
)

const ()

var (
	headerEndMarker   = []byte{0x0d, 0x0a, 0x0d, 0x0a}
	markerNotFound    = errors.New("Marker Wasnt Found")
	simulacraHandlers = map[string]func([]byte, string) ([]byte, error){
		"manifest": SimulacraManifest,
		"asset":    SimulacraAsset,
	}
)

func main() {
	workingDir, err := os.Getwd()

	if err != nil {
		fmt.Println(err)
	}

	portFlagPtr := flag.Int("port", 9191, "port for simulacra server")
	folderPathPtr := flag.String("dir", workingDir, "folder with simulacra manifest and assets")

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

func SimulacraManifest(b []byte, directory string) ([]byte, error) {
	manifest, err := os.ReadFile(directory + "/manifest.yaml")
	if err != nil {
		return nil, errors.New("No Manifest Found")
	}

	response := SimulacraResponsePacket("200", manifest, "plain/yaml")

	return response, nil
}

func SimulacraAsset(b []byte, directory string) ([]byte, error) {
	headerPath, err := GetHeaderPath(b)

	if err != nil {
		return nil, err
	}

	file, err := os.ReadFile(directory + "/assets/" + headerPath)

	if err != nil {
		return nil, err
	}

	extension := filepath.Ext(headerPath)
	contentType := fmt.Sprintf("file/%s", extension)

	response := SimulacraResponsePacket("200", file, contentType)

	return response, nil
}

func HandleHttpRequest(_ []byte, _ string) ([]byte, error) {

	responseContent := "Simulacra: 200 OK.\nThis simulacra server dosent provide a simulacra web client.\neither use the official client or use another web client before connecting here"
	responsePacketString := fmt.Sprintf(
		`HTTP/1.1 200 OK
Server: Simulacra
Date: %s
Content-Type: text/plain; charset=UTF-8
Content-Length: %d

%s
`, time.Now(), len([]byte(responseContent)), responseContent)
	return []byte(responsePacketString), nil
}

func HandleSimulacraRequest(packet []byte, directory string) ([]byte, error) {
	method, err := GetHeaderMethod(packet)

	if err != nil {
		return nil, err
	}

	response, err := simulacraHandlers[strings.ToLower(method)](packet, directory)

	if err != nil {
		return nil, errors.New("Unknown protocol method")
	}

	return response, nil
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

		isHttpRequest, err := IsHttpPacket(packet)
		fmt.Println("is http:", isHttpRequest)

		if isHttpRequest {
			res, err := HandleHttpRequest(packet, directory)
			if err != nil {
				fmt.Println(err)

			}
			response = res
			break
		}

		// fmt.Println(string(packet))
		isSimulacraRequest, err := IsSimulacraPacket(packet)
		// fmt.Println("is simulacra:", isSimulacraRequest)

		if isSimulacraRequest {
			res, err := HandleSimulacraRequest(packet, directory)
			if err != nil {
				fmt.Println(err)
			}

			response = res
			break
		}

	}
	c.Write(response)
}

// Find header end position marker
func FindMarker(marker []byte, target []byte) (int, error) {
	for index, element := range target {
		if element != marker[0] {
			continue
		}

		sliceEnd := index + len(marker)
		slice := target[index:sliceEnd]

		if !bytes.Equal(slice, marker) {
			continue
		}

		return index, nil
	}
	return 0, markerNotFound
}

// FindProtocolHeader finds and returns the contents of the protocol header,
// typically the first line of the header.
// returning Method: string, Path: string, Protocol: string, and Error
func FindProtocolHeader(raw string) (string, string, string, error) {
	headerRows := strings.Split(raw, "\n")           //seperate by new line
	protocolRow := strings.Split(headerRows[0], " ") //get protocol row
	if len(protocolRow) != 3 {
		return "", "", "", errors.New("Malformed Protocol Row")
	}

	return protocolRow[0], protocolRow[1], protocolRow[2], nil
}

func GetHeaderFromPacket(packet []byte) (string, error) {
	length, err := FindMarker(headerEndMarker, packet)
	if err != nil {
		return "", err
	}

	headers := string(packet[:length])
	return headers, nil
}

// Get raw header method text
func GetHeaderMethod(packet []byte) (string, error) {

	headers, err := GetHeaderFromPacket(packet)
	if err != nil {
		return "", err
	}
	method, _, _, err := FindProtocolHeader(headers)

	return method, err
}

// get raw header path text
func GetHeaderPath(packet []byte) (string, error) {

	length, err := FindMarker(headerEndMarker, packet)
	if err != nil {
		return "", err
	}

	headers := string(packet[:length])

	_, path, _, err := FindProtocolHeader(headers)

	return path, err
}

// Get raw Header Protocol Text
func GetHeaderProtocol(packet []byte) (string, error) {
	length, err := FindMarker(headerEndMarker, packet)

	if err != nil {
		return "", err
	}

	headers := string(packet[:length])

	_, _, protocol, err := FindProtocolHeader(headers)

	return protocol, err
}

// Test if Packet contains HTTP in the protocol
func IsHttpPacket(packet []byte) (bool, error) {
	protocol, err := GetHeaderProtocol(packet)
	if err != nil {
		return false, err
	}

	if !strings.Contains(protocol, "HTTP") {
		return false, nil
	}

	return true, nil
}

// Test if Packet contains SIMULACRA in the protocol
func IsSimulacraPacket(packet []byte) (bool, error) {
	protocol, err := GetHeaderProtocol(packet)
	if err != nil {
		return false, err
	}

	if !strings.Contains(protocol, "SIMULACRA") {
		return false, nil
	}

	return true, nil
}

func SimulacraResponsePacket(status string, payload []byte, payloadType string) []byte {
	payloadLength := len(payload)
	statusMessage := ""
	responseString := fmt.Sprintf(
		`SIMULACRA/0.1 %s %s
Server: Simulacra
Date: %s
Content-Type: %s
Content-Length: %d`,
		status,
		statusMessage,
		time.Now(),
		payloadType,
		payloadLength,
	)

	responseHeader := append([]byte(responseString), headerEndMarker...)
	responseWithPayload := append(responseHeader, payload...)
	return responseWithPayload
}
