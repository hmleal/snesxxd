# SNES Inspector

A terminal-based SNES ROM inspector written in Go for exploring cartridge
metadata and raw ROM bytes.

## Usage

Build or run `snesxxd` with one of the commands below. Replace the example ROM
filename with a ROM file you are authorized to inspect.


### Header information

```
snesxxd info Aladdin.sfc

Header (0x7FC0)
-------------------------------------------------------------
Title               : GOOF TROOP
Developer           : 0x08 Capcom
Map Mode            : 0x30 LoROM - Fast (3.58 MHz)
ROM Type            : 0x00 ROM only
ROM Size Exponent   : 0x09 512KB
RAM Size Exponent   : 0x00
Region              : 0x01 USA
Checksum            : 0x5AD0
Checksum Complement : 0xA52F

Raw (Debug)         : [71 79 79 70 32 84 82 79 79 80 32 32 32 32 32 32 32 32 32 32 32 48 0 9 0 1 8 0 47 165 208 90]
```

### Hexadecimal dump

```
snesxxd hex Aladdin.sfc
snesxxd hex --offset 1024 Aladdin.sfc

00000000: 7818 fb5c 0780 809c 0042 9c0b 429c 0c42
00000010: a98f 8d00 219c 0a00 a901 8d0d 42a2 0de0
00000020: 01f0 069e 0042 ca10 f6a9 ff8d 0142 c210
```

## Disclaimer

This project is intended for education, ROM structure research, and personal
experimentation.

Use ROM files that you legally own or are otherwise authorized to inspect. The
project does not include copyrighted ROM data.
