package hexdump

// filename := os.Args[1]
// fileInfo, err := os.Stat(filename)
// if err != nil {
// 	fmt.Fprintln(os.Stderr, fmt.Errorf("file error: %w", err))
// 	os.Exit(1)
// }
// printHeader("File")
// fmt.Printf("%-10s %s\n", "Name:", filename)
// fmt.Printf("%-10s %d bytes\n", "File size:", fileInfo.Size())
// file, fileErr := os.Open(filename)
// if fileErr != nil {
// 	fmt.Fprintln(os.Stderr, fmt.Errorf("file error: %w", fileErr))
// 	os.Exit(1)
// }
// defer file.Close()
// header := make([]byte, 32)
// _, err := file.ReadAt(header, 0x7FC0)
// if err != nil {
// 	panic("PaNiC")
// }
// sheader, err := rom.NewSNESHeader(header)
// if err != nil {
// 	panic("Panic")
// }
// fmt.Println(header)
// printData(sheader)
// header := make([]byte, 32)
// _, err = file.ReadAt(header, 0x7FC0)
// if err != nil {
// 	panic(err)
// }
// for _, v := range header {
// 	fmt.Printf("%02x ", v)
// }
// romTitle := strings.TrimRight(string(header[0:21]), "\x00")
// fmt.Printf("%-10s %s\n", "Rom:", romTitle)
// printHeader("Hexadecimal Dump")
// buffer := make([]byte, 16)
// offset := 0
// for {
// 	bytesRead, err := file.Read(buffer)
// 	if bytesRead > 0 {
// 		// fmt.Printf("%08x: ", offset)
// 		for i := range bytesRead {
// 			// fmt.Printf("%s", getColoredCode(buffer[i]))
// 			if i%2 == 1 {
// 				// fmt.Printf(" ")
// 			}
// 		}
// 		// fmt.Println()
// 		offset += bytesRead
// 	}
// 	if err == io.EOF {
// 		break
// 	}
// 	if err != nil {
// 		fmt.Fprintln(os.Stderr, fmt.Errorf("file error: %w", err))
// 		os.Exit(1)
// 	}
// }
