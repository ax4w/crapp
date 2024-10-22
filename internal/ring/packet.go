package ring

import "fmt"

func byteToBcd(b int) int {
	tens := b / 10
	ones := b % 10
	return (tens << 4) | ones
}

func makePacket(command int, subData []byte) ([]byte, error) {
	if command < 0 || command > 255 {
		return nil, fmt.Errorf("Invalid command, must be between 0 and 255")
	}

	if subData != nil && len(subData) > 14 {
		return nil, fmt.Errorf("Sub data must be less than or equal to 14 bytes")
	}

	packet := make([]byte, 16)
	packet[0] = byte(command)

	if subData != nil {
		// Copy subData into packet
		copy(packet[1:], subData)
		println("copied packet data")
	}
	// Calculate checksum and place it in the last byte
	packet[len(packet)-1] = checksum(packet)

	return packet, nil
}

// Checksum calculates the checksum by summing all bytes in the packet modulus 255.
func checksum(packet []byte) byte {
	var sum int
	for _, b := range packet {
		sum += int(b)
	}
	return byte(sum & 255)
}
