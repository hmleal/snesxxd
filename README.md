# SNES Inspector

A terminal-based SNES ROM inspector written in Go for exploring cartridge
metadata and raw ROM bytes.

## Usage

Build or run `snesxxd` with one of the commands below. Replace the example ROM
filename with a ROM file you are authorized to inspect.


### Header information

```
snesxxd info Aladdin.sfc

File
-------------------------------------------------------------
Title               : ALADDIN
Developer           : 0x08 Capcom
Map Mode            : 0x30
ROM Type            : 0x00
ROM Size Exponent   : 0x0B 2048KB
RAM Size Exponent   : 0x00
Region              : 1
Checksum            : 0x060A
Checksum Complement : 0xF9F5
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
