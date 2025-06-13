package handlers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/LotusLabsSoftworks/simulacra-server/src/helpers"
)

var (
	simulacraHandlers = map[string]func([]byte, string) ([]byte, error){
		"manifest": SimulacraManifest,
		"asset":    SimulacraAsset,
	}
)

func SimulacraManifest(b []byte, directory string) ([]byte, error) {
	manifest, err := os.ReadFile(directory + "/manifest.yaml")
	if err != nil {
		return nil, errors.New("No Manifest Found")
	}

	response := SimulacraResponsePacket("200", manifest, "plain/yaml")

	return response, nil
}

func SimulacraAsset(b []byte, directory string) ([]byte, error) {
	headerPath, err := helpers.GetHeaderPath(b)

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

func HandleSimulacraRequest(packet []byte, directory string) ([]byte, error) {
	method, err := helpers.GetHeaderMethod(packet)

	if err != nil {
		return nil, err
	}

	response, err := simulacraHandlers[strings.ToLower(method)](packet, directory)

	if err != nil {
		return nil, errors.New("Unknown protocol method")
	}

	return response, nil
}

// Find header end position marke
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

	responseHeader := helpers.WrapHeaderWithMarker([]byte(responseString))
	responseWithPayload := append(responseHeader, payload...)
	return responseWithPayload
}
