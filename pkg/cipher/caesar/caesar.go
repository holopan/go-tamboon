package caesar

import (
	"bufio"
	"io"
	"log"
	"os"
)

type CaesarReader struct {
	reader  io.Reader
	scanner *bufio.Scanner
	file    *os.File
	shift   int
}

func (r *CaesarReader) Decrept(buf []byte) {
	for idx := range buf {
		buf[idx] += byte(r.shift)
	}
}

func (r *CaesarReader) Encrept(buf []byte) {
	for idx := range buf {
		buf[idx] -= byte(r.shift)
	}
}

func NewCaesarReader(filepath string, shift int) *CaesarReader {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal("could not open ", filepath)
	}

	reader := bufio.NewReader(file)

	if reader == nil {
		file.Close()
		log.Fatal("reader is nil")
	}

	s := bufio.NewScanner(reader)
	splitFunc := NewSplitFunc(shift)
	s.Split(splitFunc.ScanLines)

	return &CaesarReader{
		reader:  reader,
		file:    file,
		scanner: s,
		shift:   shift,
	}
}

func (r *CaesarReader) ReadLine() (string, bool) {
	if ok := r.scanner.Scan(); !ok {
		return "", false
	}
	data := r.scanner.Bytes()
	r.Decrept(data)
	return string(data), true
}

func (r *CaesarReader) Write() (string, error) {
	// implement write to file later
	return "", nil
}

func (r *CaesarReader) Close() {
	r.file.Close()
}
