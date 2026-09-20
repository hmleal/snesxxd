package rom

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type SNESHeader struct {
	Raw           []byte // Used only for debug
	Title         string
	MapMode       byte
	ROMType       byte
	ROMSize       byte
	RAMSize       byte
	Region        byte
	Developer     byte
	DeveloperName string
	Version       byte
	ChecksumComp  uint16
	Checksum      uint16
}

// NewSNESHeader parses a 32-byte SNES ROM header into an SNESHeader.
//
// It returns an error when data contains fewer than 32 bytes.
func NewSNESHeader(data []byte) (SNESHeader, error) {
	if len(data) < 32 {
		return SNESHeader{}, fmt.Errorf("header too small: got %d bytes, need 32", len(data))
	}

	snesHeader := SNESHeader{
		Raw:           data,
		Title:         strings.TrimRight(string(data[0x00:0x15]), "\x00"),
		MapMode:       data[0x15],
		ROMType:       data[0x16],
		ROMSize:       data[0x17],
		RAMSize:       data[0x18],
		Region:        data[0x19],
		Developer:     data[0x1A],
		DeveloperName: developerName(data[0x1A]),
		Version:       data[0x1B],
		Checksum:      binary.LittleEndian.Uint16(data[0x1E:0x20]),
		ChecksumComp:  binary.LittleEndian.Uint16(data[0x1C:0x1E]),
	}

	return snesHeader, nil
}

func DetectHeaderOffset(data []byte) (int, error) {
	if len(data) < 32 {
		return 0, fmt.Errorf("header too small: got %d bytes, need 32", len(data))
	}

	candidates := []int{
		0x7FC0,       // LoRom
		0xFFC0,       // HiRom
		0x7FC0 + 512, // LoRom + copier header
		0xFFC0 + 512, // HiRom + copier header
	}

	for _, offset := range candidates {
		if offset+32 > len(data) {
			continue
		}

		header := data[offset : offset+32]

		if isValidHeader(header) {
			return offset, nil
		}
	}

	return 0, fmt.Errorf("SNES header not found")
}

func developerName(id byte) string {
	var developers = map[byte]string{
		0x01: "Nintendo",
		0x08: "Capcom",
		0x13: "Electronic Arts",
		0x18: "Hudson Soft",
		0x33: "Ocean",
	}

	if name, ok := developers[id]; ok {
		return fmt.Sprintf("%s (0x%02X)", name, id)
	}

	return fmt.Sprintf("Unknown developer (0x%02X)", id)
}

func regionName(id byte) string {
	var regions = map[byte]string{
		0x00: "Japan",
		0x01: "USA",
	}

	if name, ok := regions[id]; ok {
		return fmt.Sprintf("(0x%02X) %s", id, name)
	}

	return fmt.Sprintf("Unknown region (0x%02X)", id)
	// 00 = Japan
	// 01 = USA
	// 02 = Europe
	// 03 = Sweden
	// 04 = Finland
	// 05 = Denmark
	// 06 = France
	// 07 = Netherlands
	// 08 = Spain
	// 09 = Germany
	// 0A = Italy
	// 0B = China
	// 0C = Korea
	// 0D = Canada
	// 0E = Brazil
	// 0F = Australia
	// 10 = Other
}

func isValidHeader(header []byte) bool {
	mapMode := header[0x15]

	switch mapMode {
	case 0x20, 0x21, 0x30, 0x31:
		return true
	default:
		return false
	}
}
