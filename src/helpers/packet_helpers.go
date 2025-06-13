package parsers

import (
	"bytes"
	"errors"
	"strings"
)

var (
	headerEndMarker = []byte{0x0d, 0x0a, 0x0d, 0x0a}
)

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
	return 0, errors.New("Marker Wasnt Found")
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
