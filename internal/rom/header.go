package rom

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// HeaderMapInfo describes the memory mapping mode and access speed encoded in a
// SNES ROM header byte.
type HeaderMapInfo struct {
	Raw       byte
	IsFastRom bool
	MapMode   byte
	MapValue  byte
	IsValid   bool
}

// NewHeaderMapInfo decodes a SNES ROM mapping byte into HeaderMapInfo.
func NewHeaderMapInfo(headerByte byte) HeaderMapInfo {
	// Bit 4 (0x10) determines ROM speed: 0 = Slow (2.68MHz), 1 = Fast (3.58MHz)
	isFast := (headerByte & 0x10) != 0

	// Bits 0-3 (0x0F) determine the memory mapping mode
	mapBits := headerByte & 0x0F

	// Bits 5-7 (0xE0 shifted right by 5) should equal 0b001 (1)
	fixedBits := (headerByte & 0xE0) >> 5
	isValidHeader := fixedBits == 1

	return HeaderMapInfo{
		Raw:       headerByte,
		IsFastRom: isFast,
		MapMode:   mapBits,
		MapValue:  mapBits,
		IsValid:   isValidHeader,
	}
}

// String returns a human-readable description of the mapping mode and ROM
// speed.
func (h HeaderMapInfo) String() string {
	speedStr := "Slow (2.68 MHz)"
	if h.IsFastRom {
		speedStr = "Fast (3.58 MHz)"
	}

	switch h.MapMode & 0x0F {
	case 0x00:
		return fmt.Sprintf("LoROM - %s", speedStr)
	case 0x01:
		return fmt.Sprintf("HiROM - %s", speedStr)
	case 0x02:
		return fmt.Sprintf("ExHiROM + S-DD1 - %s", speedStr) // Street Fighter Alpha 2
	case 0x03:
		return fmt.Sprintf("SA-1 - %s", speedStr) // Super Mario RPG, Kirby
	case 0x05:
		return fmt.Sprintf("ExHiROM - %s", speedStr) // Tales of Phantasia
	case 0x0A:
		return fmt.Sprintf("SPC7110 - %s", speedStr) // Tengai Makyou Zero
	default:
		return fmt.Sprintf("Unknown/Custom (0x%02X) - %s", h.MapMode, speedStr)
	}
}

// SNESHeader contains the metadata parsed from a 32-byte SNES ROM header.
type SNESHeader struct {
	Raw          []byte // Used only for debug
	Title        string
	MapMode      HeaderMapInfo
	ROMType      byte
	ROMSize      byte
	RAMSize      byte
	RegionId     byte
	DeveloperId  byte
	Version      byte
	ChecksumComp uint16
	Checksum     uint16
}

var developerNames = map[byte]string{
	0x01: "Nintendo",
	0x08: "Capcom",
	0x13: "Electronic Arts",
	0x18: "Hudson Soft",
	0x33: "Extended license code",
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
		MapMode:      NewHeaderMapInfo(data[0x15]),
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

func (h *SNESHeader) ROMTypeInfo() string {
	switch h.ROMType {
	case 0x00:
		return fmt.Sprintf("0x%02X %s", h.ROMType, "ROM only")
	case 0x01:
		return fmt.Sprintf("0x%02X %s", h.ROMType, "ROM + RAM")
	case 0x02:
		return fmt.Sprintf("0x%02X %s", h.ROMType, "ROM + RAM + Battery")
	default:
		return fmt.Sprintf("0x%02X %s", h.ROMType, "Unknown")
	}
}

// DeveloperName returns the developer identifier and name from the header.
func (h *SNESHeader) DeveloperName() string {
	if name, ok := developerNames[h.DeveloperId]; ok {
		return fmt.Sprintf("0x%02X %s", h.DeveloperId, name)
	}

	return fmt.Sprintf("0x%02X Unknown developer", h.DeveloperId)
}

// RegionName returns the region identifier and name from the header.
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

// DetectHeaderOffset finds the offset of a valid SNES ROM header in data.
func DetectHeaderOffset(data []byte) (int, error) {
	if len(data) < 32 {
		return 0, fmt.Errorf("header too small: got %d bytes, need 32", len(data))
	}

	candidates := []int{
		0x7FC0,         // LoROM
		0xFFC0,         // HiROM
		0x40FFC0,       // ExHiROM (Street Fighter Alpha 2, Tales of Phantasia, etc.)
		0x7FC0 + 512,   // LoROM + copier header
		0xFFC0 + 512,   // HiROM + copier header
		0x40FFC0 + 512, // ExHiROM + copier header
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

// isValidHeader reports whether a 32-byte header uses a supported map mode.
func isValidHeader(header []byte) bool {
	checksumComp := binary.LittleEndian.Uint16(header[0x1C:0x1E])
	checksum := binary.LittleEndian.Uint16(header[0x1E:0x20])

	if (checksum ^ checksumComp) != 0xFFFF {
		return false
	}

	mapMode := header[0x15]

	switch mapMode {
	case
		0x20, 0x30, // LoROM
		0x21, 0x31, // HiROM
		0x22, 0x32, // ExLoROM (Hacks/Homebrew)
		0x23, 0x33, // S-DD1 (Street Fighter Alpha 2)
		0x25, 0x35: // ExHiROM (Tales of Phantasia)
		return true
	default:
		return false
	}
}

// func isValidHeader(header []byte) bool {
//     // Indices 28-29 are the Complement Checksum, 30-31 are the Checksum
//     compChecksum := uint16(header[28]) | (uint16(header[29]) << 8)
//     checksum     := uint16(header[30]) | (uint16(header[31]) << 8)
//     // A valid SNES header must always satisfy this bitwise condition
//     return (checksum ^ compChecksum) == 0xFFFF
// }
