package rom

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type SNESHeader struct {
	Raw          []byte // Used only for debug
	Title        string
	MapMode      byte
	ROMType      byte
	ROMSize      byte
	RAMSize      byte
	RegionId     byte
	DeveloperId  byte
	Version      byte
	ChecksumComp uint16
	Checksum     uint16
}

// NewSNESHeader parses a 32-byte SNES ROM header into an SNESHeader.
//
// It returns an error when data contains fewer than 32 bytes.
func NewSNESHeader(data []byte) (SNESHeader, error) {
	if len(data) < 32 {
		return SNESHeader{}, fmt.Errorf("header too small: got %d bytes, need 32", len(data))
	}

	snesHeader := SNESHeader{
		Raw:          data,
		Title:        strings.TrimRight(string(data[0x00:0x15]), "\x00"),
		MapMode:      data[0x15],
		ROMType:      data[0x16],
		ROMSize:      data[0x17],
		RAMSize:      data[0x18],
		RegionId:     data[0x19],
		DeveloperId:  data[0x1A],
		Version:      data[0x1B],
		Checksum:     binary.LittleEndian.Uint16(data[0x1E:0x20]),
		ChecksumComp: binary.LittleEndian.Uint16(data[0x1C:0x1E]),
	}

	return snesHeader, nil
}

func (h *SNESHeader) DeveloperName() string {
	var developers = map[byte]string{
		0x01: "Nintendo",
		0x08: "Capcom",
		0x13: "Electronic Arts",
		0x18: "Hudson Soft",
		0x33: "Ocean",
	}

	if name, ok := developers[h.DeveloperId]; ok {
		return fmt.Sprintf("0x%02X %s", h.DeveloperId, name)
	}

	return fmt.Sprintf("0x%02X Unknown developer", h.DeveloperId)
}

func (h *SNESHeader) RegionName() string {
	var regions = map[byte]string{
		0x00: "Japan",
		0x01: "USA",
		0x02: "Europe",
		0x03: "Sweden",
		0x04: "Finland",
		0x05: "Denmark",
		0x06: "France",
		0x07: "Netherlands",
		0x08: "Spain",
		0x09: "Germany",
		0x0A: "Italy",
		0x0B: "China",
		0x0C: "Korea",
		0x0D: "Canada",
		0x0E: "Brazil",
		0x0F: "Australia",
		0x10: "Other",
	}

	if name, ok := regions[h.RegionId]; ok {
		return fmt.Sprintf("0x%02X %s", h.RegionId, name)
	}

	return fmt.Sprintf("Unknown region 0x%02X", h.RegionId)
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

func isValidHeader(header []byte) bool {
	mapMode := header[0x15]

	switch mapMode {
	case 0x20, 0x21, 0x30, 0x31:
		return true
	default:
		return false
	}
}
