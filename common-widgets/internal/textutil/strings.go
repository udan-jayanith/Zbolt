package textutil

import(
	"io"
	"math"
	"bytes"
)

type stringBuilderWithRange struct {
	buf      []byte
	start    int
	endPlus1 int
	offset   int
}

func (s *stringBuilderWithRange) Reset() {
	s.buf = s.buf[:0]
	s.start = 0
	s.endPlus1 = 0
	s.offset = 0
}

func (s *stringBuilderWithRange) ResetWithRange(start, end int) {
	s.buf = s.buf[:0]
	s.start = start
	s.endPlus1 = end + 1
	s.offset = 0
}

func (s *stringBuilderWithRange) Write(b []byte) (int, error) {
	origN := len(b)
	defer func() {
		s.offset += origN
	}()

	start := s.start
	end := math.MaxInt
	if s.endPlus1 > 0 {
		end = s.endPlus1 - 1
	}

	// Calculate the intersection of [s.offset, s.offset+len(b)) and [start, end).
	idx0 := max(s.offset, start)
	idx1 := min(s.offset+len(b), end)

	if idx0 >= idx1 {
		return origN, nil
	}

	s.buf = append(s.buf, b[idx0-s.offset:idx1-s.offset]...)
	return origN, nil
}

func (s *stringBuilderWithRange) String() string {
	return string(s.buf)
}

func (s *stringBuilderWithRange) Bytes() []byte {
	return s.buf
}

type stringEqualChecker struct {
	str    string
	pos    int
	result bool
}

func (s *stringEqualChecker) Reset(str string) {
	s.str = str
	s.pos = 0
	s.result = true
}

func (s *stringEqualChecker) Result() bool {
	if s.pos != len(s.str) {
		return false
	}
	return s.result
}

func (s *stringEqualChecker) Write(b []byte) (int, error) {
	if s.pos+len(b) > len(s.str) {
		s.result = false
		return 0, io.EOF
	}
	if !bytes.Equal([]byte(s.str[s.pos:s.pos+len(b)]), b) {
		s.result = false
		return 0, io.EOF
	}
	s.pos += len(b)
	return len(b), nil
}

