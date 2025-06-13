package helpers

var (
	headerEndMarker = []byte{0x0d, 0x0a, 0x0d, 0x0a}
)

func WrapHeaderWithMarker(h []byte) []byte {
	return append(h, headerEndMarker...)
}
