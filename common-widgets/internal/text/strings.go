package text

import (
	"github.com/udan-jayanith/Zbolt/common-widgets/internal/text/internal/textutil"
	"bytes"
	"io"
	"math"
	"strings"
)

func replaceNewLinesWithSpace(text string, start, end int) (string, int, int) {
	var buf strings.Builder
	for {
		pos, len := textutil.FirstLineBreakPositionAndLen(text)
		if len == 0 {
			buf.WriteString(text)
			break
		}
		buf.WriteString(text[:pos])
		origLen := buf.Len()
		buf.WriteString(" ")
		if diff := len - 1; diff > 0 {
			if origLen < start {
				if start >= origLen+len {
					start -= diff
				} else {
					// This is a very rare case, e.g. the position is in between '\r' and '\n'.
					start = origLen + 1
				}
			}
			if origLen < end {
				if end >= origLen+len {
					end -= diff
				} else {
					end = origLen + 1
				}
			}
		}
		text = text[pos+len:]
	}
	text = buf.String()

	return text, start, end
}

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
