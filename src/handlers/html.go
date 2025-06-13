package handlers

import (
	"fmt"
	"time"
)

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
