package caesar

import "bytes"

type SplitFunc struct {
	Shift int
}

func NewSplitFunc(shift int) *SplitFunc {
	return &SplitFunc{
		Shift: shift,
	}
}

// dropCR drops a terminal \r from the data.
func (s *SplitFunc) dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == byte('\r'+s.Shift) {
		return data[0 : len(data)-1]
	}
	return data
}

// ScanLines is a split function for a Scanner that returns each line of
// text, stripped of any trailing end-of-line marker. The returned line may
// be empty. The end-of-line marker is one optional carriage return followed
// by one mandatory newline. In regular expression notation, it is `\r?\n`.
// The last non-empty line of input will be returned even if it has no
// newline.
func (s *SplitFunc) ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}
	if i := bytes.IndexByte(data, byte('\n'+s.Shift)); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, s.dropCR(data[0:i]), nil
	}
	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), s.dropCR(data), nil
	}
	// Request more data.
	return 0, nil, nil
}
