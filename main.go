package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"unicode"

	"snesxxd/internal/rom"

	"github.com/urfave/cli/v3"
)

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	White  = "\033[37m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
)

// isPrintableASCII reports whether b is a printable ASCII byte.
//
// Printable ASCII bytes are in the range 32 (space) through 126 (~).
func isPrintableASCII(b byte) bool {
	return b >= 32 && b <= 126
}

// isWhitespaceControl reports whether b is a whitespace control byte.
//
// These are rendered in yellow and include tab, line feed, and carriage return.
func isWhitespaceControl(b byte) bool {
	return b == 0x09 || b == 0x0A || b == 0x0D
}

// isNullBytes reports whether b is a null byte.
func isNullBytes(b byte) bool {
	return b == 0x00
}

// isMaxByteValue reports whether b is the maximum byte value.
func isMaxByteValue(b byte) bool {
	return b == 0xFF
}

// isPrintableUnicode reports whether b is a printable Unicode rune.
func isPrintableUnicode(b byte) bool {
	return unicode.IsPrint(rune(b))
}

// getColoredCode returns the hexadecimal value of b with a color based on its class.
func getColoredCode(b byte) string {
	switch {
	case isPrintableASCII(b):
		return fmt.Sprintf("%s%02x%s", Green, b, Reset)
	case isWhitespaceControl(b):
		return fmt.Sprintf("%s%02x%s", Yellow, b, Reset)
	case isNullBytes(b):
		return fmt.Sprintf("%s%02x%s", White, b, Reset)
	case isMaxByteValue(b):
		return fmt.Sprintf("%s%02x%s", Blue, b, Reset)
	default:
		return fmt.Sprintf("%s%02x%s", Red, b, Reset)
	}
}

// TUI -  Terminal User Interface
func printHeader(title string) {
	fmt.Println("")
	fmt.Println(title)
	fmt.Println("-------------------------------------------------------------")
}

func printData(h rom.SNESHeader) {
	const labelWidth = 19

	// fmt.Printf("%-*s : %s\n", labelWidth, "Filename", filename)
	// fmt.Printf("%-*s : %d bytes\n", labelWidth, "Filesize", fileInfo.Size())
	fmt.Printf("%-*s : %s\n", labelWidth, "Title", h.Title)
	fmt.Printf("%-*s : %s\n", labelWidth, "Developer", h.DeveloperName)
	fmt.Printf("%-*s : 0x%02X\n", labelWidth, "Map Mode", h.MapMode)
	fmt.Printf("%-*s : 0x%02X\n", labelWidth, "ROM Type", h.ROMType)
	fmt.Printf("%-*s : 0x%02X\n", labelWidth, "ROM Size Exponent", h.ROMSize)
	fmt.Printf("%-*s : 0x%02X\n", labelWidth, "RAM Size Exponent", h.RAMSize)
	fmt.Printf("%-*s : %d\n", labelWidth, "Region", h.Region)
	fmt.Printf("%-*s : 0x%04X\n", labelWidth, "Checksum", h.Checksum)
	fmt.Printf("%-*s : 0x%04X\n", labelWidth, "Checksum Complement", h.ChecksumComp)
	fmt.Printf("%-*s : %d\n", labelWidth, "Raw (Debug)", h.Raw)
	fmt.Println("")
}

func runInfo(filename string) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("runInfo: %q - %w", filename, err)
	}

	offset, err := rom.DetectHeaderOffset(data)
	if err != nil {
		return fmt.Errorf("runInfo: %w", err)
	}

	header, err := rom.NewSNESHeader(data[offset : offset+32])
	if err != nil {
		return fmt.Errorf("runInfo: %w", err)
	}

	printHeader("File")
	printData(header)

	return nil
}

func runHexDump(filename string, offset int) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("runHexDump: %w", err)
	}

	buffer := make([]byte, 16)

	for {
		bytesRead, err := file.ReadAt(buffer, int64(offset))

		if bytesRead > 0 {
			fmt.Printf("%08x: ", offset)
			for i := range bytesRead {
				fmt.Printf("%s", getColoredCode(buffer[i]))
				if i%2 == 1 {
					fmt.Printf(" ")
				}
			}
			fmt.Println()
			offset += bytesRead
		}

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("runHexDump: %w", err)
		}
	}

	return nil
}

func main() {
	cmd := &cli.Command{
		Name:    "snesxxd",
		Usage:   "A terminal-based SNES ROM inspector written in Go",
		Version: "0.0.1",
		Commands: []*cli.Command{
			{
				Name:  "info",
				Usage: "Display SNES ROM header information",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "filename", Required: true},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					if err := runInfo(cmd.StringArg("filename")); err != nil {
						return err
					}

					return nil
				},
			},
			{
				Name:  "hex",
				Usage: "Display a hexadecimal dump of the ROM",
				Arguments: []cli.Argument{
					&cli.StringArg{Name: "filename", Required: true},
				},
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "offset",
						Value: 0,
						Usage: "starting byte offset for reading the ROM",
					},
				},
				Action: func(ctx context.Context, cmd *cli.Command) error {
					offset := cmd.Int("offset")
					if err := runHexDump(cmd.StringArg("filename"), offset); err != nil {
						return err
					}

					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
